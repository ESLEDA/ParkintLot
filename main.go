package main

import (
	"estacionamiento/models"
	"image/color"
	
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"gonum.org/v1/gonum/stat/distuv"
)

type MainScene struct {
	window fyne.Window
}

func NewMainScene(window fyne.Window) *MainScene {
	return &MainScene{
		window: window,
	}
}

var contenedor = container.NewWithoutLayout()

func (s *MainScene) Show() {
	// Añade un fondo marrón al contenedor
	fondo := canvas.NewRectangle(color.RGBA{R: 139, G: 69, B: 19, A: 255})
	fondo.Resize(fyne.NewSize(885, 700))
	contenedor.Add(fondo)

	s.window.SetContent(contenedor)
}

func (s *MainScene) Run() {
	p := models.NewEstacionamiento(make(chan int, 20), &sync.Mutex{})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			auto := models.NewAuto(id)
			imagen := auto.GetImagenEntrada()
			imagen.Resize(fyne.NewSize(30, 50))

			// Ajusta la posición inicial para entrar en medio de ambas filas
			x := 10 // Ajusta según tu diseño
			y := 60 // Todos los autos en la misma fila
			imagen.Move(fyne.NewPos(float32(x), float32(y)))

			contenedor.Add(imagen)
			contenedor.Refresh()

			auto.Iniciar(p, contenedor, &wg)
		}(i)
		var poisson = generarPoisson(float64(2))
		time.Sleep(time.Second * time.Duration(poisson))
	}

	wg.Wait()
}

func generarPoisson(lambda float64) float64 {
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	return poisson.Rand()
}

type MainView struct{}

func NewMainView() *MainView {
	return &MainView{}
}

func (v *MainView) Run() {
	myApp := app.New()
	window := myApp.NewWindow("parkiglot")
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(600, 570))

	mainScene := NewMainScene(window)
	mainScene.Show()
	go mainScene.Run()
	window.ShowAndRun()
}

func main() {
	mainView := NewMainView()
	mainView.Run()
}
