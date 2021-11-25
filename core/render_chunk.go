package core

import (
	"flag"
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/chunk"
	"github.com/artheus/go-minecraft/core/item"
	mesh2 "github.com/artheus/go-minecraft/core/mesh"
	"github.com/artheus/go-minecraft/core/texture"
	. "github.com/artheus/go-minecraft/math32"
	"github.com/artheus/go-minecraft/types"
	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"sort"
	"sync"
)

var (
	texturePath  = flag.String("t", "texture.png", "texture file")
	renderRadius = flag.Int("r", 6, "render radius")
)

type ChunkRenderer struct {
	shader  *glhf.Shader
	texture *glhf.Texture

	facePool *sync.Pool

	sigch     chan struct{}
	meshcache sync.Map //map[Vec3]*Mesh

	stat Stat

	item *mesh2.Mesh
}

func NewBlockRender() (*ChunkRenderer, error) {
	var (
		err error
	)
	img, rect, err := texture.LoadImage(*texturePath)
	if err != nil {
		return nil, err
	}

	r := &ChunkRenderer{
		sigch: make(chan struct{}, 4),
	}

	mainthread.Call(func() {
		r.shader, err = glhf.NewShader(glhf.AttrFormat{
			glhf.Attr{Name: "pos", Type: glhf.Vec3},
			glhf.Attr{Name: "tex", Type: glhf.Vec2},
			glhf.Attr{Name: "normal", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			glhf.Attr{Name: "matrix", Type: glhf.Mat4},
			glhf.Attr{Name: "camera", Type: glhf.Vec3},
			glhf.Attr{Name: "fogdis", Type: glhf.Float},
		}, blockVertexSource, blockFragmentSource)

		if err != nil {
			return
		}
		r.texture = glhf.NewTexture(rect.Dx(), rect.Dy(), false, img)

	})
	if err != nil {
		return nil, err
	}
	r.facePool = &sync.Pool{
		New: func() interface{} {
			return make([]float32, 0, r.shader.VertexFormat().Size()/4*6*6)
		},
	}

	return r, nil
}

func (r *ChunkRenderer) makeChunkMesh(c *chunk.Chunk, onmainthread bool) *mesh2.Mesh {
	facedata := r.facePool.Get().([]float32)
	defer r.facePool.Put(facedata[:0])

	c.RangeBlocks(func(id Vec3, w int) {
		if w == 0 {
			log.Panicf("unexpect 0 item types on %v", id)
		}
		show := block.Sides(
			IsTransparent(game.world.Block(id.Left())),
			IsTransparent(game.world.Block(id.Right())),
			IsTransparent(game.world.Block(id.Up())),
			IsTransparent(game.world.Block(id.Down())) && id.Y != 0,
			IsTransparent(game.world.Block(id.Front())),
			IsTransparent(game.world.Block(id.Back())),
		)
		if IsPlant(game.world.Block(id)) {
			facedata = block.PlantData(facedata, show, id, item.Tex.Texture(w))
		} else {
			facedata = block.BlockData(facedata, show, id, item.Tex.Texture(w))
		}
	})
	n := len(facedata) / (r.shader.VertexFormat().Size() / 4)
	log.Printf("chunk faces:%d", n/6)
	var mesh *mesh2.Mesh
	if onmainthread {
		mesh = mesh2.NewMesh(r.shader, facedata)
	} else {
		mainthread.Call(func() {
			mesh = mesh2.NewMesh(r.shader, facedata)
		})
	}
	mesh.Id = c.ID()
	return mesh
}

// call on mainthread
func (r *ChunkRenderer) UpdateItem(w int) {
	vertices := r.facePool.Get().([]float32)
	defer r.facePool.Put(vertices[:0])
	texture := item.Tex.Texture(w)
	show := block.Sides(true, true, true, true, true, true)
	pos := Vec3{0, 0, 0}
	if IsPlant(w) {
		vertices = block.PlantData(vertices, show, pos, texture)
	} else {
		vertices = block.BlockData(vertices, show, pos, texture)
	}
	item := mesh2.NewMesh(r.shader, vertices)
	if r.item != nil {
		r.item.Release()
	}
	r.item = item
}

func frustumPlanes(mat *mgl32.Mat4) []mgl32.Vec4 {
	c1, c2, c3, c4 := mat.Rows()
	return []mgl32.Vec4{
		c4.Add(c1),          // left
		c4.Sub(c1),          // right
		c4.Sub(c2),          // top
		c4.Add(c2),          // bottom
		c4.Mul(0.1).Add(c3), // front
		c4.Mul(320).Sub(c3), // back
	}
}

func isChunkVisiable(planes []mgl32.Vec4, id types.ChunkID) bool {
	p := mgl32.Vec3{float32(id.X * ChunkWidth), 0, float32(id.Z * ChunkWidth)}
	const m = ChunkWidth

	points := []mgl32.Vec3{
		{p.X(), p.Y(), p.Z()},
		{p.X() + m, p.Y(), p.Z()},
		{p.X() + m, p.Y(), p.Z() + m},
		{p.X(), p.Y(), p.Z() + m},

		{p.X(), p.Y() + 256, p.Z()},
		{p.X() + m, p.Y() + 256, p.Z()},
		{p.X() + m, p.Y() + 256, p.Z() + m},
		{p.X(), p.Y() + 256, p.Z() + m},
	}
	for _, plane := range planes {
		var in, out int
		for _, point := range points {
			if plane.Dot(point.Vec4(1)) < 0 {
				out++
			} else {
				in++
			}
			if in != 0 && out != 0 {
				break
			}
		}
		if in == 0 {
			return false
		}
	}
	return true
}

func (r *ChunkRenderer) get3dmat() mgl32.Mat4 {
	n := float32(*renderRadius * ChunkWidth)
	width, height := game.win.GetSize()
	mat := mgl32.Perspective(Radian(45), float32(width)/float32(height), 0.01, n)
	mat = mat.Mul4(game.camera.Matrix())
	return mat
}

func (r *ChunkRenderer) get2dmat() mgl32.Mat4 {
	n := float32(*renderRadius * ChunkWidth)
	mat := mgl32.Ortho(-n, n, -n, n, -1, n)
	mat = mat.Mul4(game.camera.Matrix())
	return mat
}

func (r *ChunkRenderer) sortChunks(chunks []types.ChunkID) []types.ChunkID {
	nb := chunk.NearBlock(game.camera.Pos())
	cid := nb.ChunkID()
	x, z := cid.X, cid.Z
	mat := r.get3dmat()
	planes := frustumPlanes(&mat)

	sort.Slice(chunks, func(i, j int) bool {
		v1 := isChunkVisiable(planes, chunks[i])
		v2 := isChunkVisiable(planes, chunks[j])
		if v1 && !v2 {
			return true
		}
		if v2 && !v1 {
			return false
		}
		d1 := (chunks[i].X-x)*(chunks[i].X-x) + (chunks[i].Z-z)*(chunks[i].Z-z)
		d2 := (chunks[j].X-x)*(chunks[j].X-x) + (chunks[j].Z-z)*(chunks[j].Z-z)
		return d1 < d2
	})
	return chunks
}

func (r *ChunkRenderer) updateMeshCache() {
	// Get chunk camera is currently in
	block := chunk.NearBlock(game.camera.Pos())
	chunk := block.ChunkID()
	x, z := chunk.X, chunk.Z

	// Get which chunks to render (camera culling)
	n := *renderRadius
	needed := make(map[types.ChunkID]bool)

	for dx := -n; dx < n; dx++ {
		for dz := -n; dz < n; dz++ {
			id := types.ChunkID{X: x + dx, Z: z + dz}
			if dx*dx+dz*dz > n*n {
				continue
			}
			needed[id] = true
		}
	}

	// Make lists of which blocks are added or removed
	var added, removed []types.ChunkID
	r.meshcache.Range(func(k, v interface{}) bool {
		id := k.(types.ChunkID)
		if !needed[id] {
			removed = append(removed, id)
			return true
		}
		return true
	})

	// Rebuild cache with needed chunks
	for id := range needed {
		mesh, ok := r.meshcache.Load(id)
		// Rebuild those not cached
		if !ok {
			added = append(added, id)
		} else {
			if mesh.(*mesh2.Mesh).Dirty {
				log.Printf("update cache %v", id)
				added = append(added, id)
				removed = append(removed, id)
			}
		}
	}

	// Number of chunks constructed in batch
	const batchBuildChunk = 4
	r.sortChunks(added)
	if len(added) > 4 {
		added = added[:batchBuildChunk]
	}

	// Delete any removed mesh from meshcache
	var removedMesh []*mesh2.Mesh
	for _, id := range removed {
		log.Printf("remove cache %v", id)
		mesh, _ := r.meshcache.Load(id)
		r.meshcache.Delete(id)
		removedMesh = append(removedMesh, mesh.(*mesh2.Mesh))
	}

	newChunks := game.world.Chunks(added)
	for _, c := range newChunks {
		log.Printf("add cache %v", c.ID())
		r.meshcache.Store(c.ID(), r.makeChunkMesh(c, false))
	}

	// Release any removed mesh from VRAM
	mainthread.CallNonBlock(func() {
		for _, mesh := range removedMesh {
			mesh.Release()
		}
	})

}

// forceChunks forces any removed mesh from chunks to be released from VRAM
// must be called on main-thread
func (r *ChunkRenderer) forceChunks(ids []types.ChunkID) {
	var removedMesh []*mesh2.Mesh

	// Get requested chunks
	chunks := game.world.Chunks(ids)

	// Add any removed mesh from requested chunks to removedMesh slice
	for _, chunk := range chunks {
		id := chunk.ID()
		imesh, ok := r.meshcache.Load(id)
		var mesh *mesh2.Mesh
		if ok {
			mesh = imesh.(*mesh2.Mesh)
		}
		if ok && !mesh.Dirty {
			continue
		}
		r.meshcache.Store(id, r.makeChunkMesh(chunk, true))
		if ok {
			removedMesh = append(removedMesh, mesh)
		}
	}

	// Release any removed mesh from VRAM
	mainthread.CallNonBlock(func() {
		for _, mesh := range removedMesh {
			mesh.Release()
		}
	})
}

// forcePlayerChunks runs forceChunks on the chunk the player is currently in
func (r *ChunkRenderer) forcePlayerChunks() {
	bid := chunk.NearBlock(game.camera.Pos())
	cid := bid.ChunkID()

	var ids []types.ChunkID

	for dx := -1; dx <= 1; dx++ {
		for dz := -1; dz <= 1; dz++ {
			id := types.ChunkID{X: cid.X + dx, Z: cid.Z + dz}
			ids = append(ids, id)
		}
	}

	r.forceChunks(ids)
}

// checkChunks sends an empty struct to sigch which in turn
// is caught by UpdateLoop which will run updateMeshCache to
// update meshcache with mesh surrounding the player
func (r *ChunkRenderer) checkChunks() {
	// nonblock signal
	select {
	case r.sigch <- struct{}{}:
	default:
	}
}

// DirtyChunk marks the chunk (by types.ChunkID) as dirty (changed)
func (r *ChunkRenderer) DirtyChunk(id types.ChunkID) {
	mesh, ok := r.meshcache.Load(id)
	if !ok {
		return
	}
	mesh.(*mesh2.Mesh).Dirty = true
}

// UpdateLoop runs a loop for updating meshcache whenever a signal
// is received on the sigch (signal channel)
func (r *ChunkRenderer) UpdateLoop() {
	for {
		select {
		case <-r.sigch:
		}
		r.updateMeshCache()
	}
}

// renderChunks will render all chunks visible to player
// after running forcePlayerChunks to force any changes
// to mesh in player chunk to be updated in VRAM
func (r *ChunkRenderer) renderChunks() {
	r.forcePlayerChunks()
	r.checkChunks()
	mat := r.get3dmat()

	r.shader.SetUniformAttr(0, mat)
	r.shader.SetUniformAttr(1, game.camera.Pos())
	r.shader.SetUniformAttr(2, float32(*renderRadius)*ChunkWidth)

	planes := frustumPlanes(&mat)
	r.stat = Stat{}
	r.meshcache.Range(func(k, v interface{}) bool {
		id, mesh := k.(types.ChunkID), v.(*mesh2.Mesh)
		r.stat.CacheChunks++
		if isChunkVisiable(planes, id) {
			r.stat.RendingChunks++
			r.stat.Faces += mesh.Faces()
			mesh.Draw()
		}
		return true
	})
}

// renderItem will draw the HUD block item, currently selected
func (r *ChunkRenderer) renderItem() {
	if r.item == nil {
		return
	}
	width, height := game.win.GetSize()
	ratio := float32(width) / float32(height)
	projection := mgl32.Ortho2D(0, 15, 0, 15/ratio)
	model := mgl32.Translate3D(1, 1, 0)
	model = model.Mul4(mgl32.HomogRotate3DX(Radian(10)))
	model = model.Mul4(mgl32.HomogRotate3DY(Radian(45)))
	mat := projection.Mul4(model)
	r.shader.SetUniformAttr(0, mat)
	r.shader.SetUniformAttr(1, mgl32.Vec3{0, 0, 0})
	r.shader.SetUniformAttr(2, float32(*renderRadius)*ChunkWidth)
	r.item.Draw()
}

// Render will render all chunks and HUD block items to screen
func (r *ChunkRenderer) Render() {
	r.shader.Begin()
	r.texture.Begin()

	r.renderChunks()
	r.renderItem()

	r.shader.End()
	r.texture.End()
}

type Stat struct {
	Faces         int
	CacheChunks   int
	RendingChunks int
}

func (r *ChunkRenderer) Stat() Stat {
	return r.stat
}
