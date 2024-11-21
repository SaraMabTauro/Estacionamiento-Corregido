package domain

import(
	"fmt"
	"math/rand"
	"time"
)

type Vehicle struct {
	ID int
}


func NewVehicle(id int) *Vehicle {
	return &Vehicle{
		ID: id,
	}
}

func (v *Vehicle) SimulateStay() {
	stayDuration := rand.Intn(3) + 3
	fmt.Printf("Vehiculo %d esta estacionado por %d segundos. \n", v.ID, stayDuration)
	time.Sleep(time.Duration(stayDuration) * time.Second)
}