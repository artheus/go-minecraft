package chunk

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/math/f32"
)

type Action int

const (
	ActionAdd Action = iota
	ActionDelete
)

type ChunkAction struct {
	pos f32.Vec3
	block *block.Block
	action Action
}
