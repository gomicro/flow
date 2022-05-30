package config

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	confDir  = ".aws"
	confFile = "config"
)

type Profile struct {
	Name               string
	AwsAccessKeyID     string `profile:"aws_access_key_id"`
	AwsAccountID       string `profile:"aws_account_id"`
	AwsSecretAccessKey string `profile:"aws_secret_access_key"`
	Color              string `profile:"color"`
	Output             string `profile:"output"`
	Region             string `profile:"region"`
	RoleArn            string `profile:"role_arn"`
	SourceProfile      string `profile:"source_profile"`
	SsoAccountID       string `profile:"sso_account_id"`
	SsoRegion          string `profile:"sso_region"`
	SsoRoleName        string `profile:"sso_role_name"`
	SsoStartURL        string `profile:"sso_start_url"`
}

func GetProfiles() ([]*Profile, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("config: get home directory: %w", err)
	}

	file, err := os.Open(filepath.Join(usr.HomeDir, confDir, confFile))
	if err != nil {
		return nil, fmt.Errorf("config: opening file: %w", err)
	}
	defer file.Close()

	t := reflect.TypeOf(Profile{})

	pTagFieldMap := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		pt := f.Tag.Get("profile")

		if pt != "" {
			pTagFieldMap[pt] = f.Name
		}
	}

	ps := []*Profile{}
	var p *Profile

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(l, "[") {
			p = &Profile{}
			ps = append(ps, p)

			p.Name = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.Trim(l, "[]")), "profile"))
		} else {
			parts := strings.SplitN(l, "=", 2)

			if len(parts) != 2 {
				continue
			}

			tag, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

			field, ok := pTagFieldMap[tag]
			if ok {
				reflect.ValueOf(p).Elem().FieldByName(field).SetString(value)
			}
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("config: scanning config: %w", err)
	}

	return ps, nil
}
