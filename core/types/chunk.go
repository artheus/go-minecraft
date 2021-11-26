package types

import . "github.com/artheus/go-minecraft/math32"

type IChunk interface {
	ID() Vec3
	Block(id Vec3) int
	Add(id Vec3, w int)
	Del(id Vec3)
	RangeBlocks(f func(id Vec3, w int))
}