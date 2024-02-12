package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/schmizzel/go-graphics/pkg/scene"
)

type Config struct {
	In  string `short:"f" long:"in" description:"Input file. Either a .obj or config json file" default:"config.json"`
	Out string `short:"o" long:"out" description:"Output file" default:"out.png"`

	Scene   Scene `json:"scene"`
	Image   Image `json:"image"`
	Process Process
}

func NewDefaultConfig() Config {
	return Config{
		Out:     "out.png",
		Scene:   NewDefaultScene(),
		Image:   NewDefaultImage(),
		Process: NewDefaultProcess(),
	}
}

// Parse config from command line flags
func ParseConfig() (Config, error) {
	cfg := NewDefaultConfig()
	_, err := flags.NewParser(&cfg, flags.HelpFlag).Parse()
	if err != nil {
		return cfg, err
	}

	err = cfg.parseIn()

	// (Quick fix) Parse flags again to override config file fields
	_, err = flags.NewParser(&cfg, flags.HelpFlag).Parse()
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

func ParseConfigFile(path string) (Config, error) {
	cfg := NewDefaultConfig()
	cfg.In = path
	err := cfg.parseIn()
	return cfg, err
}

func (cfg Config) ToString() string {
	j, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(j)
}

func (cfg Config) ToScene() (*scene.Node, error) {
	return cfg.toScene()
}

type Image struct {
	Width  int `json:"width" short:"w" long:"width" description:"Output image width"`
	Height int `json:"height" short:"h" long:"height" description:"Output image height"`
}

func NewDefaultImage() Image {
	return Image{
		Width:  800,
		Height: 600,
	}
}

type Process struct {
	Threads          int     `json:"threads" short:"t" long:"threads" description:"The number of threads to use for rendering"`
	Spp              int     `json:"spp" short:"s" long:"spp" description:"Samples per pixel"`
	Alpha            float64 `json:"alpha" short:"a" long:"alpha" description:"Alpha parameter for PHR"`
	Delta            float64 `json:"delta" short:"d" long:"delta" description:"Delta parameter for PHR"`
	UsePhr           bool    `json:"usePhr" short:"p" long:"phr" description:"If present, apply PHR after initial BVH construction"`
	Heatmap          bool    `json:"heatmap" long:"heatmap" description:"If present, render heatmap of the bvh"`
	HeatmapThreshold int     `json:"heatmapThreshold" long:"heatmapThreshold" description:"Threshold at which heatmap shows red"`
}

func NewDefaultProcess() Process {
	return Process{
		Threads:          runtime.NumCPU(),
		Spp:              300,
		Alpha:            0.5,
		Delta:            6,
		UsePhr:           false,
		Heatmap:          false,
		HeatmapThreshold: 100,
	}
}

type Scene struct {
	Objects []Object `json:"objects"`
	Camera  Camera   `json:"camera"`
}

func NewDefaultScene() Scene {
	return Scene{
		Objects: []Object{},
		Camera: Camera{
			LookFrom: [3]float64{1, 1, 1},
			LookAt:   [3]float64{0, 0, 0},
			Up:       [3]float64{0, 1, 0},
			Fov:      90,
		},
	}
}

type Object struct {
	File     string     `json:"file"`
	Scale    [3]float64 `json:"scale"`
	Position [3]float64 `json:"position"`
	Material Material   `json:"material"`
}

type Material struct {
	Type       string     `json:"type"`
	Albedo     [3]float64 `json:"albedo"`
	Diffustion float64    `json:"diffusion"`
	Ratio      float64    `json:"ratio"`
	Emmitance  float64    `json:"emmitance"`
}

type Camera struct {
	LookFrom [3]float64 `json:"lookFrom"`
	LookAt   [3]float64 `json:"lookAt"`
	Up       [3]float64 `json:"up"`
	Fov      float64    `json:"fov" long:"fov" description:"Field of View"`
}

func (cfg *Config) parseIn() error {
	err := cfg.parseJsonConfig()
	if err != nil {
		return fmt.Errorf("failed to parse json config: %w", err)
	}

	cfg.parseObj()
	return nil
}

func (cfg *Config) parseObj() {
	if !strings.HasSuffix(cfg.In, ".obj") {
		return
	}

	cfg.Scene.Objects = append(cfg.Scene.Objects, Object{
		File: cfg.In,
		Material: Material{
			Type:   "diffuse",
			Albedo: [3]float64{1, 0, 0},
		},
		Position: [3]float64{0, 0, 0},
		Scale:    [3]float64{1, 1, 1},
	})
}

func (cfg *Config) parseJsonConfig() error {
	if !strings.HasSuffix(cfg.In, ".json") {
		return nil
	}

	f, err := os.Open(cfg.In)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, cfg)
}
