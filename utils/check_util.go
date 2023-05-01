package utils

import (
	"mime/multipart"
	"strconv"

	"github.com/h2non/filetype"
)

func CheckOffset(offset string) (int, error) {
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)

	if err != nil {
		return 0, err
	}

	return offsetInt, nil
}

func CheckId(id string) (int, error) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return idInt, nil
}


func CheckAudioFile(fh *multipart.FileHeader) (bool, error) {
	file, err := fh.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 261 bytes of the file to determine the file type
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return false, err
	}

	kind, _ := filetype.Match(head)
	if kind == filetype.Unknown {
		return false, nil
	}

	// Check if the detected file type is a sound file type
	if kind.MIME.Type == "audio" {
		return true, nil
	}

	return false, nil
}