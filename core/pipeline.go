package core

// PipelineResponse implementations are responses from any API Pipeline. Must convert itself into a GamblerEvent
type PipelineResponse interface {
	Convert() *GamblerEvent
}

// GamblerPipeline implementations must extract data for a set of dates, and process the responses into GamblerEvents
type GamblerPipeline interface {
	Extract(dates ...string) <-chan PipelineResponse
	Process(responses <-chan PipelineResponse) <-chan *GamblerEvent
}
