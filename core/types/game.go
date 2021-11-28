package types

import (
	"github.com/artheus/go-minecraft/math/f32"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type IGameApplication interface {
	World() IWorld
	Camera() ICamera
	Window() *glfw.Window

	CurrentBlockid() f32.Vec3
	ShouldClose() bool
	Update()

	LineRenderer() ILineRenderer
	PlayerRenderer() IPlayerRenderer
	ChunkRenderer() IChunkRenderer
}
