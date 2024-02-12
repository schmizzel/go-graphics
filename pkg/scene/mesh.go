package scene

type Mesh interface {
	Primitives() []Primitive
}

type TriangleMesh struct {
	triangles []*Triangle
}

func NewTriangleMesh(triangles []*Triangle) *TriangleMesh {
	return &TriangleMesh{triangles}
}

func (m *TriangleMesh) Primitives() []Primitive {
	primitives := make([]Primitive, len(m.triangles))

	for i, tri := range m.triangles {
		primitives[i] = tri
	}

	return primitives
}
