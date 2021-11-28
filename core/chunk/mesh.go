package chunk

type Mesh struct {
	vao     uint32
	vbos    map[uint8]uint32
	faces   int
	visible bool
}
