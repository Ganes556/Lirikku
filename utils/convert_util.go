package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"

	"github.com/Lirikku/libs"
)

func Audio2RawBase64(data *multipart.FileHeader) (string, error) {

	dat, _:= data.Open()
	
	defer dat.Close()

	buff, _ := io.ReadAll(dat)

	cmd := libs.FfmpegCmd
	
	r := bytes.NewReader(buff)

	var rawOuput bytes.Buffer
	
	err := cmd.PipeInput(r).OutputFormat("s16le").AudioCodec("pcm_s16le").AudioChannels(1).AudioRate(44100).PipeOutput(&rawOuput).Run()

	if err != nil {
		return "", err
	}

	rawBase64 := base64.StdEncoding.EncodeToString(rawOuput.Bytes())
	
	return rawBase64, nil
	
}