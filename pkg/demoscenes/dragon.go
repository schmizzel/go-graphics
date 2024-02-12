package demoscenes

import (
	"github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type DragonScene struct {
	objPath string
}

func NewDragonScene() *DragonScene {
	return &DragonScene{objPath: "assets/local/dragon.obj"}
}

func (b DragonScene) GoGraphics() (*scene.Node, *render.Camera, error) {
	mesh, err := scene.ParseFromPath(b.objPath)
	if err != nil {
		return nil, nil, err
	}

	material := scene.Diffuse{Albedo: scene.NewHexColor(0xB7CA79)}

	s := scene.NewNode()
	dragon := scene.NewNode().SetMesh(mesh).SetMaterial(material)
	dragon.FitInside(scene.NewAABB(math.NewVector3(-1, 0, -1), math.NewVector3(1, 2, 1)), math.NewVector3(0.5, 0, 0.5))
	s.AddChild(dragon)

	light := scene.NewNode().SetMesh(scene.NewSphere(3)).SetMaterial(scene.Light{Color: scene.NewColor(1, 1, 1), Emitance: 150}).SetPosition(-1, 10, 0)
	s.AddChild(light)

	mouth := scene.NewNode().SetMesh(scene.NewSphere(0.03)).SetMaterial(scene.Light{Color: scene.NewHexColor(0xFFFAD5), Emitance: 1000}).SetPosition(-0.05, 1, -0.5)
	s.AddChild(mouth)

	cam := render.NewCamera(1, 35).SetPosition(-3, 2, -1).LookAt(0, .6, -.1)
	return s, cam, nil
}

func (b DragonScene) Pt() (*pt.Scene, *pt.Camera, error) {
	scene := pt.Scene{}

	material := pt.DiffuseMaterial(pt.HexColor(0xB7CA79))
	mesh, err := pt.LoadOBJ(b.objPath, material)
	if err != nil {
		panic(err)
	}
	mesh.FitInside(pt.Box{pt.Vector{-1, 0, -1}, pt.Vector{1, 2, 1}}, pt.Vector{0.5, 0, 0.5})
	scene.Add(mesh)

	light := pt.LightMaterial(pt.White, 75)
	scene.Add(pt.NewSphere(pt.Vector{-1, 10, 0}, 3, light))

	mouth := pt.LightMaterial(pt.HexColor(0xFFFAD5), 500)
	scene.Add(pt.NewSphere(pt.V(-0.05, 1, -0.5), 0.03, mouth))

	camera := pt.LookAt(pt.Vector{-3, 2, -1}, pt.Vector{0, 0.6, -0.1}, pt.Vector{0, 1, 0}, 35)

	return &scene, &camera, nil
}
