build:
	go build -o bin/service-accounts services/accounts/main.go

certs:
	openssl genrsa -out certs/cert.pem 2048

test:
	go test ./core/auth/...