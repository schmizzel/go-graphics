package math

import (
	"math"
	"math/rand"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func NewVector3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

// Vector with random X,Y,Z components within [min,max]
func NewRandomVector(min, max float64, r *rand.Rand) Vector3 {
	return Vector3{RandFloat(min, max, r), RandFloat(min, max, r), RandFloat(min, max, r)}
}

func (v1 Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}

func (v1 Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}

func (v1 Vector3) Inverse() Vector3 {
	return Vector3{X: 1 / v1.X, Y: 1 / v1.Y, Z: 1 / v1.Z}
}

func (v1 Vector3) ElemMul(v2 Vector3) Vector3 {
	return Vector3{v1.X * v2.X, v1.Y * v2.Y, v1.Z * v2.Z}
}

func (v1 Vector3) Mul(factor float64) Vector3 {
	return Vector3{X: v1.X * factor, Y: v1.Y * factor, Z: v1.Z * factor}
}

func (v1 Vector3) ElemDiv(v2 Vector3) Vector3 {
	return Vector3{X: v1.X / v2.X, Y: v1.Y / v2.Y, Z: v1.Z / v2.Z}
}

func (v1 Vector3) LengthSquared() float64 {
	return v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z
}

func (v1 Vector3) Length() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}

func (v1 Vector3) Dot(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3) Cross(v2 Vector3) Vector3 {
	return Vector3{X: v1.Y*v2.Z - v1.Z*v2.Y, Y: v1.Z*v2.X - v1.X*v2.Z, Z: v1.X*v2.Y - v1.Y*v2.X}
}

func (v1 Vector3) Unit() Vector3 {
	l := v1.Length()
	return Vector3{v1.X / l, v1.Y / l, v1.Z / l}
}

func (v1 Vector3) Distance(v2 Vector3) float64 {
	a := v1.X - v2.X
	b := v1.Y - v2.Y
	c := v1.Z - v2.Z
	return math.Sqrt(a*a + b*b + c*c)
}

func (v Vector3) ApproxZero() bool {
	return math.Abs(v.X) <= APPROX_THRESH && math.Abs(v.Y) <= APPROX_THRESH && math.Abs(v.Z) <= APPROX_THRESH
}

// Returns homogenous representation of a point as Vector4
func (p Vector3) ToPoint() Vector4 {
	return Vector4{
		x: p.X,
		y: p.Y,
		z: p.Z,
		w: 1,
	}
}

// Returns homogenous representation of a vector as Vector4
func (p Vector3) ToVector() Vector4 {
	return Vector4{
		x: p.X,
		y: p.Y,
		z: p.Z,
		w: 0,
	}
}

func (v Vector3) MinComponent() float64 {
	return math.Min(math.Min(v.X, v.Y), v.Z)
}

func (v Vector3) Negate() Vector3 {
	return Vector3{-v.X, -v.Y, -v.Z}
}

func (p Vector3) Spread() (float64, float64, float64) {
	return p.X, p.Y, p.Z
}

type Vector4 struct {
	x, y, z, w float64
}

func (v Vector4) ToV3() Vector3 {
	return Vector3{
		X: v.x,
		Y: v.y,
		Z: v.z,
	}
}

func (p Vector4) Transformed(m Matrix4) Vector4 {
	return m.Multiply(p)
}

// Calculates a vector with the minimum components of a and b combined
func MinVec(a Vector3, b Vector3) Vector3 {
	x := a.X
	if a.X > b.X {
		x = b.X
	}
	y := a.Y
	if a.Y > b.Y {
		y = b.Y
	}
	z := a.Z
	if a.Z > b.Z {
		z = b.Z
	}
	return NewVector3(x, y, z)
}

// Calculates a vector with the maximum components of a and b combined
func MaxVec(a Vector3, b Vector3) Vector3 {
	x := a.X
	if a.X < b.X {
		x = b.X
	}
	y := a.Y
	if a.Y < b.Y {
		y = b.Y
	}
	z := a.Z
	if a.Z < b.Z {
		z = b.Z
	}
	return NewVector3(x, y, z)
}
