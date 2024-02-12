package math

type Ray struct {
	Origin         Vector3
	Direction      Vector3
	DirNormSquared float64
	InvDirection   Vector3
	Sign           [3]int
}

func NewRay(origin Vector3, direction Vector3) Ray {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	return Ray{
		Origin:         origin,
		Direction:      direction,
		InvDirection:   invDirection,
		DirNormSquared: dirNormSq,
		Sign:           sign,
	}
}

func (r Ray) At(t float64) Vector3 {
	return r.Origin.Add(r.Direction.Mul(t))
}

// Creates a new ray by overriding the already allocated ray
func (r *Ray) Reuse(origin Vector3, direction Vector3) {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	r.Origin = origin
	r.Direction = direction
	r.InvDirection = invDirection
	r.DirNormSquared = dirNormSq
	r.Sign = sign
}

// Creates a new ray by overriding the already allocated ray
func (r *Ray) ReuseSameOrigin(direction Vector3) {
	invDirection := direction.Inverse()
	sign := [3]int{}

	if invDirection.X < 0 {
		sign[0] = 1
	}
	if invDirection.Y < 0 {
		sign[1] = 1
	}
	if invDirection.Z < 0 {
		sign[2] = 1
	}

	dirNormSq := direction.LengthSquared()
	r.Direction = direction
	r.InvDirection = invDirection
	r.DirNormSquared = dirNormSq
	r.Sign = sign
}
