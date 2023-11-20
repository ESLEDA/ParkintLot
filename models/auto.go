	package models

	import (
		"image/color"
		"math/rand"
		"sync"
		"time"

		"fyne.io/fyne/v2"
		"fyne.io/fyne/v2/canvas"
		"fyne.io/fyne/v2/storage"
	)

	
	type Auto struct {
		id              int
		tiempoLim       time.Duration
		espacioAsignado int
		imagenEntrada   *canvas.Image
		imagenEspera    *canvas.Image
		imagenSalida    *canvas.Image
	}

	
	func NewAuto(id int) *Auto {
		
		imagenEntrada := canvas.NewImageFromURI(storage.NewFileURI("./assets/entrada.jpeg"))
		imagenEspera := canvas.NewImageFromURI(storage.NewFileURI("./assets/derecha.jpeg"))
		imagenSalida := canvas.NewImageFromURI(storage.NewFileURI("./assets/derecha.jpeg"))
		imagenEspera.Resize(fyne.NewSize(50, 30))
		return &Auto{
		id:              id,
		tiempoLim:       time.Duration(rand.Intn(50)+50) * time.Second,
		espacioAsignado: 0,
		imagenEntrada:   imagenEntrada,
		imagenEspera:    imagenEspera,
		imagenSalida:    imagenSalida,
	}
	}

	
	func (a *Auto) Entrar(p *Estacionamiento, contenedor *fyne.Container) {
		p.GetEspacios() <- a.GetId() 
		p.GetPuertaMu().Lock()       

		espacios := p.GetEspaciosArray()
		const (
			carrosPorFila     = 10
			espacioHorizontal = 50
			espacioVertical   = 130
		)

		
		for i := 0; i < len(espacios); i++ {
			if !espacios[i] {
				espacios[i] = true
				a.espacioAsignado = i

				fila := i / carrosPorFila
				columna := i % carrosPorFila

				x := float32(55 + columna*espacioHorizontal)
				y := float32(250 + fila*espacioVertical)

				a.imagenEntrada.Move(fyne.NewPos(x, y))
				break
			}
		}

		
		line := canvas.NewLine(color.RGBA{0, 255, 0, 255})
		line.StrokeWidth = 2
		line.Position1 = fyne.NewPos(a.imagenEntrada.Position().X+45, a.imagenEntrada.Position().Y)
		line.Position2 = fyne.NewPos(a.imagenEntrada.Position().X+35, a.imagenEntrada.Position().Y+55)
		contenedor.Add(line)

		p.SetEspaciosArray(espacios)
		p.GetPuertaMu().Unlock() 
		contenedor.Refresh()     
	}


	
	func (a *Auto) Salir(p *Estacionamiento, contenedor *fyne.Container) {
		<-p.GetEspacios()      
		p.GetPuertaMu().Lock() 

		spacesArray := p.GetEspaciosArray()
		spacesArray[a.espacioAsignado] = false
		p.SetEspaciosArray(spacesArray)

		p.GetPuertaMu().Unlock() 

		contenedor.Remove(a.imagenEspera)
		a.imagenSalida.Resize(fyne.NewSize(30, 50))
		a.imagenSalida.Move(fyne.NewPos(50, 290))

		contenedor.Add(a.imagenSalida)
		contenedor.Refresh()

		
		for i := 0; i < 10; i++ {
			a.imagenSalida.Move(fyne.NewPos(a.imagenSalida.Position().X+30, a.imagenSalida.Position().Y))
			time.Sleep(time.Millisecond * 370)
		}

		contenedor.Remove(a.imagenSalida)
		contenedor.Refresh()
	}



	
	func (a *Auto) Iniciar(p *Estacionamiento, contenedor *fyne.Container, wg *sync.WaitGroup) {
		a.Avanzar(9) 

		a.Entrar(p, contenedor) 

		time.Sleep(a.tiempoLim) 

		contenedor.Remove(a.imagenEntrada)
		a.imagenEspera.Resize(fyne.NewSize(90, 30))
		p.ColaSalida(contenedor, a.imagenEspera) 
		a.Salir(p, contenedor) 

		wg.Done() 
	}

	
	func (a *Auto) Avanzar(pasos int) {
		for i := 0; i < pasos; i++ {
			a.imagenEntrada.Move(fyne.NewPos(a.imagenEntrada.Position().X, a.imagenEntrada.Position().Y+20))
			time.Sleep(time.Millisecond * 200)
		}
	}

	
	func (a *Auto) GetId() int {
		return a.id
	}

	
	func (a *Auto) GetImagenEntrada() *canvas.Image {
		return a.imagenEntrada
	}
