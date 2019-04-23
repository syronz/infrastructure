package validator

import (
	"strconv"
	"errors"
)

func CheckId(v string) (int64, error){
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	if id < 0 {
		return 0, errors.New("ID is not eligible")
	}

	return id, nil
}
