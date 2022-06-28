package core

// PipelineResponse implementations are responses from any API Pipeline.
// Must convert itself into a GamblerEvent
type PipelineResponse interface {
	Convert() *GamblerEvent
}

// GamblerPipeline implementations must Run for a given set of dates,
// including extracting, processing and storing events on a given repository
type GamblerPipeline interface {
	Run(repo GamblerRepository, dates ...string)
}
