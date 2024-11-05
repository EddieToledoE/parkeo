package main

import (
	"carro/scenes"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(scenes.RunSimulation)
}
