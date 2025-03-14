linux:
	go build -tags x11
	./ping-pong
win:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_LDFLAGS="-static-libgcc -static -lpthread" go build
	wine ./ping-pong.exe
