package bvh

import (
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/schmizzel/go-graphics/pkg/scene"
)

const MAX_CUT_SIZE = 2500

type PhrBuilder struct {
	Alpha           float64 // How quickly cut size will shrink
	Delta           float64 // Determines size of initial cut
	BranchingFactor int
	Threshold       AreaThreshold
	Split           SplitFunction

	jobs           chan phrJob
	threadCount    int
	surface        float64
	initialCutSize int32
}

func NewDefaultPHRBuilder() PhrBuilder {
	return NewPHRBuilder(0.5, 6, 2, runtime.GOMAXPROCS(0))
}

func NewPHRBuilder(alpha float64, delta float64, branchingFactor int, threadCount int) PhrBuilder {
	return PhrBuilder{
		Alpha:           alpha,
		Delta:           delta,
		BranchingFactor: branchingFactor,
		Threshold:       DefaultThreshold,
		Split:           SweepSAH,
		threadCount:     threadCount,
	}
}

// Set the hyper parameters to fast setting
func (b *PhrBuilder) SetFastMode() {
	b.Alpha = 0.5
	b.Delta = 6
}

// Set the hyper parameters to HQ setting
func (b *PhrBuilder) SetHQMode() {
	b.Alpha = 0.55
	b.Delta = 9
}

func (builder *PhrBuilder) BuildFromLBVH(p []scene.Primitive, m []scene.Material) *BVH {
	bvh := LBVH(p, m, scene.EnclosingAABB(p), builder.threadCount)
	return builder.Refine(bvh)
}

func (p *PhrBuilder) Refine(bvh *BVH) *BVH {
	p.surface = bvh.root.aabb.Surface()
	p.initialCutSize = 0

	cut := p.findInitialCut(bvh, p.threadCount)

	wg := sync.WaitGroup{}
	p.jobs = make(chan phrJob, p.threadCount)
	for i := 0; i < p.threadCount; i++ {
		go func() {
			for job := range p.jobs {
				p.buildSubTree(job, &wg)
			}
		}()
	}

	// Temporary branch as a starting point, will be discared afterwards
	temp := newBranch(1)
	temp.aabb = bvh.root.aabb
	wg.Add(1)

	p.jobs <- phrJob{
		cut:        cut,
		parent:     temp,
		childIndex: 0,
	}

	wg.Wait()
	close(p.jobs)

	temp.children[0].parent = nil
	bvh.root = temp.children[0]
	return bvh
}

type phrJob struct {
	cut        phrCut
	parent     *node
	childIndex int
}

func (p *PhrBuilder) buildSubTree(job phrJob, wg *sync.WaitGroup) {
	if len(job.cut.nodes) <= 1 {
		job.parent.addChild(job.cut.nodes[0], job.childIndex)
		wg.Done()
		return
	}
	cuts := make([]phrCut, 1, p.BranchingFactor)
	cuts[0] = job.cut

	// Keep splitting cut until enough nodes to branch the tree are found
	for len(cuts) < p.BranchingFactor {
		// Find the biggest cut
		max := 0
		maxI := 0
		for i, cut := range cuts {
			if len(cut.nodes) > max {
				max = len(cut.nodes)
				maxI = i
			}
		}

		// If the biggest cut has size = 1, no more cuts can be split => break
		if max <= 1 {
			break
		}

		// Split biggest cut
		left, right := p.Split(cuts[maxI])
		if right != nil {
			cuts[maxI] = p.refined(*left)
			cuts = append(cuts, p.refined(*right))
		} else {
			// If cut was not split, make it a leaf node
			cuts[maxI] = phrCut{
				nodes: []*node{makeLeaf(left.bounding, left.nodes...)},
			}
		}
	}

	wg.Add(len(cuts) - 1)

	// Create a new BVH branch
	branch := newBranch(len(cuts))
	branch.parent = job.parent
	branch.aabb = job.cut.bounding
	job.parent.addChild(branch, job.childIndex)

	// Queue all new children to be processed by this or any other thread
	for i, cut := range cuts {
		job := phrJob{
			cut:        cut,
			parent:     branch,
			childIndex: i,
		}

		// If channel is full, directly process job
		select {
		case p.jobs <- job:
		default:
			p.buildSubTree(job, wg)
		}
	}
}

func (p *PhrBuilder) findInitialCut(auxilary *BVH, threadCount int) phrCut {
	queue := make(chan *node, 1024)
	cut := phrCut{
		bounding: auxilary.root.aabb,
		depth:    1,
	}
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < threadCount; i++ {
		go func() {
			for node := range queue {
				p.processNodeInitialCut(node, &wg, &m, queue, &cut.nodes)
			}
		}()
	}
	queue <- auxilary.root
	wg.Add(1)
	wg.Wait()
	close(queue)
	return cut
}

func (p *PhrBuilder) processNodeInitialCut(node *node, wg *sync.WaitGroup, m *sync.Mutex, queue chan *node, cut *[]*node) {
	if node.isLeaf {
		m.Lock()
		*cut = append(*cut, node)
		m.Unlock()
		wg.Done()
	} else {
		if node.aabb.Surface() > p.Threshold(p.surface, p.Alpha, p.Delta, 0) {
			if atomic.AddInt32(&p.initialCutSize, int32(len(node.children)-1)) >= MAX_CUT_SIZE {
				m.Lock()
				*cut = append(*cut, node)
				m.Unlock()
				wg.Done()
				return
			}
			wg.Add(len(node.children) - 1)
			for _, child := range node.children {
				// If channel is full, process node directly to avoid deadlock
				select {
				case queue <- child:
				default:
					p.processNodeInitialCut(child, wg, m, queue, cut)
				}
			}
		} else {
			m.Lock()
			*cut = append(*cut, node)
			m.Unlock()
			wg.Done()
		}
	}
}

func (p *PhrBuilder) refined(cut phrCut) phrCut {
	refinedCut := make([]*node, 0, len(cut.nodes))
	for _, node := range cut.nodes {
		if node.isLeaf {
			refinedCut = append(refinedCut, node)
		} else {
			if node.aabb.Surface() < p.Threshold(p.surface, p.Alpha, p.Delta, cut.depth) {
				// Keep node in cut
				refinedCut = append(refinedCut, node)
			} else {
				// Replace node with children
				refinedCut = append(refinedCut, node.children...)
			}
		}
	}

	if len(refinedCut) == 1 {
		refinedCut[0] = makeLeaf(cut.bounding, refinedCut...)
	}

	return phrCut{
		nodes:    refinedCut,
		bounding: cut.bounding,
		depth:    cut.depth + 1,
	}
}

func makeLeaf(bounding scene.AABB, nodes ...*node) *node {
	leaves := make([]*node, 0)
	for _, n := range nodes {
		n.collectLeaves(&leaves)
	}
	prims := make([]int, 0, len(leaves))
	for _, leaf := range leaves {
		prims = append(prims, leaf.pIds...)
	}
	leaf := newLeaf(prims)
	leaf.aabb = bounding
	return leaf
}

type AreaThreshold func(surface float64, alpha float64, delta float64, depth int) float64

func DefaultThreshold(surface float64, alpha float64, delta float64, depth int) float64 {
	return surface / math.Pow(2, alpha*float64(depth)+float64(delta))
}

type phrCut struct {
	nodes    []*node
	bounding scene.AABB
	depth    int
}

type SplitFunction func(phrCut) (*phrCut, *phrCut)

func SweepSAH(cut phrCut) (l *phrCut, r *phrCut) {
	// Sort along x and y axis using two separate slices
	sort.SliceStable(cut.nodes, func(i, j int) bool {
		return cut.nodes[i].aabb.Barycenter.X < cut.nodes[j].aabb.Barycenter.X
	})
	sorted2 := make([]*node, len(cut.nodes))
	copy(sorted2, cut.nodes)
	sort.SliceStable(sorted2, func(i, j int) bool {
		return sorted2[i].aabb.Barycenter.Y < sorted2[j].aabb.Barycenter.Y
	})

	xSAH := minCost(cut.nodes, cut.bounding, cut.depth)
	ySAH := minCost(sorted2, cut.bounding, cut.depth)
	// Keep the sorted slice with lower cost and override the other by sorting along z axis
	// Finally, return the split with the lowest cost

	if xSAH.cost < ySAH.cost {
		sort.SliceStable(sorted2, func(i, j int) bool {
			return sorted2[i].aabb.Barycenter.Z < sorted2[j].aabb.Barycenter.Z
		})
		zSAH := minCost(sorted2, cut.bounding, cut.depth)
		if xSAH.cost < zSAH.cost {
			return xSAH.left, xSAH.right
		} else {
			return zSAH.left, zSAH.right
		}
	} else {
		sort.SliceStable(cut.nodes, func(i, j int) bool {
			return cut.nodes[i].aabb.Barycenter.Z < cut.nodes[j].aabb.Barycenter.Z
		})
		zSAH := minCost(cut.nodes, cut.bounding, cut.depth)
		if ySAH.cost < zSAH.cost {
			return ySAH.left, ySAH.right
		} else {
			return zSAH.left, zSAH.right
		}
	}
}

type sah struct {
	left  *phrCut
	right *phrCut
	cost  float64
}

// Uses SAH to compute the best split
func minCost(sortedNodes []*node, bounding scene.AABB, depth int) sah {
	// Compute and track right costs by incrementally extending bounding box
	SaRight := sortedNodes[len(sortedNodes)-1].aabb
	rightCosts := make([]float64, len(sortedNodes))
	rightCuts := make([]phrCut, len(sortedNodes))
	nodeCount := 0
	for i := len(sortedNodes) - 1; i > 0; i-- {
		SaRight = SaRight.Add(sortedNodes[i].aabb)
		nodeCount += sortedNodes[i].subtreeSize()
		rightCosts[i] = SaRight.Surface() * float64(nodeCount)
		rightCuts[i] = phrCut{
			nodes:    sortedNodes[i:],
			bounding: SaRight,
			depth:    depth,
		}
	}

	nodeCount += sortedNodes[0].subtreeSize()
	min := sah{
		cost: bounding.Surface() * float64(nodeCount),
		left: &phrCut{
			nodes:    sortedNodes,
			bounding: bounding,
			depth:    depth,
		},
	}

	// Incrementally extend left box and use tracked right costs to compute full SAH cost
	nodeCount = sortedNodes[0].subtreeSize()
	SaLeft := sortedNodes[0].aabb
	for i := 1; i < len(sortedNodes); i++ {
		cost := rightCosts[i] + SaLeft.Surface()*float64(nodeCount)
		if cost < min.cost {
			min.cost = cost
			min.left = &phrCut{
				nodes:    sortedNodes[:i],
				bounding: SaLeft,
				depth:    depth,
			}
			min.right = &rightCuts[i]
		}
		SaLeft = SaLeft.Add(sortedNodes[i].aabb)
		nodeCount += sortedNodes[i].subtreeSize()
	}
	return min
}
