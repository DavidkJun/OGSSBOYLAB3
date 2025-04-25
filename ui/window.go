package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz     size.Event
	pos    image.Rectangle
	center image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 400
	pw.pos.Max.Y = 400
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event:
		pw.sz = e
		pw.center = image.Pt(pw.sz.WidthPx/2, pw.sz.HeightPx/2)
		fmt.Println("resized to", pw.sz.HeightPx, "and", pw.sz.WidthPx)

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == 3 && e.Direction == 1 {
				pw.center.Y, pw.center.X = int(e.Y), int(e.X)
			}
			pw.w.Send(paint.Event{})
		}

	case paint.Event:
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawFigure() {
	DrawFigure(pw.w, pw.center)
}

func DrawFigure(screenUploader screen.Uploader, imagePoint image.Point) {
	scale := 1
	colorT := color.RGBA{
		R: 255,
		G: 255,
		B: 0,
		A: 0,
	}

	screenUploader.Fill(
		image.Rect(imagePoint.X-225*scale, imagePoint.Y-175*scale, imagePoint.X+225*scale, imagePoint.Y),
		colorT,
		draw.Src,
	)

	screenUploader.Fill(
		image.Rect(imagePoint.X-75*scale, imagePoint.Y-175*scale, imagePoint.X+75*scale, imagePoint.Y+250*scale),
		colorT,
		draw.Src,
	)
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src)

	pw.drawFigure()

	for _, br := range imageutil.Border(pw.sz.Bounds(), 5) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
