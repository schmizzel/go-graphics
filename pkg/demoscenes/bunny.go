package demoscenes

import (
	"github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type BunnyScene struct {
	objPath string
}

func NewBunnyScene() *BunnyScene {
	return &BunnyScene{objPath: "assets/bunny.obj"}
}

func (b BunnyScene) GoGraphics() (*scene.Node, *render.Camera, error) {
	mesh, err := scene.ParseFromPath(b.objPath)
	if err != nil {
		return nil, nil, err
	}

	material := scene.Diffuse{Albedo: scene.NewColor(.5, .5, .5)}

	s := scene.NewNode()
	bunny := scene.NewNode().SetMesh(mesh).SetMaterial(material)
	light := scene.NewNode().SetMesh(scene.NewSphere(3)).SetMaterial(scene.Light{Emitance: 4, Color: scene.NewColor(1, 1, 1)}).SetPosition(0, 6, 0)

	s.AddChild(bunny).AddChild(light)

	cam := render.NewCamera(1, 60).SetPosition(0, 1, 2).LookAt(-.25, .6, 0)
	return s, cam, nil
}

func (b BunnyScene) Pt() (*pt.Scene, *pt.Camera, error) {
	scene := pt.Scene{}
	material := pt.DiffuseMaterial(pt.HexColor(0x808080))

	mesh, err := pt.LoadOBJ(b.objPath, material)
	if err != nil {
		return nil, nil, err
	}
	scene.Add(mesh)
	scene.Add(pt.NewSphere(pt.V(0, 6, 0), 3, pt.LightMaterial(pt.White, 4)))

	camera := pt.LookAt(pt.V(0, 1, 2), pt.V(-.25, .6, 0), pt.V(0, 1, 0), 60)
	return &scene, &camera, nil
}
