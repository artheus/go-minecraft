package player

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/ctx"
	"github.com/artheus/go-minecraft/core/item"
	"github.com/artheus/go-minecraft/core/texture"
	. "github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"
	"log"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/icexin/gocraft-server/proto"
)

type playerState struct {
	PlayerState
	time float64
}

type Player struct {
	s1, s2 playerState

	shader *glhf.Shader
	mesh   *Mesh
}

// Linear interpolation to calculate player position
func (p *Player) computeMat() mgl32.Mat4 {
	t1 := p.s2.time - p.s1.time
	t2 := glfw.GetTime() - p.s2.time
	t := Min(float32(t2/t1), 1)

	x := Mix(p.s1.X, p.s2.X, t)
	y := Mix(p.s1.Y, p.s2.Y, t)
	z := Mix(p.s1.Z, p.s2.Z, t)
	rx := Mix(p.s1.Rx, p.s2.Rx, t)
	ry := Mix(p.s1.Ry, p.s2.Ry, t)

	front := mgl32.Vec3{
		Cos(Radian(ry)) * Cos(Radian(rx)),
		Sin(Radian(ry)),
		Cos(Radian(ry)) * Sin(Radian(rx)),
	}.Normalize()
	right := front.Cross(mgl32.Vec3{0, 1, 0})
	up := right.Cross(front).Normalize()
	pos := mgl32.Vec3{x, y, z}
	return mgl32.LookAtV(pos, pos.Add(front), up).Inv()
}

func (p *Player) UpdateState(s playerState) {
	p.s1, p.s2 = p.s2, s
}

func (p *Player) Draw(mat mgl32.Mat4) {
	mat = mat.Mul4(p.computeMat())

	p.shader.SetUniformAttr(0, mat)
	p.mesh.Render()
}

func (p *Player) Release() {
	p.mesh.Release()
}

type PlayerRenderer struct {
	ctx     *ctx.Context
	shader  *glhf.Shader
	texture *glhf.Texture
	players map[int32]*Player
}

func NewPlayerRenderer(ctx *ctx.Context) (*PlayerRenderer, error) {
	var (
		err error
	)
	img, rect, err := texture.LoadImage(*texture.TexturePath)
	if err != nil {
		return nil, err
	}

	r := &PlayerRenderer{
		players: make(map[int32]*Player),
		ctx:     ctx,
	}
	mainthread.Call(func() {
		r.shader, err = glhf.NewShader(glhf.AttrFormat{
			glhf.Attr{Name: "pos", Type: glhf.Vec3},
			glhf.Attr{Name: "tex", Type: glhf.Vec2},
			glhf.Attr{Name: "normal", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			glhf.Attr{Name: "matrix", Type: glhf.Mat4},
		}, playerVertexSource, playerFragmentSource)

		if err != nil {
			return
		}
		r.texture = glhf.NewTexture(rect.Dx(), rect.Dy(), false, img)

	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PlayerRenderer) UpdateOrAdd(id int32, s proto.PlayerState) {
	state := playerState{
		PlayerState: PlayerState{
			X:  s.X,
			Y:  s.Y,
			Z:  s.Z,
			Rx: s.Rx,
			Ry: s.Ry,
		},
		time: glfw.GetTime(),
	}

	p, ok := r.players[id]
	if !ok {
		log.Printf("add new player %d", id)
		blockData := block.BlockData(
			[]float32{},
			block.Sides(
				true,
				true,
				true,
				true,
				true,
				true,
			),
			Vec3{
				0,
				0,
				0,
			},
			item.Tex.Texture("core:player"),
		)
		var mesh *Mesh
		mainthread.Call(func() {
			mesh = NewMesh(r.shader, blockData)
		})
		p = &Player{
			shader: r.shader,
			mesh:   mesh,
		}
		r.players[id] = p
		p.s1 = state
	}
	p.UpdateState(state)
}

func (r *PlayerRenderer) Remove(id int32) {
	log.Printf("remove player %d", id)
	p, ok := r.players[id]
	if ok {
		mainthread.CallNonBlock(func() {
			p.Release()
		})
	}
	delete(r.players, id)

}

func (r *PlayerRenderer) Render() {
	mat := r.ctx.Game().ChunkRenderer().Get3dMat()
	r.shader.Begin()
	r.texture.Begin()
	for _, p := range r.players {
		p.Draw(mat)
	}
	r.texture.End()
	r.shader.End()
}
