/*
Copyright © 2023 John Lennard <john@yakmoo.se>
*/

package service

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

// parseFile wrapper around the file parser
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

// ReadEnv reads the environment file in .env format in the order .env.local, .env, .env.<environment>, .env.<environment>.local
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

// WriteEnv writes the environment file in .env format
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

// WriteJSON writes the environment file in JSON format
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
