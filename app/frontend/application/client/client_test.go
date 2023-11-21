package client

import (
	"bytes"
	"image/png"
	"net/url"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	URL, _ := url.Parse("http://127.0.0.1:8080")
	Init(URL)
	c := New()
	buffer, err := c.Download("/image?id=0")
	if err != nil {
		t.FailNow()
	}
	reader := bytes.NewReader(buffer)
	img, err := png.Decode(reader)
	if err != nil {
		t.FailNow()
	}
	file, err := os.OpenFile("./image.png", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm|os.ModeDevice)
	if err != nil {
		t.FailNow()
	}
	defer file.Close()
	if err = png.Encode(file, img); err != nil {
		t.FailNow()
	}
}
