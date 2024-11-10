package views

import (
	"image"
	_ "image/png"

	"os"

	"github.com/EddieToledoE/parkingsimulator/models"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"
)

func SetupWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "Parking Simulator",
		Bounds: pixel.R(0, 0, 800, 600), 
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)
	return win
}

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
func DrawParkingLot(win *pixelgl.Window, parkingLot *models.ParkingLot, background *pixel.Sprite, vehicles []*models.Vehicle) {
	win.Clear(colornames.Skyblue)

	// Dibujar la imagen de fondo si no es nil
	if background != nil {
		background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	}

	imd := imdraw.New(nil)

	// entrada/salida
	entranceExitSlot := pixel.R(350, 500, 450, 450)
	imd.Color = colornames.Green
	imd.Push(entranceExitSlot.Min, entranceExitSlot.Max)
	imd.Rectangle(0)

	// tamaño de los espacios
	slotWidth, slotHeight := 100.0, 50.0
	offsetX, offsetY := 150.0, 100.0

	// Dibujar los laterales de la U (izquierda y derecha)
	for i := 0; i < 8; i++ {
		// Izquierda
		slot := pixel.R(offsetX, 500-float64(i)*slotHeight, offsetX+slotWidth, 550-float64(i)*slotHeight)
		color := colornames.Green
		if parkingLot.Occupied[i] {
			color = colornames.Red
		}
		imd.Color = color
		imd.Push(slot.Min, slot.Max)
		imd.Rectangle(0)

		// Derecha
		slot = pixel.R(600, 500-float64(i)*slotHeight, 600+slotWidth, 550-float64(i)*slotHeight)
		color = colornames.Green
		if parkingLot.Occupied[8+i] {
			color = colornames.Red
		}
		imd.Color = color
		imd.Push(slot.Min, slot.Max)
		imd.Rectangle(0)
	}

	// Dibujar la base de la U
	for i := 0; i < 4; i++ {
		slot := pixel.R(300+float64(i)*slotWidth, offsetY, 350+float64(i)*slotWidth, offsetY+slotHeight)
		color := colornames.Green
		if parkingLot.Occupied[16+i] {
			color = colornames.Red
		}
		imd.Color = color
		imd.Push(slot.Min, slot.Max)
		imd.Rectangle(0)
	}

	// Dibujar el estado actual
	imd.Draw(win)

	// Dibujar vehículos
	for _, vehicle := range vehicles {
		if vehicle.State != "done" && vehicle.Sprite != nil {
			scale := 0.3 // tamaño carro
			mat := pixel.IM.Scaled(pixel.ZV, scale).Moved(vehicle.Position)
			vehicle.Sprite.Draw(win, mat)
		}
	}
}


