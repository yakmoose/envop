/*
Copyright © 2025 John Lennard <john@yakmoo.se>
*/

package service

import (
	"os"
	"strconv"

	"github.com/hashicorp/go-envparse"
)

// parseFile wrapper around the file parser
func parseFile(path string, env *map[string]any) error {
	var fh *os.File
	var err error
	if path == "" {
		fh = os.Stdin
	} else {
		fh, err = os.Open(path)
		if err != nil {
			return nil
		}
		defer fh.Close()
	}

	parsedEnvfile, err := envparse.Parse(fh)
	if err != nil {
	}

	for k, v := range parsedEnvfile {
		(*env)[k] = v
	}
	return nil
}

// ReadEnv reads the environment file in .env format in the order .env.local, .env, .env.<environment>, .env.<environment>.local
func ReadEnv(envName, path string) (map[string]any, error) {
	// read the .environmentName file
	// .env.local .env .env.<environmentName> .env.<environmentName>.local

	var fileNames []string
	if path == "" {
		fileNames = []string{""}
	} else {
		fileNames = []string{
			path,
			path + ".local",
		}

		if envName != "" {
			fileNames = append(
				fileNames, path+"."+envName,
				path+"."+envName+".local",
			)
		}
	}
	env := make(map[string]any, 0)
	for _, fileName := range fileNames {
		err := parseFile(fileName, &env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

// WriteEnv writes the environment file in .env format
func WriteEnv(fileName string, env map[string]any) error {
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

	for k, v := range env {
		_, err = fh.WriteString(k + "=" + strconv.Quote(anyToStringish(v)) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
