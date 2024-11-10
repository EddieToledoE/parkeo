package models

import (
	"fmt"
	"sync"
)

const ParkingCapacity = 20 // Espacios totales

type ParkingLot struct {
	sync.Mutex
	Slots     chan struct{}
	Occupied  map[int]bool // encontrar espacios ocupados
	Available []int        // Espacios disponibles
}

func NewParkingLot() *ParkingLot {
	available := make([]int, ParkingCapacity)
	for i := 0; i < ParkingCapacity; i++ {
		available[i] = i
	}
	return &ParkingLot{
		Slots:     make(chan struct{}, ParkingCapacity),
		Occupied:  make(map[int]bool),
		Available: available,
	}
}

func (p *ParkingLot) Enter(vehicleID int) int {
	p.Lock()
	defer p.Unlock()

	if len(p.Slots) == ParkingCapacity {
		fmt.Printf("Vehicle %d is waiting for a slot...\n", vehicleID)
	}

	p.Slots <- struct{}{}
	slot := p.Available[0]
	p.Available = p.Available[1:]
	p.Occupied[slot] = true
	fmt.Printf("Vehicle %d has entered slot %d.\n", vehicleID, slot)
	return slot
}

func (p *ParkingLot) Exit(vehicleID, slot int) {
	p.Lock()
	defer p.Unlock()

	// Verifica el estado antes de liberar el espacio
	fmt.Printf("Vehicle %d attempting to leave slot %d. Slots occupied: %d\n", vehicleID, slot, len(p.Slots))

	<-p.Slots // Liberar un espacio en el canal
	delete(p.Occupied, slot)
	p.Available = append(p.Available, slot)

	// Confirmación de la liberación del espacio
	fmt.Printf("Vehicle %d has left slot %d. Slots now occupied: %d\n", vehicleID, slot, len(p.Slots))
}

