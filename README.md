# FractalMRI

## Descripción
Este es un proyecto para analizar y clasificar imágenes de resonancia magnética utilizando algoritmos de geometría fractal.

El algoritmo utilizado para analizar las imágenes es el espectro multifractal de exponentes de hölder y conteo de cajas (box-counting dimension). La clasificación de imágenes se realiza por medio de KNN (N Vecinos más cercanos).

## Herramientas Utilizadas

Para desarrollarlo se utilizaron las siguientes herramientas:

- Lenguaje de programación: Golang, Python y C.
- Entorno de desarrollo: Visual Studio Code.
- Gestión de versiones: Git.

## Instalación

- Instalar Golang versión 1.21.
- Instalar Python versión 3.11.
- Configurar el entorno de programación de Golang.
- Clonar el repositorio en la ruta GOPATH.

## Compilar la biblioteca

- Moverse al directrio principal del proyecto.
- Ejecutar el compilador de Golang:

  Para Windows: `GOOS=windows GOARCH=amd64 go build -o ./pyfract/fract.dll -buildmode=c-shared ./export`
  
  Para Linux: `GOOS=linux GOARCH=amd64 go build -o ./pyfract/fract.so -buildmode=c-shared ./export`

## Compilar el visor
- Compilar el frontend con:
  
  `GOOS=js GOARCH=wasm go build -o ./app/backend/presentation/web/app.wasm ./app/frontend`
- Compilar el backend con:

  Para Windows:
  `GOOS=windows GOARCH=amd64 go build -o ./app/backend/presentation/main.exe ./app/backend/presentation`

  Para Linux:
  `GOOS=linux GOARCH=amd64 go build -o ./app/backend/presentation/main ./app/backend/presentation`