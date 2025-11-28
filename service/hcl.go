package service

import (
	"os"

	"github.com/genelet/horizon/dethcl"
)

func ReadHcl(_, path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	obj := map[string]any{}
	err = dethcl.Unmarshal(raw, &obj)
	return obj, err
}

func WriteHcl(fileName string, env map[string]any) error {

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

	out, err := dethcl.Marshal(env)
	if err != nil {
		return err
	}
	_, err = fh.Write(out)
	if err != nil {
		return err
	}

	return nil
}
