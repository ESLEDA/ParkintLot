package models

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"sync"
)

type Estacionamiento struct {
	espacios      chan int
	puerta        *sync.Mutex
	espaciosArray [20]bool
}

func NewEstacionamiento(espacios chan int, puertaMu *sync.Mutex) *Estacionamiento {
	return &Estacionamiento{
		espacios:      espacios,
		puerta:        puertaMu,
		espaciosArray: [20]bool{},
	}
}

func (p *Estacionamiento) GetEspacios() chan int {
	return p.espacios
}

func (p *Estacionamiento) GetPuertaMu() *sync.Mutex {
	return p.puerta
}

func (p *Estacionamiento) GetEspaciosArray() [20]bool {
	return p.espaciosArray
}

func (p *Estacionamiento) SetEspaciosArray(espaciosArray [20]bool) {
	p.espaciosArray = espaciosArray
}

func (p *Estacionamiento) ColaSalida(contenedor *fyne.Container, imagen *canvas.Image) {
	imagen.Move(fyne.NewPos(80, 280))
	contenedor.Add(imagen)
	contenedor.Refresh()
}
