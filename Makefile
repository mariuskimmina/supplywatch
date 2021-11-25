build-warehouse:
	docker build -f Docker/warehouse/Dockerfile -t supplywatch_warehouse --build-arg project=./cmd/warehouse/ .

run-warehouse:
	docker run -p "4444:4444/udp" --name supplywatch_warehouse --rm  supplywatch_warehouse:latest

build-sensor:
	docker build -f Docker/sensor/Dockerfile -t supplywatch_sensor --build-arg project=./cmd/sensor-warehouse/ .

run-sensor:
	docker run --name supplywatch_sensor --rm  supplywatch_sensor_warehouse:latest

fmt:
	go fmt ./...

test:
	go test -race ./...
