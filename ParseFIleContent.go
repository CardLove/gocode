package comtools

import (
	"strings"
)

func GetFileData(fileName, RegularRule string) [][]string {

	var fileDate [][]string
	FileContent, err := ReadAllIntoMemory(fileName)
	if err != nil {
		return fileDate
	}
	contentArr := strings.Split(FileContent, "\n")
	for _, value := range contentArr {
		if value == "" {
			continue
		}
		dateLine := make([]string, 0)
		parts := strings.SplitN(value, RegularRule, -1)
		//fmt.Println(len(parts))
		//fmt.Println(parts)
		for i := 0; i < len(parts); i++ {
			dateLine = append(dateLine, parts[i])
		}
		fileDate = append(fileDate, dateLine)
	}

	return fileDate
}
