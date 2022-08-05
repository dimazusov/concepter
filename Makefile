swagger: swagger-init
	swag init -g ./internal/server/http/router.go -o api

mockgen-install:
	go install github.com/golang/mock/mockgen@v1.6.0

mockgen: mockgen-install
	mockgen -source=internal/domain/thinker/thinker.go -package=thinker	-destination=internal/domain/thinker/thinker_mock.go
	mockgen -source=internal/domain/inference/repository.go -package=inference	-destination=internal/domain/inference/repository_mock.go
	mockgen -source=internal/domain/inference/service.go -package=inference	-destination=internal/domain/inference/service_mock.go

