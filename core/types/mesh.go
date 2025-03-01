package types

import (
	"github.com/artheus/go-minecraft/math/f32"
	"github.com/faiface/glhf"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Mesh struct {
	vao, vbo uint32
	faces int
	Id    f32.Vec3
	Dirty bool
}

func NewMesh(shader *glhf.Shader, data []float32) *Mesh {
	m := new(Mesh)

	m.faces = len(data) / (shader.VertexFormat().Size() / 4) / 6

	if m.faces == 0 {
		return m
	}

	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	offset := 0
	for _, attr := range shader.VertexFormat() {
		loc := gl.GetAttribLocation(shader.ID(), gl.Str(attr.Name+"\x00"))
		var size int32
		switch attr.Type {
		case glhf.Float:
			size = 1
		case glhf.Vec2:
			size = 2
		case glhf.Vec3:
			size = 3
		case glhf.Vec4:
			size = 4
		}
		gl.VertexAttribPointer(
			uint32(loc),
			size,
			gl.FLOAT,
			false,
			int32(shader.VertexFormat().Size()),
			gl.PtrOffset(offset),
		)
		gl.EnableVertexAttribArray(uint32(loc))
		offset += attr.Type.Size()
	}
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return m
}

func (m *Mesh) Faces() int {
	return m.faces
}

func (m *Mesh) Render() {
	if m.vao != 0 {
		gl.BindVertexArray(m.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(m.faces)*6)
		gl.BindVertexArray(0)
	}
}

func (m *Mesh) Release() {
	if m.vao != 0 {
		gl.DeleteVertexArrays(1, &m.vao)
		gl.DeleteBuffers(1, &m.vbo)
		m.vao = 0
		m.vbo = 0
	}
}
