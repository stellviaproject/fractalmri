GOOS=windows GOARCH=amd64 go build -o ../pyfract/fract.dll -buildmode=c-shared -ldflags="-s -w" .
#GOOS=linux GOARCH=amd64 go build -o ../pyfract/fract.so -buildmode=c-shared .
#CC=x86_64-w64-mingw32-gcc GOOS=linux GOARCH=amd64 go build -o ../pyfract/fract.so -buildmode=c-shared .
