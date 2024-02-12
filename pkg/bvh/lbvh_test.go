package bvh_test

import (
	"runtime"
	"testing"

	"github.com/schmizzel/go-graphics/pkg/bvh"
	"github.com/schmizzel/go-graphics/pkg/scene"
	"github.com/stretchr/testify/require"
)

func TestLBVH(t *testing.T) {
	mesh, err := scene.ParseFromPath("../../test/bunny.obj")
	require.NoError(t, err)
	s := scene.NewNode().SetMesh(mesh).SetMaterial(scene.Diffuse{Albedo: scene.NewColor(1, 0, 0)})
	p, m := s.CollectPrimitives()
	tree := bvh.DefaultLBVH(p, m, runtime.NumCPU())
	cost := int(tree.Cost(1, 1))
	require.Equal(t, 39, cost)
}
