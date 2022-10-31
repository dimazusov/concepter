swagger: swagger-init
	swag init -g ./internal/server/http/router.go -o api

mockgen-install:
	go install github.com/golang/mock/mockgen@v1.6.0

mockgen: mockgen-install
	mockgen -source=concepter.go -destination=concepter_mock.go -package=main

