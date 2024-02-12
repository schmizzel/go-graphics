package app

import (
	"fmt"
	"image/png"
	"os"

	"github.com/schmizzel/go-graphics/pkg/bvh"
	"github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
	s "github.com/schmizzel/go-graphics/pkg/scene"
)

func SaveImage(cfg Config) error {
	buffer, err := cfg.render()
	if err != nil {
		return fmt.Errorf("failed to render image: %w", err)
	}

	err = saveImage(buffer, cfg.Out)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

func (cfg *Config) render() (*render.PixelBuffer, error) {
	scene, err := cfg.toScene()
	if err != nil {
		return nil, fmt.Errorf("failed to build scene: %w", err)
	}

	tree := cfg.Process.buildBvh(scene.CollectPrimitives())

	ar := float64(cfg.Image.Width) / float64(cfg.Image.Height)
	buffer := render.NewPixelBuffer(cfg.Image.Width, cfg.Image.Height)
	cam := cfg.Scene.Camera.toCamera(ar)

	renderer := cfg.Process.toRenderer()
	renderer.RenderBvh(tree, cam, buffer)
	return buffer, nil
}

func (c Camera) toCamera(ar float64) *render.Camera {
	return render.
		NewCamera(ar, c.Fov).
		SetPosition(c.LookFrom[0], c.LookFrom[1], c.LookFrom[2]).
		SetUp(c.Up[0], c.Up[1], c.Up[2]).
		LookAt(c.LookAt[0], c.LookAt[1], c.LookAt[2])
}

func (p Process) toRenderer() *render.ImageRenderer {
	if p.Heatmap {
		return render.NewHeatmapRenderer(p.HeatmapThreshold)
	}

	r := render.NewDefaultRenderer()
	r.Spp = p.Spp
	r.NumCPU = p.Threads
	r.MissShader = &render.SkyMissShader{}
	return r
}

func (process Process) buildBvh(p []s.Primitive, m []s.Material) *bvh.BVH {
	if process.UsePhr {
		builder := bvh.NewPHRBuilder(process.Alpha, process.Delta, 2, process.Threads)
		return builder.BuildFromLBVH(p, m)
	}

	return bvh.DefaultLBVH(p, m, process.Threads)
}

func (cfg *Config) toScene() (*s.Node, error) {
	scene := s.NewNode()

	for _, o := range cfg.Scene.Objects {
		obj, err := s.ParseFromPath(o.File)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", cfg.In, err)
		}

		scale := math.Scale(o.Scale[0], o.Scale[1], o.Scale[2])
		translate := math.Translate(o.Position[0], o.Position[1], o.Position[2])
		t := math.IdentityMatrix().MultiplyMatrix(scale).MultiplyMatrix(translate)

		node := s.NewNode().SetMesh(obj).SetMaterial(o.Material.toMaterial()).Transform(t)
		scene.AddChild(node)
	}

	return scene, nil
}

// TODO: Move the material type to scene package?
func (m Material) toMaterial() s.Material {
	switch m.Type {
	case "diffuse":
		return s.Diffuse{Albedo: scene.NewColor(m.Albedo[0], m.Albedo[1], m.Albedo[2])}
	case "reflective":
		return s.Reflective{Albedo: scene.NewColor(m.Albedo[0], m.Albedo[1], m.Albedo[2]), Diffusion: m.Diffustion}
	case "refractive":
		return s.Refractive{Albedo: scene.NewColor(m.Albedo[0], m.Albedo[1], m.Albedo[2]), Ratio: m.Ratio}
	case "light":
		return s.Light{Color: scene.NewColor(m.Albedo[0], m.Albedo[1], m.Albedo[2])}
	default:
		return s.Diffuse{Albedo: scene.NewColor(.5, .5, .5)}
	}
}

func saveImage(buff *render.PixelBuffer, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	img := buff.ToImage()
	return png.Encode(f, img)
}
