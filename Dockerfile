FROM ubuntu:latest

RUN apt-get install golang-go libxrandr-dev libgl1-mesa-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev


RUN go get github.com/hajimehoshi/ebiten/v2
RUN go get -u github.com/gorilla/websocket
RUN go get github.com/google/uuid



CMD go run main.go


