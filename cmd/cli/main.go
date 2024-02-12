package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/schmizzel/go-graphics/pkg/app"
)

func main() {
	cfg, err := app.ParseConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Infof("Started Rendering:\n%s", cfg.ToString())
	start := time.Now()

	err = app.SaveImage(cfg)
	if err != nil {
		log.Errorf("failed to render: %s", err.Error())
		return
	}

	duration := time.Now().Sub(start)
	log.Infof("Finished Rendering in %s", duration)
}
