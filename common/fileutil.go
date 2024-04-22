package common

import (
	"os"
)

func ReadFile(filepath string) ([]byte, error) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
