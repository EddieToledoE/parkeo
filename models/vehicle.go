package models

import (
	"time"

	"github.com/faiface/pixel"
)

type Vehicle struct {
	ID       int
	Sprite   *pixel.Sprite
	Position pixel.Vec
	State    string // "waiting", "entering", "parked", "exiting", "done"
	Speed    float64
	Target   pixel.Vec
	ParkedAt time.Time
}

func NewVehicle(id int, sprite *pixel.Sprite, startPosition, target pixel.Vec) *Vehicle {
	return &Vehicle{
		ID:       id,
		Sprite:   sprite,
		Position: startPosition,
		State:    "waiting",
		Speed:    150, // Ajusta la velocidad
		Target:   target,
	}
}
func (v *Vehicle) Update(dt float64) {
	if v.State == "entering" || v.State == "exiting" {
		direction := v.Target.Sub(v.Position).Unit()
		v.Position = v.Position.Add(direction.Scaled(v.Speed * dt))

		// Verificar si el vehículo ha llegado al objetivo
		if v.Position.Sub(v.Target).Len() < v.Speed*dt {
			v.Position = v.Target // Asegura que el vehículo se alinee con el objetivo
			if v.State == "entering" {
				v.State = "parked"
				v.ParkedAt = time.Now()
			} else if v.State == "exiting" {
				v.State = "done"
			}
		}
	}
}


func (v *Vehicle) ShouldExit() bool {
	if v.State == "parked" {
		// Verifica si el tiempo de estacionamiento ha pasado
		return time.Since(v.ParkedAt).Seconds() > 10
	}
	return false
}
