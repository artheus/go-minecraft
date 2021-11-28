package chunk

import (
	"sync"
)

const (
	segmentHeight = 8
)

type Segment struct {
	blocks sync.Map // map[f32.Vec3]*block.Block
	visible bool
}
