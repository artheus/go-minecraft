package block

import (
	"github.com/artheus/go-minecraft/core/texture"
	. "github.com/artheus/go-minecraft/math/f32"
)

func BlockData(vertices []float32, show Side, block Vec3, tex *texture.BlockTexture) []float32 {
	l, r := tex.Left, tex.Right
	u, d := tex.Up, tex.Down
	f, b := tex.Front, tex.Back
	x, y, z := block.X, block.Y, block.Z
	if show.Left {
		vertices = append(vertices, []float32{
			// left
			x - 0.5, y - 0.5, z - 0.5, l[0][0], l[0][1], -1, 0, 0,
			x - 0.5, y - 0.5, z + 0.5, l[1][0], l[1][1], -1, 0, 0,
			x - 0.5, y + 0.5, z + 0.5, l[2][0], l[2][1], -1, 0, 0,
			x - 0.5, y + 0.5, z + 0.5, l[3][0], l[3][1], -1, 0, 0,
			x - 0.5, y + 0.5, z - 0.5, l[4][0], l[4][1], -1, 0, 0,
			x - 0.5, y - 0.5, z - 0.5, l[5][0], l[5][1], -1, 0, 0,
		}...)
	}
	if show.Right {
		vertices = append(vertices, []float32{
			// right
			x + 0.5, y - 0.5, z + 0.5, r[0][0], r[0][1], 1, 0, 0,
			x + 0.5, y - 0.5, z - 0.5, r[1][0], r[1][1], 1, 0, 0,
			x + 0.5, y + 0.5, z - 0.5, r[2][0], r[2][1], 1, 0, 0,
			x + 0.5, y + 0.5, z - 0.5, r[3][0], r[3][1], 1, 0, 0,
			x + 0.5, y + 0.5, z + 0.5, r[4][0], r[4][1], 1, 0, 0,
			x + 0.5, y - 0.5, z + 0.5, r[5][0], r[5][1], 1, 0, 0,
		}...)
	}
	if show.Up {
		vertices = append(vertices, []float32{
			// top
			x - 0.5, y + 0.5, z + 0.5, u[0][0], u[0][1], 0, 1, 0,
			x + 0.5, y + 0.5, z + 0.5, u[1][0], u[1][1], 0, 1, 0,
			x + 0.5, y + 0.5, z - 0.5, u[2][0], u[2][1], 0, 1, 0,
			x + 0.5, y + 0.5, z - 0.5, u[3][0], u[3][1], 0, 1, 0,
			x - 0.5, y + 0.5, z - 0.5, u[4][0], u[4][1], 0, 1, 0,
			x - 0.5, y + 0.5, z + 0.5, u[5][0], u[5][1], 0, 1, 0,
		}...)
	}

	if show.Down {
		vertices = append(vertices, []float32{
			// bottom
			x - 0.5, y - 0.5, z - 0.5, d[0][0], d[0][1], 0, -1, 0,
			x + 0.5, y - 0.5, z - 0.5, d[1][0], d[1][1], 0, -1, 0,
			x + 0.5, y - 0.5, z + 0.5, d[2][0], d[2][1], 0, -1, 0,
			x + 0.5, y - 0.5, z + 0.5, d[3][0], d[3][1], 0, -1, 0,
			x - 0.5, y - 0.5, z + 0.5, d[4][0], d[4][1], 0, -1, 0,
			x - 0.5, y - 0.5, z - 0.5, d[5][0], d[5][1], 0, -1, 0,
		}...)
	}

	if show.Front {
		vertices = append(vertices, []float32{
			// front
			x - 0.5, y - 0.5, z + 0.5, f[0][0], f[0][1], 0, 0, 1,
			x + 0.5, y - 0.5, z + 0.5, f[1][0], f[1][1], 0, 0, 1,
			x + 0.5, y + 0.5, z + 0.5, f[2][0], f[2][1], 0, 0, 1,
			x + 0.5, y + 0.5, z + 0.5, f[3][0], f[3][1], 0, 0, 1,
			x - 0.5, y + 0.5, z + 0.5, f[4][0], f[4][1], 0, 0, 1,
			x - 0.5, y - 0.5, z + 0.5, f[5][0], f[5][1], 0, 0, 1,
		}...)
	}

	if show.Back {
		vertices = append(vertices, []float32{
			// back
			x + 0.5, y - 0.5, z - 0.5, b[0][0], b[0][1], 0, 0, -1,
			x - 0.5, y - 0.5, z - 0.5, b[1][0], b[1][1], 0, 0, -1,
			x - 0.5, y + 0.5, z - 0.5, b[2][0], b[2][1], 0, 0, -1,
			x - 0.5, y + 0.5, z - 0.5, b[3][0], b[3][1], 0, 0, -1,
			x + 0.5, y + 0.5, z - 0.5, b[4][0], b[4][1], 0, 0, -1,
			x + 0.5, y - 0.5, z - 0.5, b[5][0], b[5][1], 0, 0, -1,
		}...)
	}

	return vertices
}

func PlantData(vertices []float32, _ Side, block Vec3, tex *texture.BlockTexture) []float32 {
	l, r := tex.Left, tex.Right
	f, b := tex.Front, tex.Back
	x, y, z := block.X, block.Y, block.Z
	vertices = append(vertices, []float32{
		// left
		x, y - 0.5, z - 0.5, l[0][0], l[0][1], -1, 0, 0,
		x, y - 0.5, z + 0.5, l[1][0], l[1][1], -1, 0, 0,
		x, y + 0.5, z + 0.5, l[2][0], l[2][1], -1, 0, 0,
		x, y + 0.5, z + 0.5, l[3][0], l[3][1], -1, 0, 0,
		x, y + 0.5, z - 0.5, l[4][0], l[4][1], -1, 0, 0,
		x, y - 0.5, z - 0.5, l[5][0], l[5][1], -1, 0, 0,
	}...)
	vertices = append(vertices, []float32{
		// right
		x, y - 0.5, z + 0.5, r[0][0], r[0][1], 1, 0, 0,
		x, y - 0.5, z - 0.5, r[1][0], r[1][1], 1, 0, 0,
		x, y + 0.5, z - 0.5, r[2][0], r[2][1], 1, 0, 0,
		x, y + 0.5, z - 0.5, r[3][0], r[3][1], 1, 0, 0,
		x, y + 0.5, z + 0.5, r[4][0], r[4][1], 1, 0, 0,
		x, y - 0.5, z + 0.5, r[5][0], r[5][1], 1, 0, 0,
	}...)

	vertices = append(vertices, []float32{
		// front
		x - 0.5, y - 0.5, z, f[0][0], f[0][1], 0, 0, 1,
		x + 0.5, y - 0.5, z, f[1][0], f[1][1], 0, 0, 1,
		x + 0.5, y + 0.5, z, f[2][0], f[2][1], 0, 0, 1,
		x + 0.5, y + 0.5, z, f[3][0], f[3][1], 0, 0, 1,
		x - 0.5, y + 0.5, z, f[4][0], f[4][1], 0, 0, 1,
		x - 0.5, y - 0.5, z, f[5][0], f[5][1], 0, 0, 1,
	}...)

	vertices = append(vertices, []float32{
		// back
		x + 0.5, y - 0.5, z, b[0][0], b[0][1], 0, 0, -1,
		x - 0.5, y - 0.5, z, b[1][0], b[1][1], 0, 0, -1,
		x - 0.5, y + 0.5, z, b[2][0], b[2][1], 0, 0, -1,
		x - 0.5, y + 0.5, z, b[3][0], b[3][1], 0, 0, -1,
		x + 0.5, y + 0.5, z, b[4][0], b[4][1], 0, 0, -1,
		x + 0.5, y - 0.5, z, b[5][0], b[5][1], 0, 0, -1,
	}...)
	return vertices
}
