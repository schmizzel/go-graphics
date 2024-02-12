package scene

import (
	"image/color"
	"math"

	m "github.com/schmizzel/go-graphics/pkg/math"
)

type Color m.Vector3

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

func NewHexColor(x int) Color {
	r := float64((x>>16)&0xff) / 255
	g := float64((x>>8)&0xff) / 255
	b := float64((x>>0)&0xff) / 255
	return Color{r, g, b}.Pow(2.2)
}

func NewColor255(r, g, b int) Color {
	return Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
}

// Convert to image color and also Gamma correct for gamma=2.0
func (c Color) GoColor() color.Color {
	r := math.Sqrt(c.X)
	g := math.Sqrt(c.Y)
	b := math.Sqrt(c.Z)
	return color.RGBA{R: uint8(m.Clamp(r, 0.0, 1.0) * 255), G: uint8(m.Clamp(g, 0.0, 1.0) * 255), B: uint8(m.Clamp(b, 0.0, 1.0) * 255), A: 255}
}

func (c1 Color) Blend(c2 Color) Color {
	return Color{c1.X * c2.X, c1.Y * c2.Y, c1.Z * c2.Z}
}

func (c1 Color) Add(c2 Color) Color {
	return Color{c1.X + c2.X, c1.Y + c2.Y, c1.Z + c2.Z}
}

func (c1 Color) Sub(c2 Color) Color {
	return Color{c1.X - c2.X, c1.Y - c2.Y, c1.Z - c2.Z}
}

func (c Color) Div(n float64) Color {
	return Color{c.X / n, c.Y / n, c.Z / n}
}

func (c Color) Scale(n float64) Color {
	return Color{c.X * n, c.Y * n, c.Z * n}
}

func (a Color) Pow(b float64) Color {
	return Color{math.Pow(a.X, b), math.Pow(a.Y, b), math.Pow(a.Z, b)}
}
