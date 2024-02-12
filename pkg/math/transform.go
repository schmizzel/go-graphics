package math

import "math"

func ScaleUniform(factor float64) Matrix4 {
	return Scale(factor, factor, factor)
}

func Scale(x, y, z float64) Matrix4 {
	return Matrix4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

func Translate(x, y, z float64) Matrix4 {
	return Matrix4{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}

func Rotate(dir Vector3, angle float64) Matrix4 {
	u := dir.Unit()
	cosa := math.Cos(angle / 2)
	sina := math.Sin(angle / 2)
	q := NewQuanternion(cosa, NewVector3(sina*u.X, sina*u.Y, sina*u.Z))
	return q.ToRotationMatrix()
}
