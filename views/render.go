package views

import (
	"carro/models"
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	background *pixel.Sprite
	bgPicture  pixel.Picture
	vehicleSprite *pixel.Sprite
)

const LaneWidth = 150.0

func loadBackground() {
	file, err := os.Open("Assets/bg.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bgPicture = pixel.PictureDataFromImage(img)
	background = pixel.NewSprite(bgPicture, bgPicture.Bounds())
}

func loadVehicleImage(filePath string) *pixel.Sprite {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	vehiclePicture := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(vehiclePicture, vehiclePicture.Bounds())
}

func RenderParkingLot(win *pixelgl.Window, vehicles []models.Vehicle) {
	if background == nil {
		loadBackground()
	}

	if vehicleSprite == nil {
		vehicleSprite = loadVehicleImage("Assets/vehicle.png")
	}

	background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	imd := imdraw.New(nil)

	parkingWidth := 600.0
	laneWidth := parkingWidth / 10
	upperLaneY := 350.0
	lowerLaneY := 100.0

	for i := 0; i < models.NumLanes; i++ {
		xOffset := 100.0 + float64(i%10)*laneWidth
		yOffset := upperLaneY
		if i >= 10 {
			yOffset = lowerLaneY
		}
		imd.Color = colornames.Green
		if models.LaneStatus[i] {
			imd.Color = colornames.Red
		}

		imd.Push(pixel.V(xOffset, yOffset))
		imd.Push(pixel.V(xOffset+laneWidth, yOffset+150.0))
		imd.Rectangle(0)
	}

	imd.Draw(win)

	for _, vehicle := range vehicles {
		vehicleSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(vehicle.Position))
	}
}
