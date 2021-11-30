package game

import (
	"fmt"
	"github.com/artheus/go-events"
	evttypes "github.com/artheus/go-events/types"
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/chunk"
	"github.com/artheus/go-minecraft/core/ctx"
	"github.com/artheus/go-minecraft/core/game/rpc"
	"github.com/artheus/go-minecraft/core/game/world"
	"github.com/artheus/go-minecraft/core/hud"
	"github.com/artheus/go-minecraft/core/player"
	"github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"
	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"time"
)

func InitGL(w, h int) *glfw.Window {
	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	win, err := glfw.CreateWindow(w, h, "gocraft", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	win.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		log.Fatal(err)
	}
	glfw.SwapInterval(1) // enable vsync
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

	return win
}

type Application struct {
	Ctx          *ctx.Context
	evtPublisher evttypes.Publisher
	window       *glfw.Window

	camera   *player.Camera
	lx, ly   float64
	vy       float32
	prevtime float64

	lineRenderer   *hud.LineRenderer
	chunkRenderer  *chunk.ChunkRenderer
	playerRenderer *player.PlayerRenderer

	world    *world.World
	itemidx  int
	itemKeys []string
	item     *block.Block
	fps      hud.FPS

	exclusiveMouse bool
	closed         bool
}

func NewGame(w, h int) (game *Application, err error) {
	game = new(Application)

	block.InitRegister()

	block.RangeBlocks(func(b *block.Block) bool {
		if b.ID == block.AirID {
			return true
		}

		game.itemKeys = append(game.itemKeys, b.ID)
		return true
	})

	game.itemidx = 0
	game.item = block.GetBlock(game.itemKeys[game.itemidx])

	mainthread.Call(func() {
		win := InitGL(w, h)
		win.SetMouseButtonCallback(game.onMouseButtonCallback)
		win.SetCursorPosCallback(game.onCursorPosCallback)
		win.SetFramebufferSizeCallback(game.onFrameBufferSizeCallback)
		win.SetKeyCallback(game.onKeyCallback)
		game.window = win
	})

	return
}

func (g *Application) Init(ctx *ctx.Context) (err error) {
	g.Ctx = ctx

	g.evtPublisher = ctx.EventPipe().Publisher()

	g.chunkRenderer, err = chunk.NewChunkRenderer(ctx)
	if err != nil {
		return err
	}

	mainthread.Call(func() {
		g.chunkRenderer.UpdateItem(g.itemKeys[g.itemidx])
	})

	g.lineRenderer, err = hud.NewLineRenderer(ctx)
	if err != nil {
		return err
	}

	g.playerRenderer, err = player.NewPlayerRenderer(ctx)
	if err != nil {
		return err
	}

	go g.chunkRenderer.UpdateLoop()

	g.world = world.NewWorld(ctx)
	g.camera = player.NewCamera(ctx, mgl32.Vec3{0, 16, 0})

	go g.camera.EventLoop()

	go g.syncPlayerLoop()

	return nil
}

func (g *Application) setExclusiveMouse(exclusive bool) {
	if exclusive {
		g.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		g.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
	g.exclusiveMouse = exclusive
}

func (g *Application) Camera() types.ICamera {
	return g.camera
}

func (g *Application) World() types.IWorld {
	return g.world
}

func (g *Application) Window() *glfw.Window {
	return g.window
}

func (g *Application) LineRenderer() types.ILineRenderer {
	return g.lineRenderer
}

func (g *Application) PlayerRenderer() types.IPlayerRenderer {
	return g.playerRenderer
}

func (g *Application) ChunkRenderer() types.IChunkRenderer {
	return g.chunkRenderer
}

func (g *Application) dirtyBlock(id Vec3) {
	cid := id.ChunkID()
	g.chunkRenderer.DirtyChunk(cid)
	neighbors := []Vec3{id.Left(), id.Right(), id.Front(), id.Back()}
	for _, neighbor := range neighbors {
		chunkid := neighbor.ChunkID()
		if chunkid != cid {
			g.chunkRenderer.DirtyChunk(chunkid)
		}
	}
}

func (g *Application) onMouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	if !g.exclusiveMouse {
		g.setExclusiveMouse(true)
		return
	}
	head := chunk.NearBlock(g.camera.Pos())
	foot := head.Down()
	blockInWorld, prev := g.world.HitTest(g.camera.Pos(), g.camera.Front())
	if button == glfw.MouseButton2 && action == glfw.Press {
		if prev != nil && *prev != head && *prev != foot {
			g.world.UpdateBlock(*prev, g.item)
			g.dirtyBlock(*prev)
			go rpc.ClientUpdateBlock(*prev, g.item)
		}
	}
	if button == glfw.MouseButton1 && action == glfw.Press {
		if blockInWorld != nil {
			g.world.UpdateBlock(*blockInWorld, block.GetBlock(block.AirID))
			g.dirtyBlock(*blockInWorld)
			go rpc.ClientUpdateBlock(*blockInWorld, block.GetBlock(block.AirID))
		}
	}
}

func (g *Application) onFrameBufferSizeCallback(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (g *Application) onCursorPosCallback(win *glfw.Window, xpos float64, ypos float64) {
	if !g.exclusiveMouse {
		return
	}
	if g.lx == 0 && g.ly == 0 {
		g.lx, g.ly = xpos, ypos
		return
	}
	dx, dy := xpos-g.lx, g.ly-ypos
	g.lx, g.ly = xpos, ypos
	g.camera.OnAngleChange(float32(dx), float32(dy))
}

func (g *Application) onKeyCallback(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}
	switch key {
	case glfw.KeyTab:
		g.camera.FlipFlying()
	case glfw.KeySpace:
		block := g.CurrentBlockid()
		if g.world.HasBlock(Vec3{block.X, block.Y - 2, block.Z}) {
			_ = g.evtPublisher.Publish(events.Event(time.Now(), &player.EventMove{
				Move: player.MoveJump,
				Delta: 8,
			}))
		}
	case glfw.KeyE:
		g.itemidx = (1 + g.itemidx) % len(g.itemKeys)
		g.item = block.GetBlock(g.itemKeys[g.itemidx])
		g.chunkRenderer.UpdateItem(g.itemKeys[g.itemidx])
	case glfw.KeyR:
		g.itemidx--
		if g.itemidx < 0 {
			g.itemidx = len(g.itemKeys) - 1
		}
		g.item = block.GetBlock(g.itemKeys[g.itemidx])
		g.chunkRenderer.UpdateItem(g.itemKeys[g.itemidx])
	}
}

func (g *Application) handleKeyInput() {
	speed := float32(0.1)
	if g.camera.Flying() {
		speed = 0.2
	}
	if g.window.GetKey(glfw.KeyEscape) == glfw.Press {
		g.setExclusiveMouse(false)
	}
	if g.window.GetKey(glfw.KeyW) == glfw.Press {
		_ = g.evtPublisher.Publish(events.Event(time.Now(), &player.EventMove{
			Move:  player.MoveForward,
			Delta: speed,
		}))
	}
	if g.window.GetKey(glfw.KeyS) == glfw.Press {
		_ = g.evtPublisher.Publish(events.Event(time.Now(), &player.EventMove{
			Move:  player.MoveBackward,
			Delta: speed,
		}))
	}
	if g.window.GetKey(glfw.KeyA) == glfw.Press {
		_ = g.evtPublisher.Publish(events.Event(time.Now(), &player.EventMove{
			Move:  player.MoveLeft,
			Delta: speed,
		}))
	}
	if g.window.GetKey(glfw.KeyD) == glfw.Press {
		_ = g.evtPublisher.Publish(events.Event(time.Now(), &player.EventMove{
			Move:  player.MoveRight,
			Delta: speed,
		}))
	}
}

func (g *Application) CurrentBlockid() Vec3 {
	pos := g.camera.Pos()
	return chunk.NearBlock(pos)
}

func (g *Application) ShouldClose() bool {
	return g.closed
}

func (g *Application) renderStat() {
	g.fps.Update()
	p := g.camera.Pos()
	nb := chunk.NearBlock(p)
	cid := nb.ChunkID()
	stat := g.chunkRenderer.State()
	title := fmt.Sprintf("[%.2f %.2f %.2f] %v [%d/%d %d] %d", p.X(), p.Y(), p.Z(),
		cid, stat.RendingChunks, stat.CacheChunks, stat.Faces, g.fps.Fps())
	g.window.SetTitle(title)
}

func (g *Application) syncPlayerLoop() {
	tick := time.NewTicker(time.Second / 10)
	for range tick.C {
		rpc.ClientUpdatePlayerState(g.Ctx, g.camera.State())
	}
}

func (g *Application) Update() {
	mainthread.Call(func() {
		g.handleKeyInput()

		gl.ClearColor(0.57, 0.71, 0.77, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		g.chunkRenderer.Render()
		g.playerRenderer.Render()
		g.lineRenderer.Render()

		g.renderStat()

		g.window.SwapBuffers()
		glfw.PollEvents()
		g.closed = g.window.ShouldClose()
	})
}
