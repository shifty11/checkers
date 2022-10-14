mock-expected-keepers:
	mockgen -source=x/checkers/types/expected_keepers.go -destination=testutil/mock_types/expected_keepers.go
build-all:
	GOOS=linux GOARCH=amd64 go build -o ./build/checkersd-linux-amd64 ./cmd/checkersd/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build/checkersd-linux-arm64 ./cmd/checkersd/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./build/checkersd-darwin-amd64 ./cmd/checkersd/main.go

do-checksum:
	cd build && sha256sum checkersd-linux-amd64 checkersd-linux-arm64 checkersd-darwin-amd64 > checkers_checksum

build-with-checksum: build-all do-checksum