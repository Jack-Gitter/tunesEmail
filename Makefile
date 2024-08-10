docker-start: 
	docker compose --profile tunesemailservice up --build -d
docker-stop:
	docker compose --profile tunesemailservice down
