package images

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func CompressImage(reader io.Reader) ([]byte, error) {
	imageSrc, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	newImg := image.NewRGBA(imageSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), imageSrc, imageSrc.Bounds().Min, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imageSrc, imageSrc.Bounds().Min, draw.Over)
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 30})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
