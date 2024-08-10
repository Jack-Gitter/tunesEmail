docker-start: 
	docker compose --profile tunesemailservice up --build
docker-stop:
	docker compose --profile tunesemailservice down
