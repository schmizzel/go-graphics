package scene

import m "github.com/schmizzel/go-graphics/pkg/math"

type Triangle struct {
	vertecies [3]Vertex
	box       AABB

	// Cache v0v1 and v0v2 when first computed
	v0v1 m.Vector3
	v0v2 m.Vector3
}

type Vertex struct {
	Position m.Vector3
	Normal   m.Vector3
}

func NewTriangle(vertecies [3]Vertex) *Triangle {
	x := [3]float64{vertecies[0].Position.X, vertecies[1].Position.X, vertecies[2].Position.X}
	y := [3]float64{vertecies[0].Position.Y, vertecies[1].Position.Y, vertecies[2].Position.Y}
	z := [3]float64{vertecies[0].Position.Z, vertecies[1].Position.Z, vertecies[2].Position.Z}
	min := m.NewVector3(m.Min3(x), m.Min3(y), m.Min3(z))
	max := m.NewVector3(m.Max3(x), m.Max3(y), m.Max3(z))
	return &Triangle{
		box:       NewAABB(min, max),
		vertecies: vertecies,
		v0v1:      vertecies[1].Position.Sub(vertecies[0].Position),
		v0v2:      vertecies[2].Position.Sub(vertecies[0].Position),
	}
}

func NewTriangleWithoutNormals(v0 m.Vector3, v1 m.Vector3, v2 m.Vector3) *Triangle {
	vertecies := [3]Vertex{
		{
			Position: v0,
			Normal:   calcNormal(v0, v1, v2),
		},
		{
			Position: v1,
			Normal:   calcNormal(v1, v2, v0),
		},
		{
			Position: v2,
			Normal:   calcNormal(v2, v0, v1),
		},
	}
	return NewTriangle(vertecies)
}

func (t *Triangle) Bounding() AABB {
	return t.box
}

func (tri *Triangle) Transformed(t m.Matrix4) Primitive {
	tinv := t.Transpose().Inverse()

	var vertecies [3]Vertex
	vertecies[0] = Vertex{
		Position: tri.vertecies[0].Position.ToPoint().Transformed(t).ToV3(),
		Normal:   tri.vertecies[0].Normal.ToVector().Transformed(tinv).ToV3(),
	}
	vertecies[1] = Vertex{
		Position: tri.vertecies[1].Position.ToPoint().Transformed(t).ToV3(),
		Normal:   tri.vertecies[1].Normal.ToVector().Transformed(tinv).ToV3(),
	}
	vertecies[2] = Vertex{
		Position: tri.vertecies[2].Position.ToPoint().Transformed(t).ToV3(),
		Normal:   tri.vertecies[2].Normal.ToVector().Transformed(tinv).ToV3(),
	}
	return NewTriangle(vertecies)
}

func (t *Triangle) Primitives() []Primitive {
	return []Primitive{t}
}

func (tri *Triangle) Intersected(ray m.Ray, tMin, tMax float64, hitOut *Hit) bool {
	// Implementation of the Möller-Trumbore algorithm
	pvec := ray.Direction.Cross(tri.v0v2)
	det := tri.v0v1.Dot(pvec)

	// If det is close to 0, Triangle and ray are parallel => no intersection
	if m.ApproxZero(det) {
		return false
	}

	invDet := 1 / det
	tvec := ray.Origin.Sub(tri.vertecies[0].Position)
	u := tvec.Dot(pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	qvec := tvec.Cross(tri.v0v1)
	v := ray.Direction.Dot(qvec) * invDet
	if v < 0 || u+v > 1 {
		return false
	}

	t := tri.v0v2.Dot(qvec) * invDet
	if t < tMin || t > tMax {
		return false
	}

	hitOut.Point = ray.At(t)
	hitOut.FrontFace = det > 0
	hitOut.Normal = tri.normal(u, v)
	if !hitOut.FrontFace {
		hitOut.Normal = hitOut.Normal.Mul(-1)
	}
	hitOut.T = t
	return true
}

// Takes u and v barycentric coordinates and returns the normal at point p
func (tri *Triangle) normal(u, v float64) m.Vector3 {
	normalW := tri.vertecies[0].Normal.Mul(1 - u - v)
	normalU := tri.vertecies[1].Normal.Mul(u)
	normalV := tri.vertecies[2].Normal.Mul(v)
	return normalU.Add(normalV).Add(normalW)
}

func calcNormal(point m.Vector3, right m.Vector3, left m.Vector3) m.Vector3 {
	pa := left.Sub(point)
	pb := right.Sub(point)
	return pb.Cross(pa).Unit()
}
