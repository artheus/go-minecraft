package math32

const (
	ChunkWidth = 32
)

type Vec2 struct {
	X, Y float32
}

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) Left() Vec3 {
	return Vec3{v.X - 1, v.Y, v.Z}
}

func (v Vec3) Right() Vec3 {
	return Vec3{v.X + 1, v.Y, v.Z}
}

func (v Vec3) Up() Vec3 {
	return Vec3{v.X, v.Y + 1, v.Z}
}

func (v Vec3) Down() Vec3 {
	return Vec3{v.X, v.Y - 1, v.Z}
}

func (v Vec3) Front() Vec3 {
	return Vec3{v.X, v.Y, v.Z + 1}
}

func (v Vec3) Back() Vec3 {
	return Vec3{v.X, v.Y, v.Z - 1}
}

func (v Vec3) ChunkID() Vec3 {
	return Vec3{
		X: Floor(v.X / ChunkWidth),
		Z: Floor(v.Z / ChunkWidth),
	}
}

type Vec4 struct {
	X, Y, Z, W float32
}
