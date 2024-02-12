package render

import (
	m "github.com/schmizzel/go-graphics/pkg/math"
	"math"
)

type Camera struct {
	orientation    orientation
	viewportWidth  float64
	viewportHeight float64

	lowerLeftCorner m.Vector3
	horizontal      m.Vector3
	vertical        m.Vector3
}

func NewCamera(aspectRatio float64, fov float64) *Camera {
	cam := new(Camera)
	cam.viewportHeight = 2.0 * math.Tan(m.DegreesToRadians(fov)/2)
	cam.viewportWidth = aspectRatio * cam.viewportHeight
	o := newOrientation(m.NewVector3(0, 0, 0), m.NewVector3(0, 0, -1), m.NewVector3(0, 1, 0))
	cam.setOrientation(o)
	return cam
}

func (c *Camera) SetPosition(x, y, z float64) *Camera {
	c.orientation.origin = m.NewVector3(x, y, z)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
	return c
}

func (c *Camera) SetUp(x, y, z float64) *Camera {
	c.orientation.up = m.NewVector3(x, y, z)
	c.setOrientation(c.orientation)
	return c
}

func (c *Camera) LookAt(x, y, z float64) *Camera {
	c.setOrientation(newOrientation(c.orientation.origin, m.NewVector3(x, y, z), c.orientation.up))
	return c
}

func (c *Camera) Translate(x, y, z float64) *Camera {
	c.SetPosition(c.orientation.origin.X+x, c.orientation.origin.Y+y, c.orientation.origin.Z+z)
	return c
}

func (c *Camera) SetFront(v m.Vector3) *Camera {
	c.orientation.w = v.Unit()
	c.orientation.u = c.orientation.up.Cross(c.orientation.w).Unit()
	c.orientation.v = c.orientation.w.Cross(c.orientation.u)
	c.horizontal = c.orientation.u.Mul(c.viewportWidth)
	c.vertical = c.orientation.v.Mul(c.viewportHeight)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
	return c
}

func (c *Camera) Origin() m.Vector3 {
	return c.orientation.origin
}

func (c *Camera) Up() m.Vector3 {
	return c.orientation.up
}

func (c *Camera) W() m.Vector3 {
	return c.orientation.w
}

func (c *Camera) setOrientation(o orientation) *Camera {
	c.orientation = o
	c.horizontal = c.orientation.u.Mul(c.viewportWidth)
	c.vertical = c.orientation.v.Mul(c.viewportHeight)
	c.lowerLeftCorner = c.orientation.origin.Sub(c.horizontal.Mul(0.5)).Sub(c.vertical.Mul(0.5)).Sub(c.orientation.w)
	return c
}

// Cast ray in new direction, while keeping origin the same
func (c *Camera) castRayReuse(s, t float64, ray *m.Ray) {
	ray.ReuseSameOrigin(c.lowerLeftCorner.Add(c.horizontal.Mul(s)).Add(c.vertical.Mul(t)).Sub(c.orientation.origin))
}

type orientation struct {
	origin m.Vector3
	up     m.Vector3
	w      m.Vector3
	u      m.Vector3
	v      m.Vector3
}

func newOrientation(lookFrom, lookAt, up m.Vector3) orientation {
	w := lookFrom.Sub(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)
	return orientation{
		origin: lookFrom,
		w:      lookFrom.Sub(lookAt).Unit(),
		up:     up,
		u:      u,
		v:      v,
	}
}
