module job

go 1.16

require (
	common v1.0.0
	github.com/go-kit/kit v0.11.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40 // indirect
	github.com/lib/pq v1.10.2
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
)

replace common => ./../common
