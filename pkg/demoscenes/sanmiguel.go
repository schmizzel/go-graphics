package demoscenes

import (
	"github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type SanMiguelScene struct {
	objPath string
}

func NewSanMiguelScene() *SanMiguelScene {
	return &SanMiguelScene{objPath: "assets/local/sanmiguel/san-miguel.obj"}
}

func (b SanMiguelScene) GoGraphics() (*scene.Node, *render.Camera, error) {
	mesh, err := scene.ParseFromPath(b.objPath)
	if err != nil {
		return nil, nil, err
	}

	material := scene.Diffuse{Albedo: scene.NewHexColor(0x888888)}

	sponza := scene.NewNode().SetMesh(mesh).SetMaterial(material)
	sun := scene.NewNode().SetMesh(scene.NewSphere(10)).SetMaterial(scene.Light{Color: scene.NewHexColor(0xFFD780), Emitance: 30}).SetPosition(15, 25, 5)

	s := scene.NewNode()
	s.AddChild(sponza)
	s.AddChild(sun)

	cam := render.NewCamera(1, 60).SetPosition(14, 2, 9).LookAt(15, 2, 7)
	// cam := render.NewCamera(1, 60).SetPosition(26,7,-2).LookAt(5,7,-2)
	return s, cam, nil
}

func (b SanMiguelScene) Pt() (*pt.Scene, *pt.Camera, error) {
	scene := pt.Scene{}

	material := pt.DiffuseMaterial(pt.HexColor(0x888888))
	mesh, err := pt.LoadOBJ(b.objPath, material)
	if err != nil {
		return nil, nil, err
	}

	scene.Add(mesh)

	light := pt.LightMaterial(pt.HexColor(0xFFD780), 30)
	scene.Add(pt.NewSphere(pt.Vector{-15, 25, 5}, 10, light))

	camera := pt.LookAt(pt.Vector{14, 2, 9}, pt.Vector{15, 2, 7}, pt.Vector{0, 1, 0}, 60)

	return &scene, &camera, nil
}
