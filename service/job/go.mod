module job

go 1.16

require (
	github.com/go-kit/kit v0.11.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.2
	common v1.0.0
)

replace (
	common => ./../common
)
