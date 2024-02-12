package scene

import (
	"math"

	m "github.com/schmizzel/go-graphics/pkg/math"
)

type Sphere struct {
	center m.Vector3
	radius float64
	box    AABB
}

func NewSphere(radius float64) *Sphere {
	s := &Sphere{
		center: m.NewVector3(0, 0, 0),
		radius: radius,
	}

	s.updateAABB()
	return s
}

func (s *Sphere) SetCenter(x, y, z float64) *Sphere {
	s.center = m.NewVector3(x, y, z)
	s.updateAABB()
	return s
}

func (s *Sphere) Transformed(t m.Matrix4) Primitive {
	newCenter := s.center.ToPoint().Transformed(t).ToV3()
	/*
		    // TODO: Fix scaling
				    pointOnSphere := s.center.Add(m.NewVector3(s.radius, 0, 0))
					newRadius := pointOnSphere.ToPoint().Transformed(t).ToV3().Length()
	*/
	return newSphereAt(newCenter.X, newCenter.Y, newCenter.Z, s.radius)
}

func (s *Sphere) Bounding() AABB {
	return s.box
}

func (s *Sphere) Primitives() []Primitive {
	return []Primitive{s}
}

func (s *Sphere) Intersected(ray m.Ray, tMin, tMax float64, hitOut *Hit) bool {
	oc := ray.Origin.Sub(s.center)
	dirNorm := ray.Direction.Length()
	a := dirNorm * dirNorm
	halfB := oc.Dot(ray.Direction)
	ocNorm := oc.Length()
	c := ocNorm*ocNorm - s.radius*s.radius
	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}

	// Nearest intersection distance within tMin <= t <= tMax
	sqrtDiscriminant := math.Sqrt(discriminant)
	t := (-halfB - sqrtDiscriminant) / a
	if t <= tMin || t >= tMax {
		t = (-halfB + sqrtDiscriminant) / a
		if t <= tMin || t >= tMax {
			return false
		}
	}

	hitOut.Point = ray.At(t)
	hitOut.Normal = hitOut.Point.Sub(s.center).Mul(1 / s.radius)
	hitOut.FrontFace = ray.Direction.Dot(hitOut.Normal) < 0
	if !hitOut.FrontFace {
		hitOut.Normal = hitOut.Normal.Mul(-1)
	}
	hitOut.T = t
	return true
}

func newSphereAt(x, y, z, radius float64) *Sphere {
	s := &Sphere{
		center: m.NewVector3(x, y, z),
		radius: radius,
	}

	s.updateAABB()
	return s
}

func (s *Sphere) updateAABB() {
	radVec := m.NewVector3(s.radius, s.radius, s.radius)
	min := s.center.Sub(radVec)
	max := s.center.Add(radVec)
	s.box = NewAABB(min, max)
}
