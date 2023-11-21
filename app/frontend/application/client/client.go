package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var hostName string

var UseHTTPS = false

func Init(URL *url.URL) {
	hostName = URL.Hostname() + ":" + URL.Port()
}

func urlFor(URL string) string {
	if strings.HasPrefix(URL, "http://") || strings.HasPrefix(URL, "https://") {
		return URL
	}
	var base string
	if UseHTTPS {
		base = "https://" + hostName
	} else {
		base = "http://" + hostName
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(URL, "/")
}

type Client struct {
	client *http.Client
}

func New() *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}
	c := new(Client)
	c.client = &http.Client{
		Jar: jar,
	}
	return c
}

func (c *Client) get(URL string, send any) error {
	res, err := c.client.Get(URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(bytes))
	return json.Unmarshal(bytes, send)
}

func (c *Client) post(URL string, send, receive any) error {
	buffer := &bytes.Buffer{}
	bytes, err := json.Marshal(send)
	if err != nil {
		return err
	}
	if _, err = buffer.Write(bytes); err != nil {
		return err
	}
	res, err := c.client.Post(URL, "application/json", buffer)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bytes, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, receive)
}

// Envia una imagen al cliente acorde a los parametros especificados.
// Los parametros que se pueden especificar son id, type y name.
// La url toma la forma: "/image?id=<id>&name=<name>&type=<type>".
// Los valores validos de id son los numeros enteros cuyos valores esten asignados a los identificadores de las imagenes subidas.
// Los valores validos de name son cualquier combinacion de caracteres menos caracteres vacios porque al agregar las imagenes se ignoran.
// Los valores validos de type son "r", "g", "b", "gray", "mask", "image" y "", donde "" o no especificar type es equivalente a usar "image".
func (c *Client) upload(URL string, file app.Value, receive any) error {
	fileBytes, err := ReadFileBytes(file)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(fileBytes)
	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Get("name").String()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}
	writer.Close()
	pURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", URL, formData)
	if err != nil {
		return err
	}
	req.URL = pURL
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.client.Do(req)
	if err != nil {
		log.Panicln(err)
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(data))
	return json.Unmarshal(data, receive)
}

func (c *Client) Upload(path string, file app.Value, receive any) error {
	vReceive := reflect.ValueOf(receive)
	URL := urlFor(path)
	log.Printf("POST %s\n", URL)
	strType := message_struct(receive)
	value := reflect.New(strType)
	if err := c.upload(URL, file, value.Interface()); err != nil {
		return err
	}
	value = value.Elem()
	message := value.FieldByName("Message").String()
	if !value.FieldByName("Assert").Bool() {
		return fmt.Errorf(message)
	} else {
		log.Println(message)
	}
	vRes := value.FieldByName("Data")
	if !vRes.IsNil() {
		// En caso de que se devuelvan datos es necesario que la estructura para recibir los datos sea distinta de nil y un puntero
		if vReceive.IsNil() || vReceive.Kind() != reflect.Ptr {
			//Retornar el error indicando que la interfaz para recibir los datos no es un puntero
			return fmt.Errorf("expected receive as a pointer, got %v", vReceive.Type().String())
		}
		vReceive.Elem().Set(vRes.Elem()) //Establecer los valores
	}
	return nil
}

// Envia los datos de send a una URL y lee la respuesta en receive.
func (c *Client) Post(path string, send, receive any) error {
	vReceive := reflect.ValueOf(receive) //Obtener el valor de receive
	URL := urlFor(path)                  //Obtener la URL absoluta
	log.Printf("POST %s\n", URL)
	strType := message_struct(receive)                           //Obtener el tipo de dato de la estructura para recibir la información
	value := reflect.New(strType)                                //Instanciar la estructura
	if err := c.post(URL, send, value.Interface()); err != nil { //Enviar los datos en send y recibir la información en la estructura creada
		return err //Retornar el error para indicar que ocurrio
	}
	value = value.Elem()
	message := value.FieldByName("Message").String() //Obtener el mensaje del servidor
	if !value.FieldByName("Assert").Bool() {         //Comprobar si hubo error en la operacion
		return fmt.Errorf(message) //Retornar el error en un mensaje
	} else {
		log.Println(message) //Imprimir el mensaje
	}
	vRes := value.FieldByName("Data") //Obtener los datos de la operacion
	if !vRes.IsNil() {                //Comprobar si se devolvieron datos al ser distinto de nil
		//En caso de que se devuelvan datos es necesario que la estructura para recibir los datos sea distinta de nil y un puntero
		if vReceive.IsNil() || vReceive.Kind() != reflect.Ptr {
			//Retornar el error indicando que la interfaz para recibir los datos no es un puntero
			return fmt.Errorf("expected receive as a pointer, got %v", vReceive.Type().String())
		}
		vReceive.Elem().Set(vRes.Elem()) //Establecer los valores
	}
	return nil
}

// Obtiene los datos de una URL en receive.
func (c *Client) Get(path string, receive any) error {
	vReceive := reflect.ValueOf(receive)
	URL := urlFor(path)
	log.Printf("GET %s\n", URL)
	strType := message_struct(receive)
	value := reflect.New(strType)
	if err := c.get(URL, value.Interface()); err != nil {
		return err
	}
	value = value.Elem()
	message := value.FieldByName("Message").String()
	if !value.FieldByName("Assert").Bool() {
		return fmt.Errorf(message)
	} else {
		log.Println(message) //Imprimir el mensaje
	}
	vRes := value.FieldByName("Data")
	if !vRes.IsNil() {
		if vReceive.IsNil() || vReceive.Kind() != reflect.Ptr {
			return fmt.Errorf("expected data as a pointer, got %v", vReceive.Type().String())
		}
		vReceive.Elem().Set(vRes.Elem())
	}
	return nil
}

func (c *Client) uploadfile(URL string, file string, receive any) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(fileBytes)

	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}
	writer.Close()
	pURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", URL, formData)
	if err != nil {
		return err
	}
	req.URL = pURL
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.client.Do(req)
	if err != nil {
		log.Panicln(err)
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, receive)
}

/*
Crea el tipo de dato de la estructura para recibir los mensajes del servidor.

Le da al campo data el mismo tipo que la estructura deseada por el usuario.

Esto permite que el mensaje que no es de error se descarte y que el mensaje de error sea puesto en la forma de un error con fmt.Errorf(...).
*/
func message_struct(receive any) reflect.Type {
	return reflect.StructOf([]reflect.StructField{
		{
			//Mensaje enviado por el servidor.
			Name: "Message",
			Type: reflect.TypeOf(""),
			Tag:  `json:"message"`,
		},
		{
			//Si se ha realizado la operacion satisfactoriamente
			Name: "Assert",
			Type: reflect.TypeOf(true),
			Tag:  `json:"assert"`,
		},
		{
			//Información resultante de la operacion
			Name: "Data",
			Type: reflect.TypeOf(receive),
			Tag:  `json:"data"`,
		},
	})
}

func ReadFileBytes(file app.Value) ([]byte, error) {
	// Crear un canal de Go para recibir los bytes del archivo
	bytesChan := make(chan []byte)
	defer close(bytesChan)

	file.Call("arrayBuffer").Then(func(v app.Value) {
		data := app.Window().Get("Uint8Array").New(v)
		bytes := make([]byte, data.Length())
		app.CopyBytesToGo(bytes, data)
		bytesChan <- bytes
	})
	// Esperar a que los bytes del archivo sean enviados al canal de Go
	select {
	case bytes := <-bytesChan:
		return bytes, nil
	case <-time.After(time.Second * 5):
		return nil, fmt.Errorf("timeout reading file")
	}
}

func (c *Client) Download(path string) ([]byte, error) {
	res, err := c.client.Get(urlFor(path))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, err
}
