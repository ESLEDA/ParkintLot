package view

import (
	"Parking_lot/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"math/rand"
	"sync"
	"time"
)
// CallCar es una función que crea una interfaz de estacionamiento concurrente
func CallCar(numVehicles int) {
	var wg sync.WaitGroup

	a := app.New()
	icon, _ := fyne.LoadResourceFromPath("assets/icon.png")
	a.SetIcon(icon)
	w := a.NewWindow("EstacionamientoConcurrencia")
	w.Resize(fyne.NewSize(800, 500))
	w.SetFixedSize(true)
	img, _ := fyne.LoadResourceFromPath("assets/parking.png")
	image := canvas.NewImageFromResource(img)
	image.Resize(fyne.NewSize(800, 500))
	rand.Seed(time.Now().UnixNano())

	parkingLot := models.NewParkingLot(20)

	cnt := container.NewWithoutLayout(image)
	w.SetContent(cnt)
	// Inicia una goroutine para simular la llegada de vehículos al estacionamiento
	go func() {
		time.Sleep(2 * time.Second)
		// Simula la llegada de "numVehicles" vehículos al estacionamiento
		for i := 1; i <= numVehicles; i++ {
			wg.Add(1)// Incrementa el contador de WaitGroup
			go parkingLot.Vehiculos(i, &wg, cnt)
			time.Sleep(time.Duration(rand.ExpFloat64() * float64(time.Second)))// Espera un tiempo aleatorio antes de simular el siguiente vehículo
		}
		wg.Wait()
	}()
	w.ShowAndRun()
}
