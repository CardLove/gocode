package comtools

import (
	"checkManager/logger"
	"fmt"
	"github.com/beevik/etree"
	"os"
	"path/filepath"
	"strings"
)

const (
	TARGRT_PATH = "/tmp/temp/ofd/"
)

func GetOfdFileContent(filePath string) string {
	content := ""
	//判断临时文件的路径是否存在 不存在创建
	if !Exists(TARGRT_PATH) {
		err := os.MkdirAll(TARGRT_PATH, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	ret, _ := RunCmd(fmt.Sprintf("suRoot suRun.sh 7z  x  %s -r  -o%s  -aoa", filePath, TARGRT_PATH))
	if !strings.Contains(ret, "Everything is Ok") {
		return content
	}
	pathSlice := GetXMLAllPath(TARGRT_PATH, ".xml")
	for _, value := range pathSlice {
		ret := GetXmlContent(value)
		if len(ret) != 0 {
			tempStr := strings.Join(ret, "\n")
			content = content + tempStr
		}

	}
	//清空临时目录
	os.RemoveAll(TARGRT_PATH)
	return content

}

//获取指定文件夹下的指定文件类型的文件路径
func GetXMLAllPath(dir string, fileType string) []string {
	filePathSlice := make([]string, 0)
	err  := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if IsFile(path) {
			if strings.Contains(GetFileSuffix(path), fileType) {
				filePathSlice = append(filePathSlice, path)
				//fmt.Println(path)
			}
		}
		return nil
	})
	if err != nil {
		logger.GetIns().Debug("err:",err)
	}
	return filePathSlice
}

//得到xml的内容
func GetXmlContent(xmlFile string) []string {
	retCon := make([]string, 0)
	ret, err := ReadAllIntoMemory(xmlFile)
	if err != nil {
		fmt.Println(err)
		return retCon
	}
	doc := etree.NewDocument()
	_ = doc.ReadFromString(ret)

	elemEnablePath := "./ofd:Page/ofd:Content/ofd:Layer/ofd:TextObject/ofd:TextCode/"

	path, _ := etree.CompilePath(elemEnablePath)

	//找不到的节点为空
	for _, value := range doc.FindElementsPath(path) {
		retCon = append(retCon, value.Text())
	}

	return retCon
}
