package utils

import (
	"mime/multipart"
	"strconv"

	"github.com/h2non/filetype"
)


func CheckOffset(offset string) int {
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)

	if err != nil || offsetInt < 0{
		return -1
	}

	return offsetInt
}

func CheckId(id string) int {
	idInt, err := strconv.Atoi(id)

	if err != nil || idInt < 0{
		return -1
	}

	return idInt
}


func CheckAudioFile(fh *multipart.FileHeader) bool {
	file, err := fh.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 261 bytes of the file to determine the file type
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return false
	}

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false
	}

	// Check if the detected file type is a sound file type
	if kind.MIME.Type == "audio" {
		return true
	}

	return false
}