package demoscenes

import (
	"github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type SponzaScene struct {
	objPath string
}

func NewSponzaScene() *SponzaScene {
	return &SponzaScene{objPath: "assets/local/sponza/sponza.obj"}
}

func (b SponzaScene) GoGraphics() (*scene.Node, *render.Camera, error) {
	mesh, err := scene.ParseFromPath(b.objPath)
	if err != nil {
		return nil, nil, err
	}

	material := scene.Diffuse{Albedo: scene.NewHexColor(0x888888)}

	sponza := scene.NewNode().SetMesh(mesh).SetMaterial(material)
	sun := scene.NewNode().SetMesh(scene.NewSphere(4)).SetMaterial(scene.Light{Color: scene.NewHexColor(0xFFD780), Emitance: 30}).SetPosition(-10, 15, 0)

	s := scene.NewNode()
	s.AddChild(sponza)
	s.AddChild(sun)

	cam := render.NewCamera(1, 60).SetPosition(-5, 3, 0).LookAt(5, 3, 0)
	return s, cam, nil
}

func (b SponzaScene) Pt() (*pt.Scene, *pt.Camera, error) {
	scene := pt.Scene{}

	material := pt.DiffuseMaterial(pt.HexColor(0x888888))
	mesh, err := pt.LoadOBJ(b.objPath, material)
	if err != nil {
		return nil, nil, err
	}

	scene.Add(mesh)

	light := pt.LightMaterial(pt.HexColor(0xFFD780), 30)
	scene.Add(pt.NewSphere(pt.Vector{-10, 15, 0}, 4, light))

	camera := pt.LookAt(pt.Vector{-5, 3, 0}, pt.Vector{5, 3, 0}, pt.Vector{0, 1, 0}, 60)

	return &scene, &camera, nil
}
