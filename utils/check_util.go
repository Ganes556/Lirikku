package utils

import "strconv"

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
