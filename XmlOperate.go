package comtools

import (
	"fmt"
	"io/ioutil"
)

func GetXmlTenplate(name string) string {
	f, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}
