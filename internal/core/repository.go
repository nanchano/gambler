package core

type GamblerRepository interface {
	Find(coin, date string) (*GamblerEvent, error)
	Store(event *GamblerEvent) error
}
