package painter

import (
	"image"

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

	mq    MessageQueue
	state TextureState

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = MessageQueue{queue: make(chan Operation, 15)}
	l.state = TextureState{}

	go func() {
		for {
			e := l.mq.Pull()

			switch e.(type) {
			case Figure, BgRect, Move, Fill, Reset:
				e.Update(l.state)
			case Update:
				t, _ := s.NewTexture(size)
				l.state.backgroundColor.Do(t)
				l.state.backgroundRect.Do(t)
				for _, fig := range l.state.figureCenters {
					fig.Do(t)
				}
				l.Receiver.Update(t)
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
