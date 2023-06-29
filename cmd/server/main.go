package main

import (
	"laser-battle/pkg/domain"
	"laser-battle/pkg/events"
	"laser-battle/pkg/player"
	"time"
)

func main() {
	playerInfos := map[string][3]byte{
		"192.168.16.188:9000": {255, 0, 0},
		"192.168.16.185:9000": {0, 0, 255},
	}

	players := make(map[int]domain.Player, 0)
	i := 0
	for host, color := range playerInfos {
		p, err := player.New(color, host, 116, 79)
		if err != nil {
			panic(err)
		}
		players[i] = p
		i++
	}

	e := events.New("8080")

	ticker := time.NewTicker(time.Millisecond * 10)

	game := domain.New(players, e.C, ticker, e.SendScore)

	panic(game.Start())
}
