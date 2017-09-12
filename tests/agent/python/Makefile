setup:
	pip install -Ur requirements.txt

start:
	docker-compose build
	docker-compose up

stop:
	docker-compose down -v

run:
	docker run -p 8000:8000 -e ELASTIC_APM_SERVERS=http://192.168.43.30:8200 24f1636b14a2
