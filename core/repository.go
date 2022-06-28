package core

// GamblerRepository implementations must be able to Find data
// for a coin/date and store a GamblerEvent
type GamblerRepository interface {
	Find(coin, date string) (*GamblerEvent, error)
	Store(event *GamblerEvent) error
}
