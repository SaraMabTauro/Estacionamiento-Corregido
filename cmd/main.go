package main

import (
    "fmt"
    "sync"
    "os"

    "parking-simulator/internal/domain"
    "parking-simulator/internal/infrastructure/gui"
    "parking-simulator/internal/usercase"
)

func init() {
    os.Setenv("FYNE_RENDERER", "software")
}

func main() {
    fmt.Println("Iniciando simulación del estacionamiento...")

    parkingLot := domain.NewParkingLot(20)
    wg := &sync.WaitGroup{}

    // Iniciar la interfaz gráfica
    go gui.StartGUI(parkingLot)

    // Iniciar la simulación
    usecase.StartSimulation(parkingLot, wg)

    // Esperar a que todas las goroutines finalicen
    wg.Wait()
}