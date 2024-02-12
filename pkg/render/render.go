package render

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/schmizzel/go-graphics/pkg/bvh"
	m "github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type Renderer interface {
	RenderBvh(*bvh.BVH, *Camera, Buffer)
}

type ImageRenderer struct {
	NumCPU int
	Spp    int

	ClosestHitShader ClosestHitShader
	MissShader       MissShader

	Sampling Sampling
}

func NewHeatmapRenderer(threshold int) *ImageRenderer {
	return &ImageRenderer{
		NumCPU:           runtime.GOMAXPROCS(0),
		Spp:              1,
		ClosestHitShader: &HeatmapShader{Threshold: threshold},
		MissShader:       &HeatmapShader{Threshold: threshold},
		Sampling:         RandomSampling,
	}
}

func NewDefaultRenderer() *ImageRenderer {
	return &ImageRenderer{
		NumCPU:           runtime.GOMAXPROCS(0),
		Spp:              300,
		ClosestHitShader: &LitShader{MaxDepth: 5},
		MissShader:       NewDefaultMissShader(),
		Sampling:         RandomSampling,
	}
}

type context struct {
	bvh   *bvh.BVH
	rand  *rand.Rand
	depth int
}

func (r *ImageRenderer) RenderBvh(b *bvh.BVH, cam *Camera, buff Buffer) {
	jobs := make(chan int, buff.Height())
	wg := sync.WaitGroup{}
	wg.Add(r.NumCPU)
	width := buff.Width()
	height := buff.Height()

	for i := 0; i < r.NumCPU; i++ {
		go func(ctx context, w, h int) {
			ray := m.Ray{
				Origin: cam.orientation.origin,
			}
			hit := scene.Hit{}
			for y := range jobs {
				for x := 0; x < w; x++ {
					for i := 0; i < r.Spp; i++ {
						u, v := r.Sampling(ctx, x, y, w, h)
						cam.castRayReuse(u, v, &ray)
						if b.ClosestHit(ray, 0.001, math.Inf(1), &hit) {
							buff.AddSample(x, y, r.ClosestHitShader.Hit(ctx, r, ray, &hit))
						} else {
							buff.AddSample(x, y, r.MissShader.Miss(ctx, r, ray))
						}
					}
				}
			}
			wg.Done()
		}(context{
			rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
			bvh:   b,
			depth: 0,
		}, width, height)
	}

	for y := 0; y < height; y++ {
		jobs <- y
	}

	close(jobs)
	wg.Wait()
}

type Sampling func(ctx context, x, y, w, h int) (u, v float64)

func RandomSampling(ctx context, x, y, w, h int) (u, v float64) {
	u = (float64(x) + ctx.rand.Float64()) / float64(w-1)
	v = (float64(y) + ctx.rand.Float64()) / float64(h-1)
	return
}
