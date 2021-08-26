module trigger

go 1.16

require (
	common v1.0.0
	github.com/go-kit/kit v0.11.0
	github.com/google/uuid v1.3.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40 // indirect
	github.com/streadway/amqp v1.0.0
	job v1.0.0
)

replace (
	common => ../common
	job => ../job
)
