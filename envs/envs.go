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
		splits := strings.SplitN(s.Text(), "=", 2)
		if len(splits) == 2 {
			key := splits[0]
			value := splits[1]
			if key != "" && value != "" {
				envs = append(envs, Env{key, value})
			}
		}
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("envs: failed scanning file: %v", err.Error())
	}

	return envs, err
}
