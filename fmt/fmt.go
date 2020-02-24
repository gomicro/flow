package fmt

import (
	"fmt"

	"github.com/spf13/viper"
)

func Printf(f string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(f, args...))
}

func Verbosef(f string, args ...interface{}) {
	if viper.GetBool("verbose") {
		fmt.Println(fmt.Sprintf(f, args...))
	}
}
