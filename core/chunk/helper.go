package chunk

import (
	. "github.com/artheus/go-minecraft/math32"
	"github.com/go-gl/mathgl/mgl32"
)

func NearBlock(pos mgl32.Vec3) Vec3 {
	return Vec3{
		X: Round(pos.X()),
		Y: Round(pos.Y()),
		Z: Round(pos.Z()),
	}
}
