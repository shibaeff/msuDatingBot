sudo systemctl restart docker.service && docker run -d -p 27017:27017 -v echobot:/data/db mongo
go build ./cmd/main.go
nohup go run ./cmd/main.go &
