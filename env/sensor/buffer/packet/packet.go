// Package packet provides implementation of packaet for the buffer in a sensor
package packet

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"
)

var PACKET_COPIES_ERROR = errors.New("PACKET_COPIES_ERROR")

var (
	number_of_copies uint
	packet_id        int
	mux              sync.Mutex
)

func init() {
	flag.UintVar(&number_of_copies, "num-copies", 8, "maximum number of copies of a packate in a buffer")
}

// Packet is generated by the sensor and stored in the buffer
type Packet struct {
	copies    int
	timestamp time.Time
	parent_id int
	Id        int
}

// New packet generation through the sensor
func New(parent_id int) Packet {
	mux.Lock()
	defer mux.Unlock()

	pkt := Packet{
		copies:    int(number_of_copies),
		parent_id: parent_id,
		timestamp: time.Now(),
		Id:        packet_id,
	}
	packet_id++

	return pkt
}

// Zero zeros the packet
// making packet zero means that packet no longer exists
func (p *Packet) Zero() {
	p.copies = 0
}

// Exists if the packet is not zeroed
func (p Packet) Exists() bool {
	return p.copies != 0
}

// DecreaseCopies after getting ACK from the reciever with n_prime
func (p *Packet) DecreaseCopies(n_prime int) error {
	if p.copies > n_prime {
		p.copies -= n_prime
		return nil
	}
	return fmt.Errorf("[%w] n_prime should be smaller then packet copies", PACKET_COPIES_ERROR)
}

// Deliverable checks if the packet can further be sent to other bicycles
func (p Packet) Deliverable() bool {
	return p.copies > 1
}

// SetCopies of the packet in the buffer
func (p *Packet) SetCopies(n int) error {
	if n > 0 {
		p.copies = n
		return nil
	}
	return fmt.Errorf("[%w] number of copies should be positive", PACKET_COPIES_ERROR)
}

func (p Packet) GetCopies() int {
	return p.copies
}

// GetTimestamp of creation of the packet
func (p Packet) GetTimestamp() time.Time {
	return p.timestamp
}

// GetParentId of the packet, the sensor which generated the packet
func (p Packet) GetParentId() int {
	return p.parent_id
}
