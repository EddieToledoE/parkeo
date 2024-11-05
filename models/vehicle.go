package models

import (
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
)

type Vehicle struct {
	ID               int
	Position         pixel.Vec
	PreviousPosition pixel.Vec
	Lane             int
	Parked           bool
	ExitTime         time.Time
	IsEntering       bool
	Teleporting      bool
	TeleportStartTime time.Time
}

var Vehicles []Vehicle
var VehicleChannel chan Vehicle
var VehiclesMutex sync.Mutex
var CarEnteringOrExiting bool


func InitVehicleSystem() {
	VehicleChannel = make(chan Vehicle)
	go VehicleGenerator()
}

func CreateVehicle(id int) Vehicle {
	VehiclesMutex.Lock()
	defer VehiclesMutex.Unlock()
	vehicle := Vehicle{
		ID:       id,
		Position: pixel.V(0, 300),
		Lane:     -1,
		Parked:   false,
	}
	Vehicles = append(Vehicles, vehicle)
	return vehicle
}

func SetExitTime(vehicle *Vehicle) {
	rand.Seed(time.Now().UnixNano())
	exitIn := time.Duration(rand.Intn(5)+1) * time.Second
	vehicle.ExitTime = time.Now().Add(exitIn)
}

func GetVehicles() []Vehicle {
	return Vehicles
}

func AssignLaneToVehicle(id int, lane int) {
	VehiclesMutex.Lock()
	defer VehiclesMutex.Unlock()
	for i := range Vehicles {
		if Vehicles[i].ID == id {
			Vehicles[i].Lane = lane
		}
	}
}

func ResetVehiclePosition(id int) {
	VehiclesMutex.Lock()
	defer VehiclesMutex.Unlock()
	for i := range Vehicles {
		if Vehicles[i].ID == id {
			Vehicles[i].Position = pixel.V(0, 300)
		}
	}
}

func FindVehiclePosition(id int) pixel.Vec {
	VehiclesMutex.Lock()
	defer VehiclesMutex.Unlock()
	for _, vehicle := range Vehicles {
		if vehicle.ID == id {
			return vehicle.Position
		}
	}
	return pixel.Vec{}
}

func ParkVehicle(vehicle *Vehicle, targetX, targetY float64) {
	vehicle.Position.X = targetX
	vehicle.Position.Y = targetY
	vehicle.Parked = true
	SetExitTime(vehicle)
}

func removeVehicle(index int) {
	Vehicles = append(Vehicles[:index], Vehicles[index+1:]...)
}

func VehicleGenerator() {
	id := 0
	for {
		id++
		vehicle := CreateVehicle(id)
		VehicleChannel <- vehicle
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))
	}
}
