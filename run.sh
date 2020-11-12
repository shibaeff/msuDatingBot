alias restart_docker="sudo systemctl restart docker.service"
alias runmongo="docker run -p 27017:27017 mongo"

sudo restart_docker && runmongo

go build ./cmd/main.go
go run ./cmd/main.go
