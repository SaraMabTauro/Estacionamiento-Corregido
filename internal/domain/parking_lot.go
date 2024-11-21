package domain

import (
	"fmt"
	"sync"
)

type ParkingLot struct {
	Capacity int
	Vehicles map[int]*Vehicle  
	Mutex    sync.Mutex
	// canal para notificaciones de cambios
	UpdateChan chan struct{}
}

func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		Capacity:   capacity,
		Vehicles:   make(map[int]*Vehicle),
		UpdateChan: make(chan struct{}, 1),
	}
}

func (p *ParkingLot) Enter(v *Vehicle) bool {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if len(p.Vehicles) < p.Capacity {
		p.Vehicles[v.ID] = v
		fmt.Printf("Vehículo %d ha entrado al estacionamiento.\n", v.ID)
		
		// Notificar cambio
		select {
		case p.UpdateChan <- struct{}{}:
		default:
		}
		
		return true
	}
	
	fmt.Printf("Vehículo %d no pudo entrar, estacionamiento lleno.\n", v.ID)
	return false
}

func (p *ParkingLot) Exit(v *Vehicle) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	delete(p.Vehicles, v.ID)
	fmt.Printf("Vehículo %d ha salido del estacionamiento.\n", v.ID)
}

func (p *ParkingLot) GetOccupiedSpots() []*Vehicle {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	
	vehicles := make([]*Vehicle, 0, len(p.Vehicles))
	for _, v := range p.Vehicles {
		vehicles = append(vehicles, v)
	}
	return vehicles
}