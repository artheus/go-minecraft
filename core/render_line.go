package core

import (
	"github.com/artheus/go-minecraft/core/hud"
	. "github.com/artheus/go-minecraft/math32"
	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/mathgl/mgl32"
)

// LineRenderer is in charge of rendering Lines as HUD on screen
type LineRenderer struct {
	shader    *glhf.Shader
	cross     *hud.Lines
	wireFrame *hud.Lines
	lastBlock Vec3
}

// NewLineRenderer creates a new instance of LineRenderer
func NewLineRenderer() (*LineRenderer, error) {
	r := &LineRenderer{}
	var err error
	mainthread.Call(func() {
		r.shader, err = glhf.NewShader(glhf.AttrFormat{
			glhf.Attr{Name: "pos", Type: glhf.Vec3},
		}, glhf.AttrFormat{
			glhf.Attr{Name: "matrix", Type: glhf.Mat4},
		}, lineVertexSource, lineFragmentSource)

		if err != nil {
			return
		}
		r.cross = makeCross(r.shader)
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

const (
	// crossDivider is a divider for crosshairs size
	// higher value = smaller cross and vise-versa
	crossDivider = 20
)

// renderCrosshairs will render the HUD crosshairs to screen
func (r *LineRenderer) renderCrosshairs() {
	width, height := game.win.GetFramebufferSize()
	project := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	model := mgl32.Translate3D(float32(width/2), float32(height/2), 0)
	model = model.Mul4(mgl32.Scale3D(float32(height/crossDivider), float32(height/crossDivider), 0))
	r.cross.Render(project.Mul4(model))
}

// renderWireFrame will render a wireframe around blocks currently
// pointed at by player's crosshairs
func (r *LineRenderer) renderWireFrame(mat mgl32.Mat4) {
	var vertices []float32
	block, _ := game.world.HitTest(game.camera.Pos(), game.camera.Front())
	if block == nil {
		return
	}

	mat = mat.Mul4(mgl32.Translate3D(float32(block.X), float32(block.Y), float32(block.Z)))
	mat = mat.Mul4(mgl32.Scale3D(1.06, 1.06, 1.06))
	if *block == r.lastBlock {
		r.wireFrame.Render(mat)
		return
	}

	id := *block
	show := ShowSides(
		IsTransparent(game.world.Block(id.Left())),
		IsTransparent(game.world.Block(id.Right())),
		IsTransparent(game.world.Block(id.Up())),
		IsTransparent(game.world.Block(id.Down())),
		IsTransparent(game.world.Block(id.Front())),
		IsTransparent(game.world.Block(id.Back())),
	)
	vertices = WireFrameData(vertices, show)
	if len(vertices) == 0 {
		return
	}
	r.lastBlock = *block
	if r.wireFrame != nil {
		r.wireFrame.Release()
	}

	r.wireFrame = hud.NewLines(r.shader, vertices)
	r.wireFrame.Render(mat)
}

// Render lines (crosshairs and wireframe) to screen
func (r *LineRenderer) Render() {
	width, height := game.win.GetSize()
	projection := mgl32.Perspective(Radian(45), float32(width)/float32(height), 0.01, ChunkWidth)
	camera := game.camera.Matrix()
	mat := projection.Mul4(camera)

	r.shader.Begin()
	r.renderCrosshairs()
	r.renderWireFrame(mat)
	r.shader.End()
}

// makeCross creates the HUD crosshairs vao
func makeCross(shader *glhf.Shader) *hud.Lines {
	return hud.NewLines(shader, []float32{
		-0.5, 0, 0, 0.5, 0, 0,
		0, -0.5, 0, 0, 0.5, 0,
	})
}