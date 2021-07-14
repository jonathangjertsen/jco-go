// +build release

package buildconfig

import "fmt"

func PanicHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
