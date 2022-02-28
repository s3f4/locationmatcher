## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## mock: creates mock files
mock:
	cd internal/driverlocation && mockery --all && \
	cd ../matching && mockery --all

## up: docker compose up
up:
	docker-compose -f docker-compose.yaml up -d  --build --remove-orphans 

## up: docker compose down
down:
	docker-compose -f docker-compose.yaml down 

