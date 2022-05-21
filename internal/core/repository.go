package core

type GamblerRepository interface {
	Find(exchange, coin, date string) (*GamblerEvent, error)
	Store(ge *GamblerEvent) error
}
