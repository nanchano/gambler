package core

type GamblerRepository interface {
	Find(coin, date string) (*GamblerEvent, error)
	Store(events <-chan *GamblerEvent) error
}
