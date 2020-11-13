sudo systemctl restart docker.service && docker run -d -p 27017:27017 mongo

nohup go run ./cmd/main.go &
