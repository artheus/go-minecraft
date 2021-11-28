package store

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"

	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"log"

	"github.com/boltdb/bolt"
)

var (
	dbpath = flag.String("db", "gocraft.db", "db file name")
)

var (
	blockBucket  = []byte("block")
	chunkBucket  = []byte("chunk")
	cameraBucket = []byte("camera")

	Storage *Store
)

func InitStore() error {
	var path string
	if *dbpath != "" {
		path = *dbpath
	}
	//if *core.serverAddr != "" {
	//	path = fmt.Sprintf("cache_%s.db", *core.serverAddr)
	//}
	if path == "" {
		return errors.New("empty db path")
	}
	var err error
	Storage, err = NewStore(path)
	return err
}

type Store struct {
	db *bolt.DB
}

func NewStore(p string) (*Store, error) {
	db, err := bolt.Open(p, 0666, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(blockBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(chunkBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(cameraBucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	db.NoSync = true
	return &Store{
		db: db,
	}, nil
}

func (s *Store) UpdateBlock(id Vec3, w *block.Block) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		log.Printf("put %v -> %d", id, w)
		bkt := tx.Bucket(blockBucket)
		cid := id.ChunkID()
		key := encodeBlockDbKey(cid, id)
		value := encodeBlockDbValue(w)
		return bkt.Put(key, value)
	})
}

func (s *Store) UpdatePlayerState(state types.PlayerState) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(cameraBucket)
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, &state)
		bkt.Put(cameraBucket, buf.Bytes())
		return nil
	})
}

func (s *Store) GetPlayerState() types.PlayerState {
	var state types.PlayerState
	state.Y = 16
	s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(cameraBucket)
		value := bkt.Get(cameraBucket)
		if value == nil {
			return nil
		}
		buf := bytes.NewBuffer(value)
		binary.Read(buf, binary.LittleEndian, &state)
		return nil
	})
	return state
}

func (s *Store) RangeBlocks(id Vec3, f func(bid Vec3, w *block.Block)) error {
	return s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockBucket)
		startkey := encodeBlockDbKey(id, Vec3{0, 0, 0})
		iter := bkt.Cursor()
		for k, v := iter.Seek(startkey); k != nil; k, v = iter.Next() {
			cid, bid := decodeBlockDbKey(k)
			if cid != id {
				break
			}
			w := decodeBlockDbValue(v)
			f(bid, block.GetBlock(w))
		}
		return nil
	})
}

func (s *Store) UpdateChunkVersion(id Vec3, version string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(chunkBucket)
		key := encodeChunkID(id)
		return bkt.Put(key, []byte(version))
	})
}

func (s *Store) GetChunkVersion(id Vec3) string {
	var version string
	s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(chunkBucket)
		key := encodeChunkID(id)
		v := bkt.Get(key)
		if v != nil {
			version = string(v)
		}
		return nil
	})
	return version
}

func (s *Store) Close() {
	s.db.Sync()
	s.db.Close()
}

func encodeVec3(v Vec3) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, [...]float32{v.X, v.Y, v.Z})
	return buf.Bytes()
}

func encodeChunkID(v Vec3) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, [...]int32{int32(v.X), int32(v.Y), int32(v.Z)})
	return buf.Bytes()
}

func encodeBlockDbKey(cid Vec3, bid Vec3) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, [...]int32{int32(cid.X), int32(cid.Z)})
	binary.Write(buf, binary.LittleEndian, [...]int32{int32(bid.X), int32(bid.Y), int32(bid.Z)})
	return buf.Bytes()
}

func decodeBlockDbKey(b []byte) (Vec3, Vec3) {
	if len(b) != 4*5 {
		log.Panicf("bad db key length:%d", len(b))
	}
	buf := bytes.NewBuffer(b)
	var arr [5]int32
	binary.Read(buf, binary.LittleEndian, &arr)

	cid := Vec3{X: float32(arr[0]), Z: float32(arr[1])}
	bid := Vec3{X: float32(arr[2]), Y: float32(arr[3]), Z: float32(arr[4])}
	if bid.ChunkID() != cid {
		log.Panicf("bad db key: cid:%v, bid:%v", cid, bid)
	}
	return cid, bid
}

func encodeBlockDbValue(w *block.Block) []byte {
	return []byte(w.ID)
}

func decodeBlockDbValue(b []byte) string {
	return string(b)
}
