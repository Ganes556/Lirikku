package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/Lirikku/libs"
)

func Audio2RawBase64(data *multipart.FileHeader) string {

	dat, _ := data.Open()

	defer dat.Close()

	buff, _ := io.ReadAll(dat)

	cmd := libs.FfmpegCmd

	r := bytes.NewReader(buff)

	var rawOuput bytes.Buffer

	err := cmd.PipeInput(r).
		OutputFormat("s16le").
		AudioCodec("pcm_s16le").
		AudioChannels(1).
		AudioRate(44100).
		PipeOutput(&rawOuput).
		Run()

	if err != nil {
		return ""
	}

	rawBase64 := base64.StdEncoding.EncodeToString(rawOuput.Bytes())

	return rawBase64

}

func Int2String(i int) string {
	return strconv.Itoa(i)
}

func Convert2Json(data any) string {
	dataJson, _ := json.Marshal(data)
	return string(dataJson)
}

func ConvertCapitalize(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToLower(s)
	r := []rune(s)
	return string(unicode.ToUpper(r[0])) + string(r[1:])
}

func ConvertUrl2Normal(s string) string {
	decodedURL, _ := url.QueryUnescape(s)
	return decodedURL
}

func Convert2Map(d any) map[string]any{
	var data map[string]any
	dd, _ := json.Marshal(d)
	json.Unmarshal(dd, &data)
	return data
}