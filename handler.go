package boltgrpc

import (
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	bolt "github.com/coreos/bbolt"
)

var db *bolt.DB

type Handler struct {
	Path string
}

func (h *Handler) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	h.mustDB()

	rerr := UpdateResponse_NONE

	buckets := req.Buckets
	key := req.Key
	val := req.Val

	if len(buckets) == 0 || len(key) == 0 {
		rerr = UpdateResponse_FAILED
		return &UpdateResponse{Err: rerr}, nil
	}

	db.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error
		for _, item := range buckets {
			if bucket == nil {
				bucket, err = tx.CreateBucketIfNotExists([]byte(item))
			} else {
				bucket, err = bucket.CreateBucketIfNotExists([]byte(item))
			}
		}

		if len(val) > 0 {
			err = bucket.Put(key, val)
			if err != nil {
				rerr = UpdateResponse_FAILED
			}
		} else {
			err = bucket.Delete(key)
			if err != nil {
				rerr = UpdateResponse_FAILED
			}
		}

		return err
	})

	return &UpdateResponse{Err: rerr}, nil
}

func (h *Handler) View(ctx context.Context, req *ViewRequest) (*ViewResponse, error) {
	h.mustDB()

	rerr := ViewResponse_NONE
	rval := []byte("")

	buckets := req.Buckets
	key := req.Key

	if len(buckets) == 0 || len(key) == 0 {
		rerr = ViewResponse_FAILED
		return &ViewResponse{Val: rval, Err: rerr}, nil
	}

	db.View(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		for _, item := range buckets {
			if bucket == nil {
				bucket = tx.Bucket([]byte(item))
			} else {
				bucket = bucket.Bucket([]byte(item))
			}
		}

		if bucket == nil {
			rerr = ViewResponse_FAILED
			return nil
		}

		v := bucket.Get(key)
		rval = make([]byte, len(v))
		copy(rval, v)
		return nil
	})

	return &ViewResponse{Val: rval, Err: rerr}, nil
}

func (h *Handler) Backup(w http.ResponseWriter, r *http.Request) {
	h.mustDB()

	err := db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="backup.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Close() {
	db.Close()
}

func (h *Handler) mustDB() {
	if db != nil {
		return
	}

	var err error
	db, err = bolt.Open(h.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}
