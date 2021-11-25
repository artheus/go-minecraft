package thread

type InitFunc func() (err error)
type TickFunc func()

type IThread interface {
	Init() InitFunc
	Tick() TickFunc
}
