set GOOS=windows
set GOARCH=amd64
go build -o ./pyfract/fract.dll -buildmode=c-shared ./export
go build -o ./app/backend/presentation/main.exe ./app/backend/presentation
set GOOS=js
set GOARCH=wasm
go build -o ./app/backend/presentation/web/app.wasm ./app/frontend
set GOOS=linux
set GOARCH=amd64
go build -o ./pyfract/fract.so -buildmode=c-shared ./export
go build -o ./app/backend/presentation/main ./app/backend/presentation
set GOOS=windows
set GOARCH=amd64