module commandrun

go 1.16

require (
	github.com/go-kit/kit v0.11.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	job v1.0.0
	result v1.0.0
)

replace (
	job => ../job
	result => ../result
)
