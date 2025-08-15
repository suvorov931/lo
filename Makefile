start-app:
	 docker compose up --build -d lo-service

down:
	docker compose down

all: start-app