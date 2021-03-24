package status

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/boltdb/bolt"
)

var (
	bucketSystem  = []byte(`system`)
	keyCurrent    = []byte(`current`)
	bucketElement = []byte(`element`)
)

type Status struct {
	db *bolt.DB
}

func New(connect string) (status *Status, e error) {
	db, e := bolt.Open(connect, 0600, &bolt.Options{Timeout: time.Second * 10})
	if e != nil {
		return
	}
	status = &Status{
		db: db,
	}
	return
}
func (s *Status) Close() error {
	return s.db.Close()
}

func (s *Status) Next() (result *Element, e error) {
	e = s.db.Update(func(t *bolt.Tx) error {
		bucket, e := t.CreateBucketIfNotExists(bucketSystem)
		if e != nil {
			return e
		}
		id := s.getCurrentID(bucket)
		if id == 0 {
			// create new
			e = s.setCurrentID(bucket, 1)
			if e != nil {
				return e
			}
			// set element
			ele := &Element{
				ID:      1,
				Step:    Init,
				Created: time.Now(),
			}
			e = s.setElement(t, ele)
			if e != nil {
				return e
			}
			result = ele
		} else {
			ele, e := s.getElement(t, id)
			if e != nil {
				return e
			}
			if ele.Step > Finished {
				return fmt.Errorf(`unknow step %v`, ele.Step)
			} else if ele.Step < Finished {
				result = ele
			} else {
				// create new
				ele := &Element{
					ID:      id + 1,
					Step:    Init,
					Created: time.Now(),
				}
				e = s.setElement(t, ele)
				if e != nil {
					return e
				}
				result = ele
			}
		}
		return nil
	})
	return
}
func (s *Status) getCurrentID(bucket *bolt.Bucket) uint64 {
	val := bucket.Get(keyCurrent)
	if len(val) == 8 {
		return binary.LittleEndian.Uint64(val)
	}
	return 0
}
func (s *Status) setCurrentID(bucket *bolt.Bucket, id uint64) error {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, id)
	return bucket.Put(keyCurrent, key)
}
func (s *Status) getElement(t *bolt.Tx, id uint64) (*Element, error) {
	bucket := t.Bucket(bucketElement)
	if bucket == nil {
		return nil, fmt.Errorf(`bucket %s not exists`, bucketElement)
	}
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, id)
	var ele Element
	b := bucket.Get(key)
	e := Unmarshal(b, &ele)
	if e != nil {
		return nil, e
	}
	return &ele, nil
}
func (s *Status) setElement(t *bolt.Tx, ele *Element) error {
	val, e := Marshal(ele)
	if e != nil {
		return e
	}
	bucket, e := t.CreateBucketIfNotExists(bucketElement)
	if e != nil {
		return e
	}
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, ele.ID)
	return bucket.Put(key, val)
}
func (s *Status) GenerateDescription(filename string) error {
	var obj struct {
		CurrentID uint64
		Node      []exportElement
	}
	e := s.db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket(bucketSystem)
		if bucket == nil {
			return nil
		}
		currentID := s.getCurrentID(bucket)
		obj.CurrentID = currentID
		for i := uint64(1); i <= currentID; i++ {
			ele, e := s.getElement(t, i)
			if e != nil {
				return e
			}
			obj.Node = append(obj.Node, exportElement{
				ID:       ele.ID,
				Step:     ele.Step.String(),
				Created:  ele.Created,
				Backuped: ele.Created,
			})
		}
		return nil
	})
	if e != nil {
		return e
	}
	b, e := json.MarshalIndent(obj, "", "\t")
	if e != nil {
		return e
	}
	return ioutil.WriteFile(filename, b, 0666)
}
func (s *Status) Update(ele *Element) (e error) {
	return s.db.Update(func(t *bolt.Tx) error {
		return s.setElement(t, ele)
	})
}
