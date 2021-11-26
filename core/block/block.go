package block

import (
	"github.com/artheus/go-minecraft/core/types"
)

type Block struct {
	mesh *types.Mesh

	visible bool
}

func (b *Block) Invisible() bool {
	return !b.visible
}