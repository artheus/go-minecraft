package world

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/chunk"
	"github.com/artheus/go-minecraft/core/ctx"
	"github.com/artheus/go-minecraft/core/game/rpc"
	"github.com/artheus/go-minecraft/core/game/store"
	"github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"
	"log"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hashicorp/golang-lru"
)

type World struct {
	mutex  sync.Mutex
	chunks *lru.Cache // map[Vec3]*Chunk
	ctx    *ctx.Context
}

func NewWorld(ctx *ctx.Context) *World {
	m := (*chunk.RenderRadius) * (*chunk.RenderRadius) * 4
	chunks, _ := lru.New(m)
	return &World{
		chunks: chunks,
		ctx:    ctx,
	}
}

func (w *World) loadChunk(id Vec3) (*chunk.Chunk, bool) {
	c, ok := w.chunks.Get(id)
	if !ok {
		return nil, false
	}
	return c.(*chunk.Chunk), true
}

func (w *World) storeChunk(id Vec3, chunk *chunk.Chunk) {
	w.chunks.Add(id, chunk)
}

func (w *World) Collide(pos mgl32.Vec3) (mgl32.Vec3, bool) {
	x, y, z := pos.X(), pos.Y(), pos.Z()
	nx, ny, nz := Round(pos.X()), Round(pos.Y()), Round(pos.Z())
	const pad = 0.25

	head := Vec3{
		X: nx,
		Y: ny,
		Z: nz,
	}
	foot := head.Down()

	stop := false
	for _, b := range []Vec3{foot, head} {
		if w.Block(b.Left()).Obstacle && x < nx && nx-x > pad {
			x = nx - pad
		}
		if w.Block(b.Right()).Obstacle && x > nx && x-nx > pad {
			x = nx + pad
		}
		if w.Block(b.Down()).Obstacle && y < ny && ny-y > pad {
			y = ny - pad
			stop = true
		}
		if w.Block(b.Up()).Obstacle && y > ny && y-ny > pad {
			y = ny + pad
			stop = true
		}
		if w.Block(b.Back()).Obstacle && z < nz && nz-z > pad {
			z = nz - pad
		}
		if w.Block(b.Front()).Obstacle && z > nz && z-nz > pad {
			z = nz + pad
		}
	}
	return mgl32.Vec3{x, y, z}, stop
}

func (w *World) HitTest(pos mgl32.Vec3, vec mgl32.Vec3) (*Vec3, *Vec3) {
	var (
		maxLen = float32(8.0)
		step   = float32(0.125)

		block, prev Vec3
		pprev       *Vec3
	)

	for len := float32(0); len < maxLen; len += step {
		block = chunk.NearBlock(pos.Add(vec.Mul(len)))
		if prev != block && w.HasBlock(block) {
			return &block, pprev
		}
		prev = block
		pprev = &prev
	}
	return nil, nil
}

func (w *World) Block(pos Vec3) *block.Block {
	chunk := w.BlockChunk(pos)
	if chunk == nil {
		return block.GetBlock(block.AirID)
	}

	return chunk.Block(pos)
}

func (w *World) BlockChunk(block Vec3) types.IChunk {
	cid := block.ChunkID()
	chunk, ok := w.loadChunk(cid)
	if !ok {
		return nil
	}
	return chunk
}

func (w *World) UpdateBlock(id Vec3, tp *block.Block) {
	chunk := w.BlockChunk(id)
	if chunk != nil {
		if tp.ID != block.AirID {
			chunk.Add(id, tp)
		} else {
			chunk.Del(id)
		}
	}
	store.Storage.UpdateBlock(id, tp)
}

func (w *World) HasBlock(id Vec3) bool {
	tp := w.Block(id)
	return tp != nil && tp.ID != block.AirID
}

func (w *World) Chunk(id Vec3) types.IChunk {
	p, ok := w.loadChunk(id)
	if ok {
		return p
	}
	chunk := chunk.NewChunk(id)
	blocks := makeChunkMap(id)
	for block, tp := range blocks {
		chunk.Add(block, tp)
	}
	err := store.Storage.RangeBlocks(id, func(bid Vec3, w *block.Block) {
		if w.ID == block.AirID {
			chunk.Del(bid)
			return
		}
		chunk.Add(bid, w)
	})
	if err != nil {
		log.Printf("fetch chunk(%v) from db error:%s", id, err)
		return nil
	}
	rpc.ClientFetchChunk(id, func(bid Vec3, w *block.Block) {
		if w.ID == block.AirID {
			chunk.Del(bid)
			return
		}
		chunk.Add(bid, w)
		store.Storage.UpdateBlock(bid, w)
	})
	w.storeChunk(id, chunk)
	return chunk
}

func (w *World) Chunks(ids []Vec3) []types.IChunk {
	ch := make(chan types.IChunk)
	var chunks []types.IChunk
	for _, id := range ids {
		go func(id Vec3) {
			ch <- w.Chunk(id)
		}(id)
	}
	for range ids {
		chunk := <-ch
		if chunk != nil {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func makeChunkMap(cid Vec3) map[Vec3]*block.Block {
	var (
		grassBlock = block.GetBlock(block.GrassBlockID)
		dirtBlock  = block.GetBlock(block.DirtID)
		waterBlock = block.GetBlock(block.SandID)
		grass      = block.GetBlock(block.GrassID)
		leaves     = block.GetBlock(block.LeavesID)
		wood       = block.GetBlock(block.WoodID)
		dandelion  = block.GetBlock(block.DandelionID)
		cloud      = block.GetBlock(block.CloudID)
	)
	m := make(map[Vec3]*block.Block)
	p, q := cid.X, cid.Z
	for dx := 0; dx < ChunkWidth; dx++ {
		for dz := 0; dz < ChunkWidth; dz++ {
			x, z := int(p)*ChunkWidth+dx, int(q)*ChunkWidth+dz
			f := Noise2(float32(x)*0.01, float32(z)*0.01, 4, 0.5, 2)
			g := Noise2(float32(-x)*0.01, float32(-z)*0.01, 2, 0.9, 2)
			mh := int(g*32 + 16)
			h := int(f * float32(mh))
			w := dirtBlock
			if h <= 12 {
				h = 12
				w = waterBlock
			}
			// grass and sand
			for y := 0; y < h; y++ {
				if y == h-1 && w == dirtBlock {
					m[Vec3{X: float32(x), Y: float32(y), Z: float32(z)}] = grassBlock
					continue
				}

				m[Vec3{X: float32(x), Y: float32(y), Z: float32(z)}] = w
			}

			// flowers
			if w == dirtBlock {
				if Noise2(-float32(x)*0.1, float32(z)*0.1, 4, 0.8, 2) > 0.6 {
					m[Vec3{X: float32(x), Y: float32(h), Z: float32(z)}] = grass

				}
				if Noise2(float32(x)*0.05, float32(-z)*0.05, 4, 0.8, 2) > 0.7 {
					//w1 := 18 + int(Noise2(float32(x)*0.1, float32(z)*0.1, 4, 0.8, 2)*7)
					m[Vec3{X: float32(x), Y: float32(h), Z: float32(z)}] = dandelion
				}
			}

			// tree
			if w == dirtBlock {
				ok := true
				if dx-4 < 0 || dz-4 < 0 ||
					dx+4 > ChunkWidth || dz+4 > ChunkWidth {
					ok = false
				}
				if ok && Noise2(float32(x), float32(z), 6, 0.5, 2) > 0.79 {
					for y := h + 3; y < h+8; y++ {
						for ox := -3; ox <= 3; ox++ {
							for oz := -3; oz <= 3; oz++ {
								d := ox*ox + oz*oz + (y-h-4)*(y-h-4)
								if d < 11 {
									m[Vec3{X: float32(x + ox), Y: float32(y), Z: float32(z + oz)}] = leaves
								}
							}
						}
					}
					for y := h; y < h+7; y++ {
						m[Vec3{X: float32(x), Y: float32(y), Z: float32(z)}] = wood
					}
				}
			}

			// cloud
			for y := 64; y < 72; y++ {
				if Noise3(float32(x)*0.01, float32(y)*0.1, float32(z)*0.01, 8, 0.5, 2) > 0.69 {
					m[Vec3{X: float32(x), Y: float32(y), Z: float32(z)}] = cloud
				}
			}
		}
	}
	return m
}
