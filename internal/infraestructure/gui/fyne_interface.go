package gui

import (
    "fmt"
    "os"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
    "image/color"
    "time"
    "simulador/internal/domain"
)

type ParkingSpot struct {
    container *fyne.Container
    carImage  *canvas.Image
    label     *widget.Label
    occupied  bool
}

func loadCarImage() *canvas.Image {
    img := canvas.NewImageFromFile("resources/car.png")
    if img == nil {
        rect := canvas.NewRectangle(color.NRGBA{R: 100, G: 100, B: 100, A: 255})
        rect.Resize(fyne.NewSize(50, 30))
        return canvas.NewImageFromResource(nil)
    }
    img.FillMode = canvas.ImageFillContain
    img.SetMinSize(fyne.NewSize(50, 30))
    return img
}

func StartGUI(parkingLot *domain.ParkingLot) {
    os.Setenv("FYNE_RENDERER", "software")
    
    application := app.New()
    window := application.NewWindow("Simulador de Estacionamiento")

    mainContainer := container.NewVBox()
    
    title := widget.NewLabel("Estado del Estacionamiento")
    title.TextStyle = fyne.TextStyle{Bold: true}
    mainContainer.Add(title)

    statusLabel := widget.NewLabel("Espacios ocupados: 0/" + fmt.Sprint(parkingLot.Capacity))
    mainContainer.Add(statusLabel)

    spotsContainer := container.NewGridWithColumns(5)
    spots := make([]*ParkingSpot, parkingLot.Capacity)

    for i := 0; i < parkingLot.Capacity; i++ {
        spot := &ParkingSpot{
            carImage: loadCarImage(),
            label:    widget.NewLabel(fmt.Sprintf("Espacio %d\nLibre", i+1)),
            occupied: false,
        }
        
        spot.carImage.Hide()
        
        spotBg := canvas.NewRectangle(color.NRGBA{G: 200, A: 255})
        spotBg.Resize(fyne.NewSize(100, 60))
        
        spot.container = container.NewVBox(
            container.NewStack(
                spotBg,
                spot.carImage,
            ),
            spot.label,
        )
        
        spots[i] = spot
        spotsContainer.Add(spot.container)
    }

    mainContainer.Add(spotsContainer)

    go func() {
        for {
            select {
            case <-parkingLot.UpdateChan:
                occupiedVehicles := parkingLot.GetOccupiedSpots()
                occupiedCount := len(occupiedVehicles)
                
                statusLabel.SetText(fmt.Sprintf("Espacios ocupados: %d/%d", 
                    occupiedCount, parkingLot.Capacity))

                for i := 0; i < parkingLot.Capacity; i++ {
                    spot := spots[i]
                    if i < occupiedCount && !spot.occupied {
                        spot.carImage.Show()
                        spot.label.SetText(fmt.Sprintf("Espacio %d\nOcupado\nVehículo %d", 
                            i+1, occupiedVehicles[i].ID))
                        spot.occupied = true
                    } else if i >= occupiedCount && spot.occupied {
                        spot.carImage.Hide()
                        spot.label.SetText(fmt.Sprintf("Espacio %d\nLibre", i+1))
                        spot.occupied = false
                    }
                }
            default:
                time.Sleep(time.Millisecond * 100)
            }
        }
    }()

    window.SetContent(mainContainer)
    window.Resize(fyne.NewSize(600, 400))
    window.ShowAndRun()
}