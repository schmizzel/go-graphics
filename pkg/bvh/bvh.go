package bvh

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/schmizzel/go-graphics/pkg/internal/printer"
	"github.com/schmizzel/go-graphics/pkg/internal/stack"
	m "github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type BVH struct {
	root *node

	// TODO: Abscract primivive to arbitrary input?
	// TODO: Are prim mat pairs more efficient?
	primitives []scene.Primitive
	materials  []scene.Material
}

func (bvh *BVH) ClosestHit(ray m.Ray, tMin, tMax float64, hitOut *scene.Hit) (ok bool) {
	hitOut.T = tMax
	return bvh.root.ClosestHit(ray, bvh, tMin, tMax, hitOut)
}

func (bvh *BVH) ToString() string {
	return printer.PrintTree(bvh.root)
}

func (bvh *BVH) TraversalSteps(ray m.Ray, tMin, tMax float64) int {
	stack := stack.New(bvh.root)
	hit := scene.Hit{T: tMax}
	count := 0

	for {
		count++
		node, hasValue := stack.Pop()
		if !hasValue {
			return count
		}

		if !node.aabb.Intersected(ray, tMin, hit.T) {
			continue
		}

		if node.isLeaf {
			for _, pId := range node.pIds {
				count++
				primitive := bvh.primitives[pId]
				primitive.Intersected(ray, tMin, hit.T, &hit)
				continue
			}
		}

		if len(node.children) == 2 {
			distA := node.children[0].aabb.Barycenter.Distance(ray.Origin)
			distB := node.children[1].aabb.Barycenter.Distance(ray.Origin)

			if distA < distB {
				stack.Push(node.children[1])
				stack.Push(node.children[0])
			} else {
				stack.Push(node.children[0])
				stack.Push(node.children[1])
			}
		} else {
			// TODO: Also sort by distance with higher number of children?
			for _, child := range node.children {
				stack.Push(child)
			}
		}
	}
}

func (bvh *BVH) Cost(ct, ci float64) float64 {
	if bvh.root == nil {
		return 0
	}

	return bvh.root.costSAH(ct, ci)
}

// TODO: There is probably a better way to do this
func (bvh *BVH) updateBounding(threads int) {
	if bvh.root.isLeaf {
		bvh.root.aabb = enclosingSlice(bvh.root.pIds, bvh.primitives)
		return
	}

	// TODO: Cache Leaves or add on construction?
	leaves := make([]*node, 0, len(bvh.primitives))
	bvh.root.collectLeaves(&leaves)

	wg := sync.WaitGroup{}
	wg.Add(threads)
	jobs := make(chan *node)
	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for leaf := range jobs {
				leaf.updateAABB(bvh.primitives)
			}
		}()
	}

	for _, leaf := range leaves {
		jobs <- leaf
	}

	close(jobs)
	wg.Wait()
}

type primitiveId = int

// TODO: Test if polymorph node is slower (Node interface)
type node struct {
	aabb     scene.AABB
	parent   *node
	children []*node
	isLeaf   bool
	// TODO: Test is it is faster to store pointers to primitives?
	pIds []primitiveId

	childAABBset uint32 // Couter used for updating AABB
	size         int
}

func newBranch(size int) *node {
	return &node{
		isLeaf:   false,
		children: make([]*node, size),
	}
}

func newLeaf(pIds []primitiveId) *node {
	return &node{
		isLeaf: true,
		pIds:   pIds,
	}
}

func (n *node) ClosestHit(ray m.Ray, bvh *BVH, tMin, tMax float64, hitOut *scene.Hit) (ok bool) {
	if !n.aabb.Intersected(ray, tMin, hitOut.T) {
		return false
	}

	if n.isLeaf {
		for i := 0; i < len(n.pIds); i++ {
			pId := n.pIds[i]
			primitive := bvh.primitives[pId]
			h := primitive.Intersected(ray, tMin, hitOut.T, hitOut)
			if h {
				ok = true
				hitOut.Material = bvh.materials[pId]
			}
		}

		return
	}

	distA := n.children[0].aabb.Barycenter.Distance(ray.Origin)
	distB := n.children[1].aabb.Barycenter.Distance(ray.Origin)

	if distA < distB {
		a := n.children[0].ClosestHit(ray, bvh, tMin, hitOut.T, hitOut)
		b := n.children[1].ClosestHit(ray, bvh, tMin, hitOut.T, hitOut)
		return a || b
	}

	a := n.children[1].ClosestHit(ray, bvh, tMin, hitOut.T, hitOut)
	b := n.children[0].ClosestHit(ray, bvh, tMin, hitOut.T, hitOut)
	return a || b
}

func (n *node) addChild(node *node, index int) {
	n.children[index] = node
	node.parent = n
}

func (n *node) GetName() string {
	if n.isLeaf {
		return fmt.Sprintf("%d primitives", len(n.pIds))
	}
	return "node"
}

func (n *node) GetChildren() []printer.Node {
	children := make([]printer.Node, len(n.children))
	for i, child := range n.children {
		children[i] = child
	}

	return children
}

func (node *node) collectLeaves(acc *[]*node) {
	if node.isLeaf {
		*acc = append(*acc, node)
		return
	}
	for _, child := range node.children {
		child.collectLeaves(acc)
	}
}

func (node *node) subtreeSize() int {
	if node.size == 0 {
		node.size = 1
		for _, child := range node.children {
			node.size += child.subtreeSize()
		}
	}

	return node.size
}

func enclosingSlice(indeces []int, primitives []scene.Primitive) scene.AABB {
	enclosing := primitives[indeces[0]].Bounding()
	for i := 1; i < len(indeces); i++ {
		prim := primitives[indeces[i]]
		enclosing = enclosing.Add(prim.Bounding())
	}
	return enclosing
}

// This fails when called on a leaf without parent, i.e. a bvh with a single node.
func (node *node) updateAABB(primitives []scene.Primitive) {
	if node.isLeaf {
		node.aabb = enclosingSlice(node.pIds, primitives)
		// Atomic counter. after all child bounding boxes have been computed the parents bounding box can be calculated
		if atomic.AddUint32(&node.parent.childAABBset, 1)%uint32(len(node.parent.children)) == 0 {
			node.parent.updateAABB(primitives)
		}
		return
	}

	node.aabb = node.children[0].aabb
	for i := 1; i < len(node.children); i++ {
		node.aabb = node.aabb.Add(node.children[i].aabb)
	}
	node.aabb.Update()

	if node.parent == nil {
		return
	}
	if atomic.AddUint32(&node.parent.childAABBset, 1)%uint32(len(node.parent.children)) == 0 {
		node.parent.updateAABB(primitives)
	}
}

func (node *node) costSAH(ct, ci float64) float64 {
	if node.isLeaf {
		return ci * float64(len(node.pIds))
	}

	cost := 0.0
	for _, child := range node.children {
		p := float64(child.aabb.Surface()) / float64(node.aabb.Surface())
		cost += p * child.costSAH(ci, ct)
	}

	return ct + cost
}
