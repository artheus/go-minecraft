package types

import (
	"github.com/artheus/go-minecraft/core/block"
	. "github.com/artheus/go-minecraft/math/f32"
)

type IChunk interface {
	ID() Vec3
	Block(id Vec3) *block.Block
	Add(id Vec3, w *block.Block)
	Del(id Vec3)
	RangeBlocks(f func(id Vec3, w *block.Block))
}