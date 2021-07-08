package comtools

import (
	"regexp"
)

func GetMustCompileValue(RegularRule string, str string) []string {
	reg := regexp.MustCompile(RegularRule)
	return reg.FindStringSubmatch(str)
}

func GetAllMustCompileValue(RegularRule string, str string) [][]string {
	reg := regexp.MustCompile(RegularRule)
	return reg.FindAllStringSubmatch(str, -1)
}

//正则替换所有
func RegexpReplaceAll(regularRule, oldStr, replaceStr string) string {
	re, _ := regexp.Compile(regularRule)
	str := re.ReplaceAllString(oldStr, replaceStr)
	return str
}
