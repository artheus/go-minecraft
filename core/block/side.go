package block

type Side struct {
	Left, Right,
	Up, Down,
	Front, Back bool
}

func Sides(left, right, up, down, front, back bool) Side {
	return Side{left, right, up, down, front, back}
}
