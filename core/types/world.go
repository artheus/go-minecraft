package types

import (
	. "github.com/artheus/go-minecraft/math32"
	"github.com/go-gl/mathgl/mgl32"
)

type IWorld interface {
	Collide(pos mgl32.Vec3) (mgl32.Vec3, bool)
	HitTest(pos mgl32.Vec3, vec mgl32.Vec3) (*Vec3, *Vec3)
	Block(id Vec3) int
	BlockChunk(block Vec3) IChunk
	UpdateBlock(id Vec3, tp int)
	IsPlant(tp int) bool
	IsTransparent(tp int) bool
	IsObstacle(tp int) bool
	HasBlock(id Vec3) bool
	Chunk(id Vec3) IChunk
	Chunks(ids []Vec3) []IChunk
}