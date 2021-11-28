package block

import (
	"github.com/pkg/errors"
	"sync"
)

type register struct {
	blocks map[string]*Block
	mx *sync.Mutex
}

var instance *register

func InitRegister() {
	instance = &register{
		blocks: map[string]*Block{},
		mx: &sync.Mutex{},
	}

	LoadBlocks()
}

func AddBlock(block *Block) error {
	if block.ID == "" {
		return errors.New("can't register block with an empty id")
	}

	instance.mx.Lock()
	defer instance.mx.Unlock()

	if _, ok := instance.blocks[block.ID]; ok {
		return errors.Errorf("block with id %s is already registered", block.ID)
	}

	instance.blocks[block.ID] = block

	return nil
}

func GetBlock(id string) (block *Block) {
	instance.mx.Lock()
	defer instance.mx.Unlock()

	return instance.blocks[id]
}

func RangeBlocks(rangeFunc func(block *Block) bool) {
	instance.mx.Lock()
	defer instance.mx.Unlock()

	for _, block := range instance.blocks {
		if f := rangeFunc(block); !f {
			break
		}
	}
}
