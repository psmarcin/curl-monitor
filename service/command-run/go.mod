module commandrun

go 1.16

require (
	github.com/go-kit/kit v0.11.0
	github.com/streadway/amqp v1.0.0
	job v1.0.0
	common v1.0.0
)

replace (
	job => ../job
	result => ../result
	common => ../common
)
