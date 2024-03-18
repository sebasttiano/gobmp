compile-validator:
	CGO_ENABLED=0 GOOS=$(OS) GO111MODULE=on go build -a -ldflags '-extldflags "-static"' -o ../../bin/validator ./validator.go
