package player

import (
	"fmt"
	"laser-battle/pkg/domain"
	"log"
	"net"
)

type player struct {
	conn net.Conn

	len    int
	center int

	color [3]byte
	mask  []bool

	enabled bool

	score int
}

func (p *player) Score() int {
	return p.score
}

func (p *player) GetColor() [3]byte {
	return p.color
}

func (p *player) Colorize(color [3]byte) {
	led := make([]byte, 0)
	for i := 0; i < p.len; i++ {
		led = append(led, color[:]...)
	}

	_, err := p.conn.Write(led)
	if err != nil {
		log.Println("write:", err)
	}
}

func New(color [3]byte, host string, len, center int) (domain.Player, error) {
	conn, err := net.Dial("udp", host)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	p := &player{
		color:  color,
		conn:   conn,
		len:    len,
		center: center,
	}

	p.Reset()

	return p, nil
}

func (p *player) SetLaserEnabled(enable bool) {
	p.enabled = enable
}

func (p *player) IsInCenter() bool {
	return !p.IsHeadInCenter() && p.mask[p.center]
}

func (p *player) IsHeadInCenter() bool {
	return p.mask[p.center] && !p.mask[p.center+1]
}

func (p *player) Step() {
	p.mask = p.mask[:len(p.mask)-1]
	p.mask = append([]bool{p.enabled}, p.mask...)
	p.send()
}

func (p *player) Reset() {
	p.mask = make([]bool, p.len)
	p.enabled = false
	p.score = 0
	p.send()
}

func (p *player) send() {
	led := make([]byte, 0)
	if p.enabled {
		p.score++
	}
	for i := 0; i < p.len; i++ {
		if p.mask[i] {
			led = append(led, p.color[:]...)
		} else {
			led = append(led, 0, 0, 0)
		}
	}

	_, err := p.conn.Write(led)
	if err != nil {
		log.Println("write:", err)
	}
}
