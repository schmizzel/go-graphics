package scene

import (
	"math"
	"math/rand"

	m "github.com/schmizzel/go-graphics/pkg/math"
)

// TODO: Probably better to have material types and move material logic to renderer

type Material interface {
	Scatter(*m.Ray, *Hit, *rand.Rand) (bool, Color)
	EmittedLight() Color
}

type Light struct {
	Color    Color
	Emitance float64
}

func (Light) Scatter(*m.Ray, *Hit, *rand.Rand) (bool, Color) {
	return false, Color{}
}

func (l Light) EmittedLight() Color {
	return l.Color.Scale(l.Emitance)
}

type Diffuse struct {
	Albedo Color
}

func (d Diffuse) Scatter(ray *m.Ray, intersec *Hit, r *rand.Rand) (bool, Color) {
	scatterDirection := intersec.Normal.Add(m.RandomUnitVector(r))

	if scatterDirection.ApproxZero() {
		scatterDirection = intersec.Normal
	}

	ray.Reuse(intersec.Point, scatterDirection)
	return true, d.Albedo
}

func (Diffuse) EmittedLight() Color {
	return NewColor(0, 0, 0)
}

type Reflective struct {
	Albedo    Color
	Diffusion float64 // diffusion in range [0,1]
}

func (d Reflective) Scatter(ray *m.Ray, intersec *Hit, r *rand.Rand) (bool, Color) {
	reflected := reflect(ray.Direction.Unit(), intersec.Normal)
	ray.Reuse(intersec.Point, reflected.Add(m.RandomUnitVector(r).Mul(d.Diffusion)))
	return ray.Direction.Dot(intersec.Normal) > 0, d.Albedo
}

func (Reflective) EmittedLight() Color {
	return NewColor(0, 0, 0)
}

type Refractive struct {
	Albedo Color
	Ratio  float64
}

func (d Refractive) Scatter(ray *m.Ray, intersec *Hit, r *rand.Rand) (bool, Color) {
	refractionRatio := d.Ratio
	if intersec.FrontFace {
		refractionRatio = 1 / d.Ratio
	}

	unitDir := ray.Direction.Unit()
	cos_theta := math.Min(unitDir.Mul(-1).Dot(intersec.Normal), 1.0)
	sin_theta := math.Sqrt(1.0 - cos_theta*cos_theta)

	cannot_refract := refractionRatio*sin_theta > 1.0

	var direction m.Vector3
	if cannot_refract || reflectance(cos_theta, refractionRatio) > rand.Float64() {
		direction = reflect(unitDir, intersec.Normal)
	} else {
		direction = refract(unitDir, intersec.Normal, refractionRatio)
	}
	ray.Reuse(intersec.Point, direction)
	return true, d.Albedo
}

func (Refractive) EmittedLight() Color {
	return NewColor(0, 0, 0)
}

func reflect(v, n m.Vector3) m.Vector3 {
	return v.Sub(n.Mul(v.Dot(n) * 2))
}

func refract(uv, n m.Vector3, etai_over_etat float64) m.Vector3 {
	cos_theta := math.Min(uv.Mul(-1).Dot(n), 1.0)
	rOutPerp := uv.Add(n.Mul(cos_theta)).Mul(etai_over_etat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

// Schlick Approximation
func reflectance(cosine, defractionRatio float64) float64 {
	r0 := (1 - defractionRatio) / (1 + defractionRatio)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
