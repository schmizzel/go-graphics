package render

import (
	m "github.com/schmizzel/go-graphics/pkg/math"
	"github.com/schmizzel/go-graphics/pkg/scene"
	"math"
)

type ClosestHitShader interface {
	Hit(context, *ImageRenderer, m.Ray, *scene.Hit) scene.Color
}

type LitShader struct {
	MaxDepth int
}

func (shader *LitShader) Hit(ctx context, renderer *ImageRenderer, r m.Ray, h *scene.Hit) scene.Color {
	if ctx.depth > shader.MaxDepth {
		return scene.NewColor(0, 0, 0)
	}

	ctx.depth++
	light := h.Material.EmittedLight()
	if b, attenuation := h.Material.Scatter(&r, h, ctx.rand); b {
		if ctx.bvh.ClosestHit(r, 0.0001, math.Inf(1), h) {
			return light.Add(shader.Hit(ctx, renderer, r, h).Blend(attenuation))
		} else {
			return light.Add(renderer.MissShader.Miss(ctx, renderer, r).Blend(attenuation))
		}
	} else {
		return light
	}

}

type MissShader interface {
	Miss(context, *ImageRenderer, m.Ray) scene.Color
}

type DefaultMissShader struct {
	Color scene.Color
}

func (shader *DefaultMissShader) Miss(ctx context, renderer *ImageRenderer, r m.Ray) scene.Color {
	return shader.Color
}

func NewDefaultMissShader() *DefaultMissShader {
	return &DefaultMissShader{
		Color: scene.NewColor(0, 0, 0),
	}
}

func NewSunMissShader() *DefaultMissShader {
	return &DefaultMissShader{
		Color: scene.NewColor(1, .96, .78),
	}
}

type SkyMissShader struct{}

func (shader *SkyMissShader) Miss(ctx context, renderer *ImageRenderer, r m.Ray) scene.Color {
	unit := r.Direction.Unit()
	t := 0.5 * (unit.Y + 1)
	white := scene.NewColor(0.8, 0.8, 0.8)
	blue := scene.NewColor(0.25, 0.35, 0.5)
	return white.Scale(1.0 - t).Add(blue.Scale(t))
}

type HeatmapShader struct {
	Threshold int
}

func (shader *HeatmapShader) Hit(ctx context, renderer *ImageRenderer, r m.Ray, h *scene.Hit) scene.Color {
	return shader.shade(ctx, renderer, r)
}

func (shader *HeatmapShader) Miss(ctx context, renderer *ImageRenderer, r m.Ray) scene.Color {
	return shader.shade(ctx, renderer, r)
}

func (shader *HeatmapShader) shade(ctx context, renderer *ImageRenderer, r m.Ray) scene.Color {
	count := ctx.bvh.TraversalSteps(r, 0.001, math.MaxFloat64)

	if count > shader.Threshold {
		factor := float64(count) / float64(shader.Threshold) * 2
		return scene.NewColor(factor, 0, 0)
	}

	factor := float64(count) / (float64(shader.Threshold) * 1.25)
	return scene.NewColor(0, factor, 1-factor)
}
