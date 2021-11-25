package block

import (
	"github.com/artheus/go-minecraft/core/mesh"
)

type Block struct {
	mesh *mesh.Mesh

	visible bool
}

func (b *Block) Invisible() bool {
	return !b.visible
}