package scene

import m "github.com/schmizzel/go-graphics/pkg/math"

// TODO: Move AABB and all intersection code to seperate module?

type AABB struct {
	Bounds     [2]m.Vector3 // bounds[0] = min, bounds[1] = max
	Width      float64
	Height     float64
	Depth      float64
	Barycenter m.Vector3
}

func NewAABB(min, max m.Vector3) AABB {
	bouding := AABB{
		Bounds: [2]m.Vector3{min, max},
	}
	bouding.Update()
	return bouding
}

func EnclosingAABB(primitives []Primitive) AABB {
	enclosing := primitives[0].Bounding()
	for i := 1; i < len(primitives); i++ {
		enclosing = enclosing.Add(primitives[i].Bounding())
	}
	return enclosing
}

func (a *AABB) Update() {
	min := a.Bounds[0]
	max := a.Bounds[1]
	a.Barycenter = min.Add(max).Mul(1.0 / 2.0)
	a.Width = max.X - min.X
	a.Height = max.Y - min.Y
	a.Depth = max.Z - min.Z
}

func (a *AABB) Surface() float64 {
	return 2*a.Width*a.Height + 2*a.Width*a.Depth + 2*a.Height*a.Depth
}

func (a AABB) Add(b AABB) AABB {
	return NewAABB(m.MinVec(a.Bounds[0], b.Bounds[0]), m.MaxVec(a.Bounds[1], b.Bounds[1]))
}

func (a AABB) Size() m.Vector3 {
	return a.Bounds[1].Sub(a.Bounds[0])
}

func (AABB AABB) Intersected(ray m.Ray, tMin, tMax float64) bool {
	tXmin := (AABB.Bounds[ray.Sign[0]].X - ray.Origin.X) * ray.InvDirection.X
	tXmax := (AABB.Bounds[1-ray.Sign[0]].X - ray.Origin.X) * ray.InvDirection.X
	tYmin := (AABB.Bounds[ray.Sign[1]].Y - ray.Origin.Y) * ray.InvDirection.Y
	tYmax := (AABB.Bounds[1-ray.Sign[1]].Y - ray.Origin.Y) * ray.InvDirection.Y

	if tXmin > tYmax || tYmin > tXmax {
		return false
	}
	if tYmin > tXmin {
		tXmin = tYmin
	}
	if tYmax < tXmax {
		tXmax = tYmax
	}

	tZmin := (AABB.Bounds[ray.Sign[2]].Z - ray.Origin.Z) * ray.InvDirection.Z
	tZmax := (AABB.Bounds[1-ray.Sign[2]].Z - ray.Origin.Z) * ray.InvDirection.Z

	if tXmin > tZmax || tZmin > tXmax {
		return false
	}

	// Check if the intersection lies outside tMin and tMax
	if tZmin > tXmin {
		tXmin = tZmin
	}

	if tZmax < tXmax {
		tXmax = tZmax
	}

	if tXmax < tMin || tXmin > tMax {
		return false
	}

	return true
}
