sudo systemctl restart docker.service && docker run -p 27017:27017 mongo

go build ./cmd/main.go
go run ./cmd/main.go
