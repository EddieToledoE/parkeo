package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/EddieToledoE/parkingsimulator/controllers"
	"github.com/EddieToledoE/parkingsimulator/models"
	"github.com/EddieToledoE/parkingsimulator/views"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	rand.Seed(time.Now().UnixNano())

	parkingLot := models.NewParkingLot()
	door := models.NewDoor()

	cfg := pixelgl.WindowConfig{
		Title:  "Parking Simulator",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Cargar imageness
	backgroundPic, err := views.LoadPicture("assets/bg.png")
	if err != nil {
		log.Fatalf("Failed to load background image: %v", err)
	}
	background := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	vehiclePic, err := views.LoadPicture("assets/car.png")
	if err != nil {
		log.Fatalf("Failed to load vehicle image: %v", err)
	}
	vehicleSprite := pixel.NewSprite(vehiclePic, vehiclePic.Bounds())

	var wg sync.WaitGroup
	var vehicles []*models.Vehicle

	// Iniciar los carros
	go controllers.SimulateVehicles(parkingLot, door, vehicleSprite, &vehicles, &wg,20)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		for _, vehicle := range vehicles {
			vehicle.Update(dt)
		}

		views.DrawParkingLot(win, parkingLot, background, vehicles)
		win.Update()
	}

	wg.Wait()
}

func main() {
	pixelgl.Run(run)
}
