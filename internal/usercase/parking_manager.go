package usecase

import (
	"fmt"
	"sync"
	"time"

	"simulador/internal/domain"
	"simulador/pkg"
)

func StartSimulation(parkingLot *domain.ParkingLot, wg *sync.WaitGroup) {
	entryGate := make(chan struct{}, 1)

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(vehicleID int) {
			defer wg.Done()
			v := domain.NewVehicle(vehicleID)

			time.Sleep(pkg.PoissonInterval(1.5))

			fmt.Printf("Vehículo %d intenta entrar al estacionamiento.\n", vehicleID)

			entryGate <- struct{}{}
			if parkingLot.Enter(v) {
				<-entryGate
				v.SimulateStay()

				entryGate <- struct{}{}
				parkingLot.Exit(v)
				<-entryGate
			} else {
				<-entryGate
			}
		}(i)
	}
}