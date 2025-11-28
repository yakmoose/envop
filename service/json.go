package service

import (
	"encoding/json"
	"os"
)

func ReadJson(_, path string) (map[string]any, error) {

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	obj := map[string]any{}
	err = json.Unmarshal(raw, &obj)
	return obj, err
}

// WriteJSON writes the environment file in JSON format
func WriteJSON(fileName string, env map[string]any) error {
	var fh *os.File
	var err error

	if fileName == "" {
		fh = os.Stdout
	} else {
		fh, err = os.Create(fileName)
		if err != nil {
			return err
		}
		defer fh.Close()
	}

	out, err := json.Marshal(env)
	if err != nil {
		return err
	}
	
	_, err = fh.Write(out)
	if err != nil {
		return err
	}

	return nil
}
