package fmt

import (
	"fmt"

	"github.com/spf13/viper"
)

// Printf handles all general log messages for commands
func Printf(f string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(f, args...))
}

// Verbosef handles all log messages noted as verbose and does not show unless
// the verbose flag has been provided to the command
func Verbosef(f string, args ...interface{}) {
	if viper.GetBool("verbose") {
		fmt.Println(fmt.Sprintf(f, args...))
	}
}
