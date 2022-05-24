package core

type GamblerService interface {
	Find(coin, date string) (*GamblerEvent, error)
	Store(ge *GamblerEvent) error
}

type gamblerService struct {
	pipeline          GamblerPipeline
	gamblerRepository GamblerRepository
}

func NewGamblerService(gp GamblerPipeline, gr GamblerRepository) GamblerService {
	return &gamblerService{
		pipeline:          gp,
		gamblerRepository: gr,
	}
}

func (gs *gamblerService) Find(coin, date string) (*GamblerEvent, error) {
	return gs.gamblerRepository.Find(coin, date)
}

func (gs *gamblerService) Store(event *GamblerEvent) error {
	return gs.gamblerRepository.Store(event)
}
