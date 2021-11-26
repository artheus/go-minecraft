package types

import "github.com/go-gl/mathgl/mgl32"

type ICamera interface {
	Matrix() mgl32.Mat4
	SetPos(pos mgl32.Vec3)
	Pos() mgl32.Vec3
	Front() mgl32.Vec3
	FlipFlying()
	Flying() bool
	OnAngleChange(dx, dy float32)

	State() PlayerState
	Restore(state PlayerState)
}