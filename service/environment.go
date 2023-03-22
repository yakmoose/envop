package service

import (
	"os"

	"github.com/subosito/gotenv"
)

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

	env := gotenv.Env{}
	for _, fileName := range fileNames {
		fh, err := os.Open(fileName)
		if err != nil {
			continue
		}
		defer fh.Close()
		e := gotenv.Parse(fh)

		// boop
		for k, v := range e {
			env[k] = v
		}
	}
	return env
}
