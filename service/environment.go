package service

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

func parseFile(file string, env *map[string]string) error {
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fh.Close()
	for k, v := range gotenv.Parse(fh) {
		(*env)[k] = v
	}
	return nil
}

func ReadEnv(environment, path string) map[string]string {
	// read the .environment file
	// .environment.local .environment .environment.<environment> .environment.<environment>.local
	// create 1password items and store em
	fileNames := []string{
		path + ".local",
		path,
		path + "." + environment,
		path + "." + environment + ".local",
	}

	env := make(map[string]string, 0)
	for _, fileName := range fileNames {
		parseFile(fileName, &env)
	}
	return env
}

func WriteEnv(environment, path string, env map[string]string) error {
	fh, err := os.Create(path + "." + environment + ".local")
	if err != nil {
		return err
	}
	defer fh.Close()

	for k, v := range env {
		fh.WriteString(k + "=" + strconv.Quote(v) + "\n")
	}

	return nil
}

func WriteJSON(environment, path string, env map[string]string) error {
	fh, err := os.Create(path + "." + environment + ".local.json")
	if err != nil {
		return err
	}
	defer fh.Close()

	out, err := json.Marshal(env)
	if err != nil {
		return err
	}
	fh.Write(out)

	return nil
}
