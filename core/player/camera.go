package player

import (
	"github.com/artheus/go-minecraft/core/ctx"
	. "github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type CameraMovement int

const (
	MoveForward CameraMovement = iota
	MoveBackward
	MoveLeft
	MoveRight
	MoveJump
)

type Camera struct {
	ctx    *ctx.Context

	pos    mgl32.Vec3
	up     mgl32.Vec3
	right  mgl32.Vec3
	front  mgl32.Vec3
	wfront mgl32.Vec3

	prevtime         float64
	velocityY        float32
	rotateX, rotateY float32

	Sens float32

	flying bool
}

func NewCamera(ctx *ctx.Context, pos mgl32.Vec3) *Camera {
	c := &Camera{
		ctx:     ctx,
		pos:     pos,
		front:   mgl32.Vec3{0, 0, -1},
		rotateY: 0,
		rotateX: -90,
		Sens:    0.14,
		flying:  false,
	}
	c.updateAngles()
	return c
}

func (c *Camera) Restore(state PlayerState) {
	c.pos = mgl32.Vec3{state.X, state.Y, state.Z}
	c.rotateX = state.Rx
	c.rotateY = state.Ry
	c.updateAngles()
}

func (c *Camera) State() PlayerState {
	return PlayerState{
		X:  c.pos.X(),
		Y:  c.pos.Y(),
		Z:  c.pos.Z(),
		Rx: c.rotateX,
		Ry: c.rotateY,
	}
}

func (c *Camera) Matrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.pos, c.pos.Add(c.front), c.up)
}

func (c *Camera) SetPos(pos mgl32.Vec3) {
	c.pos = pos
}

func (c *Camera) Pos() mgl32.Vec3 {
	return c.pos
}

func (c *Camera) Front() mgl32.Vec3 {
	return c.front
}

func (c *Camera) FlipFlying() {
	c.flying = !c.flying
}

func (c *Camera) Flying() bool {
	return c.flying
}

func (c *Camera) OnAngleChange(dx, dy float32) {
	if mgl32.Abs(dx) > 200 || mgl32.Abs(dy) > 200 {
		return
	}
	c.rotateX += dx * c.Sens
	c.rotateY += dy * c.Sens
	if c.rotateY > 89 {
		c.rotateY = 89
	}
	if c.rotateY < -89 {
		c.rotateY = -89
	}
	c.updateAngles()
}

func (c *Camera) EventLoop() {
	running := true
	subscriber := c.ctx.EventPipe().Subscriber()

	go c.GravityLoop()

	for running {
		select {
		case <-c.ctx.Context().Done():
			running = false
			break
		case evt := <-subscriber.Get():
			if move, ok := evt.Object().(*EventMove); ok {
				c.OnMoveChange(move.Move, move.Delta)
			}
		}
	}
}

func (c *Camera) GravityLoop() {
	var running = true

	for running {
		select {
		case <- c.ctx.Context().Done():
			running = false
			break
		default:
			if c.flying {
				continue
			}

			var dt float64

			mainthread.Call(func() {
				now := glfw.GetTime()
				dt = now - c.prevtime
				c.prevtime = now
				if dt > 0.02 {
					dt = 0.02
				}
			})

			c.velocityY -= float32(dt * 20)
			if c.velocityY < -50 {
				c.velocityY = -50
			}

			// TODO: Fix "laggy" walking, due to (kinda) collisions with floor..
			y := c.pos.Y()
			ny := Round(c.pos.Y())
			const pad = 0.25

			head := Vec3{
				X: Round(c.pos.X()),
				Y: ny,
				Z: Round(c.pos.Z()),
			}
			feet := head.Down()

			if c.ctx.Game().World().Block(feet.Down()).Obstacle && y < ny && ny-y > pad && c.velocityY < 0 {
				c.velocityY = 0 //c.pos.Y() - ny - pad
			}

			c.pos = c.pos.Add(mgl32.Vec3{0, c.velocityY*float32(dt), 0})
		}
	}
}

func (c *Camera) OnMoveChange(dir CameraMovement, delta float32) {
	if c.flying {
		delta = 5 * delta
	}

	switch dir {
	case MoveForward:
		if c.flying {
			c.pos = c.pos.Add(c.front.Mul(delta))
		} else {
			c.pos = c.pos.Add(c.wfront.Mul(delta))
		}
	case MoveBackward:
		if c.flying {
			c.pos = c.pos.Sub(c.front.Mul(delta))
		} else {
			c.pos = c.pos.Sub(c.wfront.Mul(delta))
		}
	case MoveLeft:
		c.pos = c.pos.Sub(c.right.Mul(delta))
	case MoveRight:
		c.pos = c.pos.Add(c.right.Mul(delta))
	case MoveJump:
		c.velocityY = delta
	}

	pos := c.Pos()

	pos, _ = c.ctx.Game().World().Collide(pos)
	c.SetPos(pos)
}

func (c *Camera) updateAngles() {
	front := mgl32.Vec3{
		Cos(Radian(c.rotateY)) * Cos(Radian(c.rotateX)),
		Sin(Radian(c.rotateY)),
		Cos(Radian(c.rotateY)) * Sin(Radian(c.rotateX)),
	}
	c.front = front.Normalize()
	c.right = c.front.Cross(mgl32.Vec3{0, 1, 0}).Normalize()
	c.up = c.right.Cross(c.front).Normalize()
	c.wfront = mgl32.Vec3{0, 1, 0}.Cross(c.right).Normalize()
}
