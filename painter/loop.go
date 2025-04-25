package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type MessageQueue struct {
	queue chan Operation
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	mq       MessageQueue
	state    TextureState
	doneFunc func()
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.mq = MessageQueue{queue: make(chan Operation)}
	l.state = TextureState{backgroundColor: &Fill{Color: color.White}}

	go func() {
		for {
			e := l.mq.Pull()

			switch e.(type) {
			case Figure, BgRect, Move, Fill, Reset:
				e.Update(&l.state)
			case Update:
				l.state.backgroundColor.Do(l.next)

				if l.state.backgroundRect != nil {
					l.state.backgroundRect.Do(l.next)
				}

				for _, fig := range l.state.figureCenters {
					fig.Do(l.next)
				}
				l.prev = l.next
				l.Receiver.Update(l.next)
				l.next, _ = s.NewTexture(size)
			}

			if l.doneFunc != nil {
				l.doneFunc()
			}
		}
	}()
}

func (l *Loop) Post(ol OperationList) {

	for _, op := range ol {
		l.mq.Push(op)
	}
}

func (l *Loop) StopAndWait() {
}

func (mq *MessageQueue) Push(op Operation) {
	mq.queue <- op
}

func (mq *MessageQueue) Pull() Operation {
	return <-mq.queue
}
