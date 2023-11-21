cls
set GOOS=windows
set GOARCH=amd64
go build -o ./main.exe .
set GOOS=js
set GOARCH=wasm
go build -o web/app.wasm ../../frontend
main.exe
set GOOS=windows
set GOARCH=amd64