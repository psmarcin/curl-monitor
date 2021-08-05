package common

const (
	CommandRunRoutingKey = "commandrun"
	CommandRunQueueName  = "CM.Job.CommandRun"
	ResultRoutingKey     = "result"
	ResultQueueName      = "CM.Job.Result"
	TriggerRoutingKey    = "trigger"
	TriggerQueueName     = "CM.Job.Trigger"
	JobExchangeName      = "CM.Job"
)
