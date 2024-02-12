package benchmark

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/schmizzel/go-graphics/pkg/app"
	"github.com/schmizzel/go-graphics/pkg/bvh"
	"github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
	"github.com/stretchr/testify/require"
)

const (
	WIDTH  = 512
	HEIGHT = 512
	FOV    = 60
)

func BenchmarkBunny(b *testing.B) {
	c := config{
		obj: "../assets/bunny.obj",
		views: []view{
			{
				position: math.NewVector3(0, 1, 2),
				lookAt:   math.NewVector3(-.25, .6, 0),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func BenchmarkDragon(b *testing.B) {
	c := config{
		obj: "../assets/local/dragon.obj",
		views: []view{
			{
				position: math.NewVector3(-.75, .75, -.3),
				lookAt:   math.NewVector3(.1, 0, 0),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func BenchmarkBuddha(b *testing.B) {
	c := config{
		obj: "../assets/local/buddha.obj",
		views: []view{
			{
				position: math.NewVector3(0, 0, -1),
				lookAt:   math.NewVector3(0, 0, 0),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func BenchmarkSponza(b *testing.B) {
	c := config{
		obj: "../assets/local/sponza/sponza.obj",
		views: []view{
			{
				position: math.NewVector3(-5, 3, 0),
				lookAt:   math.NewVector3(5, 3, 0),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func BenchmarkSibenik(b *testing.B) {
	c := config{
		obj: "../assets/local/sibek/sibenik.obj",
		views: []view{
			{
				position: math.NewVector3(-16, -10, 0),
				lookAt:   math.NewVector3(1, -10, 0),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func BenchmarkSanMiguel(b *testing.B) {
	c := config{
		obj: "../assets/local/sanmiguel/san-miguel.obj",
		views: []view{
			{
				position: math.NewVector3(14, 2, 9),
				lookAt:   math.NewVector3(15, 2, 7),
				up:       math.NewVector3(0, 1, 0),
			},
		},
	}

	benchmarkScene(c, b)
}

func benchmarkScene(c config, b *testing.B) {
	buffer := render.NewPixelBuffer(WIDTH, HEIGHT)
	renderer := render.NewDefaultRenderer()
	renderer.Spp = 1

	hq := bvh.NewDefaultPHRBuilder()
	hq.Alpha = 0.55
	hq.Delta = 9

	fast := bvh.NewDefaultPHRBuilder()
	fast.Alpha = 0.5
	fast.Delta = 6

	c.benchLBVH(b, renderer, buffer)
	c.benchPHR(b, "phr-hq", hq, renderer, buffer)
	c.benchPHR(b, "phr-fast", fast, renderer, buffer)
}

type view struct {
	position math.Vector3
	lookAt   math.Vector3
	up       math.Vector3
}

func (v view) toCam() *render.Camera {
	return render.
		NewCamera(WIDTH/HEIGHT, FOV).
		SetUp(v.up.Spread()).
		SetPosition(v.position.Spread()).
		LookAt(v.lookAt.Spread())
}

type config struct {
	obj   string
	views []view
}

func (o config) benchLBVH(b *testing.B, renderer render.Renderer, buff render.Buffer) {
	p, m := prepareScene(b, o.obj)

	var tree *bvh.BVH

	b.Run("lbvh/build", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree = bvh.DefaultLBVH(p, m, runtime.NumCPU())
		}
	})

	o.benchRender(b, "lbvh/render", tree, renderer, buff)
}

func (o config) benchPHR(b *testing.B, name string, builder bvh.PhrBuilder, renderer render.Renderer, buff render.Buffer) {
	p, m := prepareScene(b, o.obj)

	var tree *bvh.BVH
	lbvh := bvh.DefaultLBVH(p, m, runtime.NumCPU())

	b.Run(name+"/build", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tree = builder.Refine(lbvh)
		}
	})

	o.benchRender(b, name+"/render", tree, renderer, buff)
}

func (o config) benchRender(b *testing.B, name string, tree *bvh.BVH, renderer render.Renderer, buff render.Buffer) {
	for i, view := range o.views {
		n := fmt.Sprintf("%s/view%d", name, i)
		cam := view.toCam()
		b.Run(n, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderer.RenderBvh(tree, cam, buff)
			}
		})
	}
}

func prepareScene(b *testing.B, path string) ([]scene.Primitive, []scene.Material) {
	cfg, err := app.ParseConfigFile(path)
	require.NoError(b, err)

	s, err := cfg.ToScene()
	require.NoError(b, err)

	return s.CollectPrimitives()
}
