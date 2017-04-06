all: linux windows osx

linux:
	GOOS=linux GOARCH=386 go build -o srf-fwu-linux *.go

windows:
	GOOS=windows GOARCH=386 go build -o srf-fwu-win.exe *.go

osx:
	GOOS=darwin GOARCH=386 go build -o srf-fwu-osx *.go

clean:
	rm -f srf-fwu-linux srf-fwu-win.exe srf-fwu-osx
