package core

// GamblerService represents the main behaviour for our use case: retrieve and store cryptocurrency info
type GamblerService interface {
	Extract(dates ...string) <-chan PipelineResponse
	Process(responses <-chan PipelineResponse) <-chan *GamblerEvent
	Find(coin, date string) (*GamblerEvent, error)
	Store(events <-chan *GamblerEvent) error
}

// gamblerService contains a pipeline and repository to implement the GamblerService interface
type gamblerService struct {
	pipeline GamblerPipeline
	repo     GamblerRepository
}

// NewGamblerService builds a new gamblerService with a given pipeline and repo
func NewGamblerService(pipeline GamblerPipeline, repo GamblerRepository) GamblerService {
	return &gamblerService{
		pipeline: pipeline,
		repo:     repo,
	}
}

func (service *gamblerService) Extract(dates ...string) <-chan PipelineResponse {
	return service.pipeline.Extract(dates...)
}

func (service *gamblerService) Process(responses <-chan PipelineResponse) <-chan *GamblerEvent {
	return service.pipeline.Process(responses)
}

// Find retrieves a GamblerEvent from the repository for a given coin and date
func (service *gamblerService) Find(coin, date string) (*GamblerEvent, error) {
	return service.repo.Find(coin, date)
}

// Store saves each GamblerEvent on a channel on the given repository
func (service *gamblerService) Store(events <-chan *GamblerEvent) error {
	return service.repo.Store(events)
}
