build-warehouse:
	docker build -f build/warehouse/Dockerfile -t supplywatch_warehouse --build-arg project=./cmd/warehouse/ .

run-warehouse:
	docker run -p "4444:4444/udp" --name supplywatch_warehouse --rm  supplywatch_warehouse:latest

build-sensor:
	docker build -f Docker/sensor/Dockerfile -t supplywatch_sensor --build-arg project=./cmd/sensor-warehouse/ .

run-sensor:
	docker run --name supplywatch_sensor --rm  supplywatch_sensor_warehouse:latest

fmt:
	go fmt ./...

test:
	go test ./... -race -cover | grep -v "\[no test files\]"

grpc-gen:
	protoc -I proto/warehouse/ proto/warehouse/*.proto --go_out=internal/warehouse/grpc --go-grpc_out=internal/warehouse/grpc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

grpc-clean:
	rm grpc/pb/warehouse/*.pb.*
