docker-run-dev:
	docker rmi -f saas-message || true
	docker build -t saas-message:dev .
	docker stop saas-message-dev || true
	docker rm saas-message-dev || true
	docker run -d --rm -p 8888:8888 -e PORT=8888 -e PRODUCTION=false --name saas-message-dev saas-message:dev
docker-dev-logs:
	docker logs -f saas-message-dev
docker-run-prod:
	docker rmi -f saas-message || true
	docker build -t saas-message:prod .
	docker stop saas-message-prod || true
	docker rm saas-message-prod || true
	docker run -d --rm -p 8888:8888 -e PORT=8888 -e PRODUCTION=true --name saas-message-prod saas-message:prod