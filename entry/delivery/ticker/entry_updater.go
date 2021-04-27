package ticker

import (
	"cleanrss/domain"
	"log"
	"time"
)

type TickerEntryUpdater struct {
	u        domain.EntryUsecase
	t        *time.Ticker
	done     chan bool
	lastTick *time.Time
}

func NewTickerEntryUpdater(ticker *time.Ticker, u domain.EntryUsecase) TickerEntryUpdater {
	if ticker != nil {
		this := TickerEntryUpdater{t: ticker, u: u, done: make(chan bool)}
		this.Mount(nil)
		return this
	} else {
		return TickerEntryUpdater{done: make(chan bool)}
	}
}

func (t *TickerEntryUpdater) tick(time time.Time) {
	if t.lastTick == nil || time.Minute() == 0 {
		t.lastTick = &time
		log.Println("Refreshing all providers...")
		err := t.u.TriggerRefreshAll()
		if err != nil {
			log.Println(err)
		}
	}
}

func (t TickerEntryUpdater) Mount(ticker *time.Ticker) {
	if ticker != nil {
		t.t = ticker
	}
	go func(t *TickerEntryUpdater) {
		for {
			select {
			case <-t.done:
				return
			case tick := <-t.t.C:
				t.tick(tick)
			}
		}
	}(&t)
}

func (t TickerEntryUpdater) Unmount() {
	t.done <- true
}
