package demoscenes

import (
	"github.com/fogleman/pt/pt"
	"github.com/schmizzel/go-graphics/pkg/render"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type DemoScene interface {
	GoGraphics() (*scene.Node, *render.Camera, error)
	Pt() (*pt.Scene, *pt.Camera, error)
}
