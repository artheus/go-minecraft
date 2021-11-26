package hud

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/ctx"
	. "github.com/artheus/go-minecraft/math32"
	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/mathgl/mgl32"
)

// LineRenderer is in charge of rendering Lines as HUD on screen
type LineRenderer struct {
	ctx       *ctx.Context
	shader    *glhf.Shader
	cross     *Lines
	wireFrame *Lines
	lastBlock Vec3
}

// NewLineRenderer creates a new instance of LineRenderer
func NewLineRenderer(ctx *ctx.Context) (*LineRenderer, error) {
	r := &LineRenderer{
		ctx: ctx,
	}
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
	width, height := r.ctx.Game().Window().GetFramebufferSize()
	project := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	model := mgl32.Translate3D(float32(width/2), float32(height/2), 0)
	model = model.Mul4(mgl32.Scale3D(float32(height/crossDivider), float32(height/crossDivider), 0))
	r.cross.Render(project.Mul4(model))
}

// renderWireFrame will render a wireframe around blocks currently
// pointed at by player's crosshairs
func (r *LineRenderer) renderWireFrame(mat mgl32.Mat4) {
	var vertices []float32
	b, _ := r.ctx.Game().World().HitTest(r.ctx.Game().Camera().Pos(), r.ctx.Game().Camera().Front())
	if b == nil {
		return
	}

	mat = mat.Mul4(mgl32.Translate3D(float32(b.X), float32(b.Y), float32(b.Z)))
	mat = mat.Mul4(mgl32.Scale3D(1.06, 1.06, 1.06))
	if *b == r.lastBlock {
		r.wireFrame.Render(mat)
		return
	}

	id := *b
	show := block.Sides(
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Left())),
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Right())),
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Up())),
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Down())),
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Front())),
		r.ctx.Game().World().IsTransparent(r.ctx.Game().World().Block(id.Back())),
	)
	vertices = WireFrameData(vertices, show)
	if len(vertices) == 0 {
		return
	}
	r.lastBlock = *b
	if r.wireFrame != nil {
		r.wireFrame.Release()
	}

	r.wireFrame = NewLines(r.shader, vertices)
	r.wireFrame.Render(mat)
}

// Render lines (crosshairs and wireframe) to screen
func (r *LineRenderer) Render() {
	width, height := r.ctx.Game().Window().GetSize()
	projection := mgl32.Perspective(Radian(45), float32(width)/float32(height), 0.01, ChunkWidth)
	camera := r.ctx.Game().Camera().Matrix()
	mat := projection.Mul4(camera)

	r.shader.Begin()
	r.renderCrosshairs()
	r.renderWireFrame(mat)
	r.shader.End()
}

// makeCross creates the HUD crosshairs vao
func makeCross(shader *glhf.Shader) *Lines {
	return NewLines(shader, []float32{
		-0.5, 0, 0, 0.5, 0, 0,
		0, -0.5, 0, 0, 0.5, 0,
	})
}
