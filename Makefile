build-warehouse:
	docker build -f build/warehouse/Dockerfile -t registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/warehouse --build-arg project=./cmd/warehouse .

publish-image-warehouse:
	docker push registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/warehouse

run-warehouse:
	docker run -p "4444:4444/udp" --name supplywatch_warehouse --rm  supplywatch_warehouse:latest

# docker build -f build/sensor/Dockerfile -t supplywatch_sensor --build-arg project=./cmd/sensor/ .

run-sensor:
	docker run --name supplywatch_sensor --rm  supplywatch_sensor_warehouse:latest

build-sensor:
	docker build -f build/sensor/Dockerfile -t registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/sensor --build-arg project=./cmd/sensor .

publish-image-sensor:
	docker push registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-gr/sensor

fmt:
	go fmt ./...

test:
	go test ./... -race -cover | grep -v "\[no test files\]"

grpc-gen:
	protoc -I proto/warehouse/ proto/warehouse/*.proto --go_out=internal/warehouse/grpc --go-grpc_out=internal/warehouse/grpc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

grpc-clean:
	rm grpc/pb/warehouse/*.pb.*
