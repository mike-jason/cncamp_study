export tag=v0.1

build:

	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release:build
        echo "building httpserver container"
        docker build -t jiajia1003/cncamp:${tag} .

push:release
        echo "pushing jiajia1003/cncamp"
        docker push jiajia1003/cncamp:v0.1
