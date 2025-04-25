package painter

import (
	"errors"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"testing"
)

type FakeReceiver struct {
	updateCounter int
}

func (r *FakeReceiver) Update(_ screen.Texture) {
	r.updateCounter++
}

type DummyScreen struct{}

func (ds DummyScreen) NewBuffer(_ image.Point) (screen.Buffer, error) {
	return nil, errors.New("not implemented")
}
func (ds DummyScreen) NewTexture(_ image.Point) (screen.Texture, error) {
	return nil, errors.New("not implemented")
}
func (ds DummyScreen) NewWindow(_ *screen.NewWindowOptions) (screen.Window, error) {
	return nil, errors.New("not implemented")
}

type TextureStub struct{}

func (ts TextureStub) Release() {}
func (ts TextureStub) Size() image.Point {
	return image.Point{}
}
func (ts TextureStub) Bounds() image.Rectangle {
	return image.Rectangle{}
}
func (ts TextureStub) Upload(_ image.Point, _ screen.Buffer, _ image.Rectangle) {}
func (ts TextureStub) Fill(_ image.Rectangle, _ color.Color, _ draw.Op)         {}

type SignalMonitor struct {
	expected int
	signal   chan struct{}
}

func (sm SignalMonitor) trigger() {
	sm.signal <- struct{}{}
}

func (sm SignalMonitor) await() {
	for i := 0; i < sm.expected; i++ {
		<-sm.signal
	}
}

func newSignalMonitor(expected int) SignalMonitor {
	return SignalMonitor{expected: expected, signal: make(chan struct{}, expected)}
}

func TestFinalColorIsBlack(t *testing.T) {
	cmds := OperationList{
		Fill{Color: color.RGBA{R: 255, G: 100, B: 0, A: 255}},
		Fill{Color: color.Black},
	}

	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: &FakeReceiver{}, doneFunc: monitor.trigger}

	loop.Start(DummyScreen{})
	loop.Post(cmds)

	monitor.await()

	if loop.state.backgroundColor.Color != color.Black {
		t.Error("Expected final background color to be black")
	}
}

func TestInitialStateIsCorrect(t *testing.T) {
	loop := Loop{Receiver: &FakeReceiver{}}
	loop.Start(DummyScreen{})
	loop.Post(OperationList{})

	if loop.state.backgroundColor.Color != color.White ||
		loop.state.backgroundRect != nil ||
		loop.state.figureCenters != nil {
		t.Error("Default state not properly initialized")
	}
}

func TestFinalFillWins(t *testing.T) {
	cmds := OperationList{
		Fill{Color: color.RGBA{R: 200, G: 100, B: 100, A: 255}},
		Fill{Color: color.Gray{Y: 10}},
		Fill{Color: color.White},
		Fill{Color: color.Black},
	}

	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: &FakeReceiver{}, doneFunc: monitor.trigger}
	loop.Start(DummyScreen{})
	loop.Post(cmds)
	monitor.await()

	if loop.state.backgroundColor.Color != color.Black {
		t.Error("Latest fill should determine the final background color")
	}
}

func TestLastBgRectStored(t *testing.T) {
	cmds := OperationList{
		BgRect{X1: 0.0, Y1: 0.0, X2: 0.2, Y2: 0.3},
		BgRect{X1: 0.6, Y1: 0.7, X2: 0.9, Y2: 0.95},
	}

	expected := BgRect{X1: 0.6, Y1: 0.7, X2: 0.9, Y2: 0.95}
	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: &FakeReceiver{}, doneFunc: monitor.trigger}
	loop.Start(DummyScreen{})
	loop.Post(cmds)
	monitor.await()

	if *loop.state.backgroundRect != expected {
		t.Error("Expected most recent BgRect to be stored")
	}
}

func TestFigurePlacement(t *testing.T) {
	cmds := OperationList{
		Figure{X: 0.05, Y: 0.15},
		Figure{X: 0.25, Y: 0.35},
	}

	expected1 := Figure{X: 0.05, Y: 0.15}
	expected2 := Figure{X: 0.25, Y: 0.35}

	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: &FakeReceiver{}, doneFunc: monitor.trigger}
	loop.Start(DummyScreen{})
	loop.Post(cmds)
	monitor.await()

	if *loop.state.figureCenters[0] != expected1 || *loop.state.figureCenters[1] != expected2 {
		t.Error("Figure coordinates do not match expected values")
	}
}

func TestFullReset(t *testing.T) {
	cmds := OperationList{
		Figure{X: 0.1, Y: 0.1},
		BgRect{X1: 0.2, Y1: 0.2, X2: 0.3, Y2: 0.3},
		Fill{Color: color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 255}},
		Reset{},
	}

	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: &FakeReceiver{}, doneFunc: monitor.trigger}
	loop.Start(DummyScreen{})
	loop.Post(cmds)
	monitor.await()

	if loop.state.figureCenters != nil || loop.state.backgroundRect != nil || loop.state.backgroundColor.Color != color.Black {
		t.Error("Reset failed to restore initial state")
	}
}

func TestReceiverUpdate(t *testing.T) {
	cmds := OperationList{
		Fill{Color: color.RGBA{R: 0xaa, G: 0xbb, B: 0xcc, A: 255}},
		Update{},
	}

	receiver := &FakeReceiver{}
	monitor := newSignalMonitor(len(cmds))
	loop := Loop{Receiver: receiver, doneFunc: monitor.trigger}
	loop.Start(DummyScreen{})
	loop.next = TextureStub{}
	loop.Post(cmds)
	monitor.await()

	if receiver.updateCounter != 1 {
		t.Error("Update command did not notify receiver")
	}
}
