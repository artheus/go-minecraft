package main

import (
	. "github.com/artheus/go-minecraft/math32"
	"github.com/artheus/go-minecraft/types"

	"log"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

func NearBlock(pos mgl32.Vec3) Vec3 {
	return Vec3{
		X: Round(pos.X()),
		Y: Round(pos.Y()),
		Z: Round(pos.Z()),
	}
}

type Chunk struct {
	id     types.ChunkID
	blocks sync.Map // map[Vec3]int
}

func NewChunk(id types.ChunkID) *Chunk {
	c := &Chunk{
		id: id,
	}
	return c
}

func (c *Chunk) ID() types.ChunkID {
	return c.id
}

func (c *Chunk) Block(id Vec3) int {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	w, ok := c.blocks.Load(id)
	if ok {
		return w.(int)
	}
	return 0
}

func (c *Chunk) add(id Vec3, w int) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	c.blocks.Store(id, w)
}

func (c *Chunk) del(id Vec3) {
	if id.ChunkID() != c.id {
		log.Panicf("id %v chunk %v", id, c.id)
	}
	c.blocks.Delete(id)
}

func (c *Chunk) RangeBlocks(f func(id Vec3, w int)) {
	c.blocks.Range(func(key, value interface{}) bool {
		f(key.(Vec3), value.(int))
		return true
	})
}
