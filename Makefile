createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root seckill_system

server:
	cd cmd/app && go run main.go

.PHONY: createdb server