setup:
	pip install -Ur requirements.txt

start:
	docker-compose build
	docker-compose up

stop:
	docker-compose down -v

run:
	ELASTIC_APM_SERVERS=http://localhost:8200
