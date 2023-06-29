package domain

import (
	"time"
)

type Player interface {
	SetLaserEnabled(enable bool)

	IsInCenter() bool
	IsHeadInCenter() bool

	Step()
	Reset()

	GetColor() [3]byte
	Colorize([3]byte)

	Score() int
}

type Game struct {
	players   map[int]Player
	events    chan Event
	ticker    *time.Ticker
	sendScore func(score int, winnerId int)
}

func New(players map[int]Player, events chan Event, ticker *time.Ticker, sendScore func(score int, winnerId int)) *Game {
	return &Game{players: players, events: events, ticker: ticker, sendScore: sendScore}
}

type Event struct {
	PlayerId int  `json:"player_id,omitempty"`
	Enabled  bool `json:"enabled,omitempty"`
}

func (g *Game) Start() error {
	for {
		select {
		case <-g.ticker.C:
			for _, player := range g.players {
				player.Step()
			}

			someBodyInCenter := false
			for _, player := range g.players {
				if player.IsInCenter() {
					someBodyInCenter = true
					break
				}
			}

			if someBodyInCenter {
				for id, player := range g.players {
					if player.IsHeadInCenter() {
						score := 0
						for _, player := range g.players {
							score += player.Score()
						}
						g.sendScore(score, id)
						color := player.GetColor()
						time.Sleep(time.Second / 2)
						for _, player := range g.players {
							player.Colorize(color)
						}
						time.Sleep(time.Second)
						for _, player := range g.players {
							player.Reset()
						}
						break
					}
				}
			}

		case event := <-g.events:
			g.players[event.PlayerId].SetLaserEnabled(event.Enabled)
		}
	}
	return nil
}
