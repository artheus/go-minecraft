package core

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/item"
	mesh2 "github.com/artheus/go-minecraft/core/mesh"
	"github.com/artheus/go-minecraft/core/texture"
	. "github.com/artheus/go-minecraft/math32"
	"log"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/icexin/gocraft-server/proto"
)

type PlayerState struct {
	X, Y, Z float32
	Rx, Ry  float32
}

type playerState struct {
	PlayerState
	time float64
}

type Player struct {
	s1, s2 playerState

	shader *glhf.Shader
	mesh   *mesh2.Mesh
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
	p.mesh.Draw()
}

func (p *Player) Release() {
	p.mesh.Release()
}

type PlayerRender struct {
	shader  *glhf.Shader
	texture *glhf.Texture
	players map[int32]*Player
}

func NewPlayerRender() (*PlayerRender, error) {
	var (
		err error
	)
	img, rect, err := texture.LoadImage(*texturePath)
	if err != nil {
		return nil, err
	}

	r := &PlayerRender{
		players: make(map[int32]*Player),
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

func (r *PlayerRender) UpdateOrAdd(id int32, s proto.PlayerState) {
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
			item.Tex.Texture(64),
		)
		var mesh *mesh2.Mesh
		mainthread.Call(func() {
			mesh = mesh2.NewMesh(r.shader, blockData)
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

func (r *PlayerRender) Remove(id int32) {
	log.Printf("remove player %d", id)
	p, ok := r.players[id]
	if ok {
		mainthread.CallNonBlock(func() {
			p.Release()
		})
	}
	delete(r.players, id)

}

func (r *PlayerRender) Draw() {
	mat := game.blockRender.get3dmat()
	r.shader.Begin()
	r.texture.Begin()
	for _, p := range r.players {
		p.Draw(mat)
	}
	r.texture.End()
	r.shader.End()
}
