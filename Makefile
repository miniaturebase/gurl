windows:
	@GOOS=windows GOARCH=amd64 go build -o bin/gurl.exe

mac:
	@GOOS=darwin GOARCH=amd64 go build -o bin/gurl

linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/gurl

# We don't buld 32-bit, but if we need to ..

# @GOOS=windows GOARCH=386 go build -o bin/gurl.exe
# @GOOS=darwin GOARCH=386 go build -o bin/gurl
# @GOOS=linux GOARCH=386 go build -o bin/gurl
