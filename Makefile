monitor:
	docker build -f build/monitor/Dockerfile -t registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/monitor --build-arg project=./cmd/monitor .
	docker push registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/monitor

warehouse:
	docker build -f build/warehouse/Dockerfile -t registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/warehouse --build-arg project=./cmd/warehouse .
	docker push registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/warehouse

sensor:
	docker build -f build/sensor/Dockerfile -t registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/sensor --build-arg project=./cmd/sensor .
	docker push registry.code.fbi.h-da.de/distributed-systems/2020_vsprakt_moore/vs_ws21-22_mi5x-gimbel/mi5x-vfbi-003-team-g/sensor

fmt:
	go fmt ./...

test:
	go test ./... -race -cover | grep -v "\[no test files\]"

grpc-gen:
	protoc -I proto/warehouse/ proto/warehouse/*.proto --go_out=internal/warehouse/grpc --go-grpc_out=internal/warehouse/grpc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

grpc-clean:
	rm grpc/pb/warehouse/*.pb.*
