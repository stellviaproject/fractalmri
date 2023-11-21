package msg

import (
	md "fractalmri/app/frontend/domain/model"
	"net/http"
)

const (
	ImageFormatErrStr   = "formato de imágen incorrecto"
	ExceededFileSizeStr = "tamaño de archivo excedido"
)

type Msg struct {
	Message string `json:"Message"`
	Assert  bool   `json:"Assert"`
	Data    any    `json:"Data"`
}

func NewUploadFailedMsg() (int, *Msg) {
	return http.StatusBadRequest, &Msg{
		Message: "falló la subida",
		Assert:  false,
	}
}

func NewFileSizeExceededMsg() (int, *Msg) {
	return http.StatusBadRequest, &Msg{
		Message: ExceededFileSizeStr,
		Assert:  false,
	}
}

func NewImageFormatErrMsg() (int, *Msg) {
	return http.StatusBadRequest, &Msg{
		Message: ImageFormatErrStr,
		Assert:  false,
	}
}

func NewImageNotFoundMsg() (int, *Msg) {
	return http.StatusNotFound, &Msg{
		Message: "la imágen no se pudo encontrar",
		Assert:  false,
	}
}

func NewFileNotFound() (int, *Msg) {
	return http.StatusNotFound, &Msg{
		Message: "el archivo no se pudo encontrar",
		Assert:  false,
	}
}

func NewEncodeErrMsg() (int, *Msg) {
	return http.StatusInternalServerError, &Msg{
		Message: "el servidor no ha podido codificar la imágen",
		Assert:  false,
	}
}

func NewImageUploadOk(img *md.ImageModel) (int, *Msg) {
	return http.StatusOK, &Msg{
		Message: "la imagen se ha subido",
		Assert:  true,
		Data:    img,
	}
}

func NewImageListOk(list []*md.ImageModel) (int, *Msg) {
	return http.StatusOK, &Msg{
		Message: "listado de imagenes satisfactorio",
		Assert:  true,
		Data:    list,
	}
}

func NewResultListOk(list []*md.ResultModel) (int, *Msg) {
	return http.StatusOK, &Msg{
		Message: "listado de resultados satisfactorio",
		Assert:  true,
		Data:    list,
	}
}

func NewRunFailedMsg() (int, *Msg) {
	return http.StatusBadRequest, &Msg{
		Message: "no se pudo ejecutar el análisis",
		Assert:  false,
	}
}

func NewResultMsg(result *md.ResultModel) (int, *Msg) {
	return http.StatusOK, &Msg{
		Message: "el análisis ha comenzado",
		Assert:  true,
		Data:    result,
	}
}

func NewInvalidURLParamMgs() (int, *Msg) {
	return http.StatusBadRequest, &Msg{
		Message: "los parámetros de la url son inválidos",
		Assert:  false,
	}
}

func NewImageNoExistMsg() (int, *Msg) {
	return http.StatusNotFound, &Msg{
		Message: "la imágen no existe",
		Assert:  false,
	}
}

func NewMFSDecodeErrMsg() (int, *Msg) {
	return http.StatusNotFound, &Msg{
		Message: "el espectro multifractal no se pudo decodificar",
		Assert:  false,
	}
}

func NewUnauthorizedMsg() (int, *Msg) {
	return http.StatusUnauthorized, &Msg{
		Message: "usuario no autorizado",
		Assert:  false,
	}
}

func NewProfileFailMsg() (int, *Msg) {
	return http.StatusNotFound, &Msg{
		Message: "fallo al decodificar profile",
	}
}

func NewUserIDAuthFailMsg() (int, *Msg) {
	return http.StatusUnauthorized, &Msg{
		Message: "autorizacion fallida",
		Assert:  false,
	}
}

func NewTokenAuthFailMsg() (int, *Msg) {
	return http.StatusUnauthorized, &Msg{
		Message: "autorizacion de token fallida",
		Assert:  false,
	}
}

func NewStartMsg(profile *md.Profile) (int, *Msg) {
	return http.StatusOK, &Msg{
		Message: "autenticacion satisfactoria",
		Assert:  true,
		Data:    profile,
	}
}
