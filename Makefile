docker-run-dev:
	docker rmi -f saas-message || true
	docker build -t saas-message:dev .
	docker stop saas-message-dev || true
	docker run -d --rm -p 8888:8888 --name saas-message-dev saas-message:dev