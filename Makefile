createdb:
	 docker exec -it postgres12 createdb --username=root --owner=root seckill_system

.PHONY: createdb