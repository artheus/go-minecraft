package chunk

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/texture"
	"github.com/artheus/go-minecraft/core/thread"
	. "github.com/artheus/go-minecraft/math32"
	"github.com/faiface/glhf"
	"image"
	"sync"
)

// Thread for managing all active and visible chunks around the player
type Thread struct {
	shader *glhf.Shader
	texture *glhf.Texture
	chunks sync.Map // map[types.ChunkID]*Chunk
}

const (
	texturePath = "texture.png"
)

func (t *Thread) Init() thread.InitFunc {
	return func() (err error) {
		var img []uint8
		var rect image.Rectangle

		if img, rect, err = texture.LoadImage(texturePath); err != nil {
			return err
		}

		t.shader, err = glhf.NewShader(glhf.AttrFormat{
			glhf.Attr{Name: "pos", Type: glhf.Vec3},
			glhf.Attr{Name: "tex", Type: glhf.Vec2},
			glhf.Attr{Name: "normal", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			glhf.Attr{Name: "matrix", Type: glhf.Mat4},
			glhf.Attr{Name: "camera", Type: glhf.Vec3},
			glhf.Attr{Name: "fogdis", Type: glhf.Float},
		}, blockVertexSource, blockFragmentSource)

		if err != nil {
			return err
		}
		t.texture = glhf.NewTexture(rect.Dx(), rect.Dy(), false, img)

		return nil
	}
}

func (t *Thread) Tick() thread.TickFunc {
	return func() {
		// Render chunks
		t.chunks.Range(func(id, chunk interface{}) bool {
			c := chunk.(*Chunk)

			// Render blocks in chunk
			c.blocks.Range(func(pos, b interface{}) bool {
				t.renderBlock(pos.(Vec3), b.(*block.Block))
				return true
			})

			return true
		})
	}
}

// renderBlock should render visible blocks of a chunk
func (t *Thread) renderBlock(pos Vec3, b *block.Block) {
	t.shader.Begin()
	t.texture.Begin()
	defer t.shader.End()
	defer t.texture.End()

	// DRAW / RENDER Chunks
}



