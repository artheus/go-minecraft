package types

import (
	"github.com/artheus/go-minecraft/math/f32"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/icexin/gocraft-server/proto"
)

type IRenderer interface {
	Render()
}

type IPlayerRenderer interface {
	IRenderer

	UpdateOrAdd(id int32, s proto.PlayerState)
	Remove(id int32)
}

type IChunkRenderer interface {
	IRenderer

	UpdateItem(w string)
	Get3dMat() mgl32.Mat4
	Get2dMat() mgl32.Mat4
	DirtyChunk(id f32.Vec3)
	UpdateLoop()
}

type ILineRenderer interface {
	IRenderer
}