up:
	docker-compose up --build -d
down:
	docker-compose down
conndb:
	docker exec -it my_postgres psql -U user -d mydatabase