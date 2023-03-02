docker-run-dev:
	docker rmi -f saas-message || true
	docker build -t saas-message:dev .
	docker run -d --rm -p 8888:8888 saas-message:dev