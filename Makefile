build-warehouse:
	docker build -f Docker/warehouse/Dockerfile -t supplywatch_warehouse --build-arg project=./cmd/warehouse/ .

run-warehouse:
	docker run -p "4444:4444/udp" --name supplywatch_warehouse --rm  supplywatch_warehouse:latest

build-sensor-warehouse:
	docker build -f Docker/sensor-warehouse/Dockerfile -t supplywatch_sensor_warehouse --build-arg project=./cmd/sensor-warehouse/ .

run-sensor-warehouse:
	docker run --name supplywatch_sensor_warehouse --rm  supplywatch_sensor_warehouse:latest
