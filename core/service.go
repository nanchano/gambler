package core

// GamblerService represents the main behaviour for our use case;
// retrieve and store cryptocurrency info
type GamblerService interface {
	Run(dates ...string)
	Find(coin, date string) (*GamblerEvent, error)
	Store(event *GamblerEvent) error
}

// gamblerService contains a pipeline and repository to implement
// the GamblerService interface
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

// Run extracts, processes and stores crytpo data for a set of dates
func (service *gamblerService) Run(dates ...string) {
	service.pipeline.Run(service.repo, dates...)
}

// Find retrieves a GamblerEvent from the repository for a given coin and date
func (service *gamblerService) Find(coin, date string) (*GamblerEvent, error) {
	return service.repo.Find(coin, date)
}

// Store saves a GamblerEvent on the a repository
func (service *gamblerService) Store(event *GamblerEvent) error {
	return service.repo.Store(event)
}
