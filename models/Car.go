package models

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"math/rand"
	"sync"
	"time"
)

// ParkingLot represents a parking lot structure.
type ParkingLot struct {
	ParkingSpaces []bool      // ParkingSpaces stores availability of parking spaces.
	CapacityMutex *sync.Mutex // CapacityMutex synchronizes access to parking lot capacity.
	arriveSta     int         // Counter for arriving vehicles.
	exitSta       int         // Counter for exiting vehicles.
}

// NewParkingLot creates a new instance of ParkingLot with the given number of parking spaces.
func NewParkingLot(numSpaces int) *ParkingLot {
	parkingSpaces := make([]bool, numSpaces)
	capacityMutex := &sync.Mutex{}
	arriveSta := 0
	exitSta := 0
	return &ParkingLot{
		exitSta:       exitSta,
		arriveSta:     arriveSta,
		ParkingSpaces: parkingSpaces,
		CapacityMutex: capacityMutex,
	}
}

// Vehiculos simulates the behavior of vehicles arriving and departing from the parking lot.
func (pl *ParkingLot) Vehiculos(id int, wg *sync.WaitGroup, w *fyne.Container) {
	defer wg.Done()

	fmt.Printf("Vehicle %d arrives.\n", id)
	pl.arriveSta++
	img, _ := fyne.LoadResourceFromPath("assets/coche-rojo.png")
	auto := canvas.NewImageFromResource(img)
	auto.Resize(fyne.NewSize(32, 32))
	w.Add(auto)
	w.Refresh()

	// Set initial position of the vehicle.
	if pl.arriveSta != 0 {
		xs := 120 - 20*pl.arriveSta
		ys := 260
		auto.Move(fyne.NewPos(float32(xs), float32(ys)))
	} else {
		auto.Move(fyne.NewPos(120, 260))
	}
	time.Sleep(200 * time.Millisecond)

	for {
		// Check if the parking lot is full, otherwise wait until a space is available.
		pl.CapacityMutex.Lock()
		if pl.isFull() {
			pl.CapacityMutex.Unlock()
			fmt.Printf("Vehicle %d is blocked as the parking lot is full.\n", id)
			if pl.arriveSta != 0 {
				xs := 120 - 20*pl.arriveSta
				ys := 260
				auto.Move(fyne.NewPos(float32(xs), float32(ys)))
			} else {
				auto.Move(fyne.NewPos(120, 260))
			}
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Find an available parking space and park the vehicle.
		for i, occupied := range pl.ParkingSpaces {
			if !occupied {
				pl.ParkingSpaces[i] = true
				fmt.Printf("Vehicle %d parked in space %d.\n", id, i)
				pl.arriveSta--

				// Adjust position of the vehicle in the parking space.
				xh := 175
				yh := 60
				if i > 10 {
					xh = (65 * (i - 11)) + xh
					yh = 430
				} else {
					xh = 65*i + xh
				}
				// Move the vehicle to the parking space.
				auto.Move(fyne.NewPos(float32(xh), float32(yh-2)))
				// Simulate parking animation with delays.
				for j := 0; j < 5; j++ {
					time.Sleep(200 * time.Millisecond)
					auto.Move(fyne.NewPos(float32(xh), float32(yh-2)))
				}
				pl.CapacityMutex.Unlock()

				// Simulate vehicle parked for a random duration.
				min := 1
				max := 5
				randomDuration := time.Duration(min+rand.Intn(max-min+1)) * time.Second
				time.Sleep(randomDuration)

				pl.CapacityMutex.Lock()
				// Vehicle departs from the parking lot.
				pl.ParkingSpaces[i] = false
				fmt.Printf("Vehicle %d goes to the exit.\n", id)
				pl.exitSta++

				// Adjust position of the vehicle when leaving.
				if pl.exitSta != 0 {
					xs := 160 + 20*pl.arriveSta
					ys := 200
					auto.Move(fyne.NewPos(float32(xs), float32(ys)))
				} else {
					auto.Move(fyne.NewPos(160, 200))
				}
				// Simulate exit delay before leaving.
				time.Sleep(time.Second * time.Duration(rand.Intn(2)))
				fmt.Printf("Vehicle %d leaves the parking lot.\n", id)
				pl.exitSta--
				auto.Move(fyne.NewPos(120, 200))
				time.Sleep(200 * time.Millisecond)
				auto.Hidden = true
				pl.CapacityMutex.Unlock()
				return
			}
		}

		pl.CapacityMutex.Unlock()
	}
}

// isFull checks if all parking spaces are occupied.
func (pl *ParkingLot) isFull() bool {
	for _, occupied := range pl.ParkingSpaces {
		if !occupied {
			return false
		}
	}
	return true
}
