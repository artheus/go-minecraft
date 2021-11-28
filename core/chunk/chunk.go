package chunk

import (
	. "github.com/artheus/go-minecraft/core/block"
	. "github.com/artheus/go-minecraft/math/f32"
	"log"
	"math"
	"sync"
)

type Chunk struct {
	id       Vec3
	segments sync.Map // map[uint8]*Segment
}

func NewChunk(id Vec3) *Chunk {
	c := &Chunk{
		id: id,
		segments: sync.Map{},
	}
	return c
}

func (c *Chunk) ID() Vec3 {
	return c.id
}

func (c *Chunk) Block(pos Vec3) (block *Block) {
	if pos.ChunkID() != c.id {
		return GetBlock(AirID)
		//log.Fatalf("block %v is not in chunk %v", pos, c.id)
	}

	var ok bool
	var seg interface{}

	if seg, ok = c.segments.Load(segmentId(pos.Y)); !ok {
		return GetBlock(AirID)
	}

	var w interface{}

	if w, ok = seg.(*Segment).blocks.Load(pos); !ok || w.(*Block) == nil {
		return GetBlock(AirID)
	}

	return w.(*Block)
}

func (c *Chunk) Add(id Vec3, w *Block) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}

	var ok bool
	var seg interface{}
	var segID = segmentId(id.Y)

	if seg, ok = c.segments.Load(segID); !ok {
		seg = &Segment{
			blocks: sync.Map{},
		}
		c.segments.Store(segID, seg)
	}

	seg.(*Segment).blocks.Store(id, w)
}

func (c *Chunk) Del(id Vec3) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}

	if seg, ok := c.segments.Load(segmentId(id.Y)); ok {
		seg.(*Segment).blocks.Delete(id)
	}
}

func (c *Chunk) RangeBlocks(f func(id Vec3, w *Block)) {
	c.segments.Range(func(key, value interface{}) bool {
		value.(*Segment).blocks.Range(func(sk, sv interface{}) bool {
			f(sk.(Vec3), sv.(*Block))
			return true
		})
		return true
	})
}

func segmentId(y float32) uint8 {
	return uint8(math.Mod(float64(y), segmentHeight))
}