package scene

import (
	"github.com/schmizzel/go-graphics/pkg/math"
)

type Node struct {
	mesh           Mesh
	material       Material
	children       []*Node
	transformation math.Matrix4
}

func NewNode() *Node {
	return &Node{
		transformation: math.IdentityMatrix(),
	}
}

func (n *Node) SetMesh(mesh Mesh) *Node {
	n.mesh = mesh
	return n
}

func (n *Node) SetMaterial(material Material) *Node {
	n.material = material
	return n
}

func (n *Node) Transform(t math.Matrix4) *Node {
	n.transformation = n.transformation.MultiplyMatrix(t)
	return n
}

func (n *Node) SetPosition(x, y, z float64) *Node {
	n.transformation.Set(3, 0, x)
	n.transformation.Set(3, 1, y)
	n.transformation.Set(3, 2, z)
	return n
}

func (n *Node) Translate(x, y, z float64) *Node {
	m := math.Translate(x, y, z)
	return n.Transform(m)
}

func (n *Node) FitInside(box AABB, anchor math.Vector3) *Node {
	aabb := EnclosingAABB(n.mesh.Primitives())

	scale := box.Size().ElemDiv(aabb.Size()).MinComponent()
	extra := box.Size().Sub(aabb.Size().Mul(scale))

	m := math.IdentityMatrix()
	m = m.Translate(aabb.Bounds[0].Negate().Spread())
	m = m.Scale(scale, scale, scale)
	m = m.Translate(box.Bounds[0].Add(extra.ElemMul(anchor)).Spread())
	n.Transform(m)
	return n
}

func (n *Node) AddChild(child *Node) *Node {
	n.children = append(n.children, child)
	return n
}

func (n *Node) CollectPrimitives() ([]Primitive, []Material) {
	return n.collect(math.IdentityMatrix())
}

func (n *Node) collect(t math.Matrix4) (primitives []Primitive, materials []Material) {
	t = t.MultiplyMatrix(n.transformation)

	if n.mesh != nil {
		for _, p := range n.mesh.Primitives() {
			primitives = append(primitives, p.Transformed(t))
			materials = append(materials, n.material)
		}
	}

	for _, child := range n.children {
		p, m := child.collect(t)
		primitives = append(primitives, p...)
		materials = append(materials, m...)
	}

	return
}
