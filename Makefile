run:
	go run main.go

mocks:
	mockery --all --keeptree --dir=internal/repository --output=internal/repository/mocks --case underscore
	mockery --all --keeptree --dir=internal/service --output=internal/service/mocks --case underscore

test:
	go test -v -coverprofile cover.out ./internal/...
	go tool cover -html cover.out -o cover.html 
