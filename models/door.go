package models

import (
	"fmt"
	"sync"
)

type request struct {
	vehicleID int
	direction string
	allowPass chan bool
}

type Door struct {
	queue     chan request
	mu        sync.Mutex
	cond      *sync.Cond
	count     int
	direction string // "entering" o "exiting"
}

func NewDoor() *Door {
	d := &Door{
		queue: make(chan request, 100),
	}
	d.cond = sync.NewCond(&d.mu)
	return d
}

func (d *Door) Manage() {
	for req := range d.queue {
		d.mu.Lock()

		// Esperar hasta que la puerta esté libre o en la dirección correcta
		for d.count > 0 && d.direction != req.direction {
			d.cond.Wait()
		}

		// Cambia la dirección y permite el paso del vehículo
		d.direction = req.direction
		d.count++
		req.allowPass <- true // Enviar señal de permiso para que el vehículo pase
		close(req.allowPass)  // Cerrar el canal para indicar que se dio la señal

		fmt.Printf("Vehicle %d is %s through the door. Count: %d\n", req.vehicleID, req.direction, d.count)
		d.mu.Unlock()

		// Simulate vehicle passing through the door
		d.CompleteRequest(req.vehicleID)

		// Esperar a que el vehículo complete su paso
		d.mu.Lock()
		for d.count > 0 && d.direction == req.direction {
			d.cond.Wait()
		}
		d.mu.Unlock()
	}
}

func (d *Door) CompleteRequest(vehicleID int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.count--
	fmt.Printf("Vehicle %d has completed %s through the door. Remaining count: %d\n", vehicleID, d.direction, d.count)

	// Cambia la dirección solo cuando no hay más vehículos en la dirección actual
	if d.count == 0 {
		d.direction = "" // Liberar la dirección actual
		d.cond.Broadcast() // Despertar a otros vehículos en cola para que intenten pasar
	}
}

// Solicitud para entrar o salir con un canal de confirmación
func (d *Door) Request(vehicleID int, direction string) <-chan bool {
	allowPass := make(chan bool)
	d.queue <- request{vehicleID, direction, allowPass}
	return allowPass
}
