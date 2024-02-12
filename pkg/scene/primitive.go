package scene

import m "github.com/schmizzel/go-graphics/pkg/math"

type Hit struct {
	Point     m.Vector3 // intersection Point
	Normal    m.Vector3 // normal at the intersection Point always pointing agains the ray
	FrontFace bool      // Wheter or not the ray hit from the outside or the inside
	T         float64   // distance along the intersection ray
	Primitive Primitive
	Material  Material
}

type Intersectable interface {
	Intersected(ray m.Ray, tMin, tMax float64, hitOut *Hit) bool
	Bounding() AABB
}

type Primitive interface {
	Mesh
	Intersectable
	Transformed(m.Matrix4) Primitive
}
