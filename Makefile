
all:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-d -w -s -extldflags "-static"' slackit.go
	docker build -t antonipx/slackit:latest .
