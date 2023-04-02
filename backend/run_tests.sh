docker image build --tag backend:v1.0.0
docker image build ../dispatcher --tag dispatcher:v1.0.0
docker-compose up -d --remove-orphans --scale backend=2
go clean -testcache
go test -v ./test/integration_test.go
docker-compose down