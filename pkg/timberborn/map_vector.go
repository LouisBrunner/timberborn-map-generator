package timberborn

type Vector2 struct {
	X int
	Y int
}

func NewVector2(x, y int) Vector2 {
	return Vector2{
		X: x,
		Y: y,
	}
}

type Vector3 struct {
	Vector2
	Z int
}

func NewVector3(x, y, z int) Vector3 {
	return Vector3{
		Vector2: NewVector2(x, y),
		Z:       z,
	}
}
