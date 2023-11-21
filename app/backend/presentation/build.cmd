go build -o ./main.exe .
set GOOS=js
set GOARCH=wasm
go build -o web/app.wasm ../../frontend