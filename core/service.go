package core

// GamblerService represents the main behaviour for our use case: retrieve and store cryptocurrency info
type GamblerService interface {
	Run(dates ...string)
	Find(coin, date string) (*GamblerEvent, error)
	Store(event *GamblerEvent) error
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

// Run transforms the responses from the pipeline into GamblerEvents
func (service *gamblerService) Run(dates ...string) {
	service.pipeline.Run(service.repo, dates...)
}

// Find retrieves a GamblerEvent from the repository for a given coin and date
func (service *gamblerService) Find(coin, date string) (*GamblerEvent, error) {
	return service.repo.Find(coin, date)
}

// Store saves each GamblerEvent on a channel on the given repository
func (service *gamblerService) Store(event *GamblerEvent) error {
	return service.repo.Store(event)
}
