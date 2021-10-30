build-warehouse:
	docker build -f Docker/warehouse/Dockerfile -t supplywatch_warehouse --build-arg project=./cmd/warehouse/ .

run-warehouse:
	docker run -p "4444:4444/udp" --name supplywatch_warehouse --rm  supplywatch_warehouse:latest
