# GOOS：darwin、freebsd、linux、windows
# GOARCH：386、amd64、arm、s390x

all: darwin

publish: darwin linux windows arm64 armv7

# for raspbian
armv7:
	CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm go build -o bin/srun-armv7 ./cmd/srun

# for raspberry with arm64 ubuntu
arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/srun-arm64 ./cmd/srun

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/srun ./cmd/srun

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/srun ./cmd/srun

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/srun.exe ./cmd/srun

clean:
	rm -rf ./bin

.PHONY:publish
.PHONY:darwin
.PHONY:linux
.PHONY:windows
.PHONY:armv7
.PHONY:arm64
.PHONY:clean
