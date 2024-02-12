package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/apex/log"
	. "github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/bvh"
	"github.com/schmizzel/go-graphics/pkg/demoscenes"
	"github.com/schmizzel/go-graphics/pkg/render"
)

const (
	spp    = 100
	width  = 512
	height = 512

	name = "sponza"
)

var (
	pathA = fmt.Sprintf("renders/pt/%s.png", name)
	pathB = fmt.Sprintf("renders/go-graphics/%s.png", name)
)

func main() {
	// Comment out the scene you want to compare

	// scene := demoscenes.NewBunnyScene()
	// scene := demoscenes.NewDragonScene()
	scene := demoscenes.NewSponzaScene()
	// scene := demoscenes.NewSanMiguelScene()

	now := time.Now()
	log.Info("rendering with pt")
	renderPt(scene, pathA)
	log.Infof("stored image at %s", pathA)
	log.Infof("pt finished in %v", time.Since(now))

	now = time.Now()
	log.Info("rendering with go-graphics")
	renderGoGraphics(scene, pathB)
	log.Infof("stored image at %s", pathB)
	log.Infof("go-graphics finished in %v", time.Since(now))
}

func renderPt(scene demoscenes.DemoScene, path string) {
	s, cam, err := scene.Pt()
	if err != nil {
		log.Errorf("failed to render scene with pt: %s", err.Error())
		return
	}

	sampler := NewSampler(1, 5)
	sampler.DirectLighting = false

	renderer := NewRenderer(s, cam, sampler, width, height)
	renderer.SamplesPerPixel = spp
	renderer.Verbose = false
	renderer.IterativeRender(path, 1)
}

func renderGoGraphics(scene demoscenes.DemoScene, path string) {
	s, cam, err := scene.GoGraphics()
	if err != nil {
		log.Errorf("failed to render scene with go-graphics: %s", err.Error())
		return
	}

	r := render.NewDefaultRenderer()
	r.Spp = spp

	builder := bvh.NewDefaultPHRBuilder()
	builder.SetHQMode()

	p, m := s.CollectPrimitives()
	tree := builder.BuildFromLBVH(p, m)
	buff := render.NewPixelBuffer(width, height)
	r.RenderBvh(tree, cam, buff)

	f, err := os.Create(path)
	if err != nil {
		log.Errorf("failed to create file for go-graphics: %s", err.Error())
		return
	}

	img := buff.ToImage()
	err = png.Encode(f, img)
	if err != nil {
		log.Errorf("failed to encode file for go-graphics: %s", err.Error())
		return
	}
}
