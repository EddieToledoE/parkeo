package controllers

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/EddieToledoE/parkingsimulator/models"

	"github.com/faiface/pixel"
)

func SimulateVehicles(parkingLot *models.ParkingLot, door *models.Door, vehicleSprite *pixel.Sprite, vehicles *[]*models.Vehicle, wg *sync.WaitGroup, maxVehicles int) {
	go door.Manage() // Iniciar el manejo de la puerta

	activeVehicles := make(chan struct{}, maxVehicles)
	var mu sync.Mutex

	for i := 1; i <= 100; i++ {
		activeVehicles <- struct{}{}
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() { <-activeVehicles }()
			// crear carro
			vehicle := &models.Vehicle{
				ID:       id,
				State:    "waiting",
				Position: pixel.V(400, 525),
				Sprite:   vehicleSprite,
				Speed:    100, // Ajusta la velocidad del vehículo
			}

			mu.Lock()
			*vehicles = append(*vehicles, vehicle)
			mu.Unlock()

			fmt.Printf("Vehicle %d created and waiting to enter.\n", id)

			// Solicitar entrada y esperar confirmacion
			allowEnter := door.Request(id, "entering")
			<-allowEnter // Espera la señal de permiso para entrar
			slot := parkingLot.Enter(id)
			vehicle.Target = calculateSlotPosition(slot)
			vehicle.State = "entering"

			// Mueve el vehículo hacia el slot
			for vehicle.State == "entering" {
				vehicle.Update(0.1) // Actualizar posicion
				time.Sleep(100 * time.Millisecond) // simular movimiento
			}

			// Simulación de estacionamiento
			vehicle.ParkedAt = time.Now()
			time.Sleep(time.Duration(10) * time.Second) // Tiempo estacionado

			for {
				// Verificar y solicitar salida si el vehículo está listo
				if vehicle.ShouldExit() {
					fmt.Printf("Vehicle %d ready to exit from slot %d.\n", id, slot)

					// Solicitar salida y esperar confirmacion
					allowExit := door.Request(id, "exiting")
					select {
					case <-allowExit: // Espera la señal de permiso para salir
						vehicle.State = "exiting"
						vehicle.Target = pixel.V(400, 525)

						// Mueve hacia la salida
						for vehicle.State == "exiting" {
							vehicle.Update(0.1) // Actualizar posicion
							time.Sleep(100 * time.Millisecond)
						}

						// Liberar el espacio en el estacionamiento después de recibir la confirmación de salida
						parkingLot.Exit(id, slot)
						fmt.Printf("Vehicle %d is exiting from slot %d.\n", id, slot)

						// Completa la solicitud de salida en la puerta
						door.CompleteRequest(id)
						fmt.Printf("Vehicle %d has left.\n", id)
						return // Rompe el bucle y termina la goroutine
					case <-time.After(5 * time.Second): // Timeout para evitar bloqueo
						fmt.Printf("Vehicle %d exit request timed out.\n", id)
					}
				} else {
					fmt.Printf("Vehicle %d not ready to exit yet.\n", id)
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				}
			}
		}(i)

		// Introduce un retardo entre la generación
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}


func calculateSlotPosition(slot int) pixel.Vec {
	slotWidth, slotHeight := 100.0, 50.0
	offsetX, offsetY := 150.0, 100.0

	if slot < 8 {
		return pixel.V(offsetX+slotWidth/2, 500-float64(slot)*slotHeight+slotHeight/2)
	} else if slot < 16 {
		return pixel.V(600+slotWidth/2, 500-float64(slot-8)*slotHeight+slotHeight/2)
	} else {
		return pixel.V(300+float64(slot-16)*slotWidth+slotWidth/2, offsetY+slotHeight/2)
	}
}
