package envs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env represents a key and value of an environment variable to use when
// defining configs.
type Env struct {
	Key   string
	Value string
}

// ParseFile takes a file to look in for env var definitions
func ParseFile(file string) ([]Env, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("envs: failed opening file: %v", err.Error())
	}
	defer fh.Close()

	envs := []Env{}

	s := bufio.NewScanner(fh)
	for s.Scan() {
		key, value := split(s.Text())
		if key != "" && value != "" {
			envs = append(envs, Env{key, value})
		}
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("envs: failed scanning file: %v", err.Error())
	}

	return envs, err
}

// ParseSlice takes a slice of strings capable of making up key value pairs and
// returns split keys and values. It will ignore entries it cannot do anything
// with.
func ParseSlice(slice []string) []Env {
	envs := []Env{}

	for _, kp := range slice {
		key, value := split(kp)
		if key != "" && value != "" {
			envs = append(envs, Env{key, value})
		}
	}

	return envs
}

func split(s string) (string, string) {
	splits := strings.SplitN(s, "=", 2)
	if len(splits) == 2 {
		return splits[0], splits[1]
	}

	return "", ""
}
