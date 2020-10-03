.PHONY: dev
dev:
	cd app && \
	HTTPPORT=8888 \
	gin -i -p 3002 run main.go
