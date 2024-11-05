package scenes

import (
	"carro/models"
	"carro/views"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func RunSimulation() {
	models.InitVehicleSystem()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for vehicle := range models.VehicleChannel {
			go models.ManageLane(vehicle.ID)
		}
	}()

	for !win.Closed() {
		win.Clear(colornames.White)
		views.RenderParkingLot(win, models.GetVehicles())
		win.Update()

		models.VehiclesMutex.Lock()
		models.MoveVehiclesLogic()
		models.VehiclesMutex.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}
