package GoTicker

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrTickerIsRunning = errors.New("ticker is running")
var ErrTickerIsStoped = errors.New("ticker is stoped")

const (
	TickerStatusStoped  = 101
	TickerStatusRunning = 102
)

type Ticker struct {
	uuid   string
	done   chan bool
	status int
	ticker *time.Ticker
	fnc    func()
}

func (t *Ticker) Stop() error {
	if t.status == TickerStatusStoped {
		return ErrTickerIsStoped
	}
	t.done <- true
	t.ticker.Stop()
	t.status = TickerStatusStoped
	return nil
}

func (t *Ticker) Status() int {
	return t.status
}

func (t *Ticker) Start() error {
	if t.status == TickerStatusRunning {
		return ErrTickerIsStoped
	}

	go func(tt *Ticker) {
		for {
			select {
			case <-tt.done:
				return
			case <-tt.ticker.C:
				tt.fnc()
			}
		}
	}(t)
	t.status = TickerStatusRunning
	return nil
}

func New(interval int, fnc func()) Ticker {

	uuid := uuid.New().String()
	ticker := Ticker{
		uuid:   uuid,
		status: TickerStatusStoped,
		done:   make(chan bool),
		ticker: time.NewTicker(time.Duration(interval) * time.Second),
		fnc:    fnc,
	}

	return ticker
}
