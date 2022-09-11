package config

import (
	"fmt"
	"io/ioutil"
)

func ReadFile(f string) ([]byte, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("unable to load specified file %s: %s", f, err)
	}
	return data, nil
}
