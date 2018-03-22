package bus


import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// intervalTable is the schedule of backoff intervals in seconds
var intervalTable = []int{5, 5, 5, 15, 15, 15,
	30, 30, 30, 60, 60, 60, 90}
var intervals = byte(len(intervalTable) - 1)

// Backoff implements a basic increasing backoff strategy
// with a small amount of random jitter
type Backoff struct {
	interval byte
}

// NewBackoff returns a freshly initialized Backoff
func NewBackoff() *Backoff {
	return &Backoff{
		interval: 0,
	}
}

// Wait calculates the next backoff interval, sleeps
// for that amount and returns.
func (b *Backoff) Wait() {
	interval := b.nextInterval()
	log.Infof("Waiting %d seconds before reconnecting.", interval/time.Second)
	time.Sleep(interval)
}

// Reset restarts wait interval escalation
func (b *Backoff) Reset() {
	b.interval = byte(0)
}

func (b *Backoff) nextInterval() time.Duration {
	if b.interval < intervals {
		b.interval++
	}
	return time.Duration(intervalTable[b.interval]+b.jitter()) * time.Second
}

func (b *Backoff) jitter() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(4)
}
