package core

// Generic Repository, must be able to Find data for a coin/date, and store GamblerEvents
type GamblerRepository interface {
	Find(coin, date string) (*GamblerEvent, error)
	Store(events <-chan *GamblerEvent) error
}
