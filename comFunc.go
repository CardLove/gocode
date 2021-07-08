package comtools

import "C"
import (
	"checkManager/logger"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

func GoStrings(argc C.int, argv **C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

func charpp2string(charpp **C.char, n int) []string {
	var b *C.char
	ptrSize := unsafe.Sizeof(b)
	gostring := make([]string, n)
	if n > 0 {
		for i := 0; i < n; i++ {
			element := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(charpp)) + uintptr(i)*ptrSize))
			gostring[i] = C.GoString((*C.char)(*element))
		}
		//(*C.char)(*(**C.char)(unsafe.Pointer(uintptr( unsafe.Pointer(job.exHosts)) + uintptr(1)*ptrSize ) ) ))
	}
	return gostring
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.GetIns().Error("%#v", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) error {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.GetIns().Error("%#v", err)
		return err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		logger.GetIns().Error("%#v", err)
		return err
	}
	return nil
}

func TestWriteJsonContentToFile(filePath, fileContent string) {
	err := ioutil.WriteFile(filePath, []byte(fileContent), 0666)
	if err != nil {
		logger.GetIns().Debug("写入%s失败:%s", filePath, fileContent)
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetString(FileFD *os.File, pos *int64, count int) string {
	//FileFD, err := os.OpenFile(FindDevName, os.O_RDONLY, os.ModeDevice)
	//if err != nil {
	//	logger.GetIns().Error("failed open:%#v, err:%#v", FindDevName, err)
	//	return ""
	//}
	//defer FileFD.Close()
	content := make([]byte, count)
	FileFD.Seek(*pos, io.SeekStart)
	n, err := FileFD.Read(content)
	if err != nil {
		logger.GetIns().Error("Read:%#v, %#v, %#v, err:%#v", FileFD.Name(), *pos, count, err)
		return ""
	}
	*pos = int64(n) + *pos
	retString := make([]byte, n)
	//copy(content,retString)
	retString = append(retString, content[:n]...)
	return string(retString)
}

func GetPartitionSize(dev string) (int64, error) {
	RunCmd(fmt.Sprintf(`suRoot suRun.sh chmod go+r %s`, dev))
	//logger.GetIns().Debug("FindDevName", dev)
	file, err := os.OpenFile(dev, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		fmt.Printf("error opening %s: %s\n", dev, err)
		return 0, err
	}
	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Printf("error seeking to end of %s: %s\n", dev, err)
		return 0, err
	}
	//fmt.Printf("%s is %d bytes.\n", dev, pos)
	file.Seek(0, io.SeekStart)
	file.Close()
	return pos, err
}

//获取磁盘列表总任务个数
func GetTaskNumber(devList []string, pageSize int64) int64 {
	var taskNumber int64 = 0
	for _, devName := range devList {
		if !Exists(devName) {
			continue
		}
		diskSize, err := GetPartitionSize(devName)
		if err != nil {
			continue
		}
		taskNumber = taskNumber + diskSize/pageSize
		if diskSize%pageSize != 0 {
			taskNumber = taskNumber + 1
		}
	}
	return taskNumber
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配前缀过滤。
func ListDir(dirPth string, prefix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	prefix = strings.ToUpper(prefix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasPrefix(strings.ToUpper(fi.Name()), prefix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, nil
}

//func RunCmd(cmdStr string) (content string, err error) {
//	content = ""
//	//list := strings.Split(cmdStr, " ")
//	//cmd := exec.Command(list[0],list[1:]...)
//	cmd := exec.Command("bash", "-c", cmdStr)
//	var out bytes.Buffer
//	var stderr bytes.Buffer
//	cmd.Stdout = &out
//	cmd.Stderr = &stderr
//	err = cmd.Run()
//	//fmt.Println("Process PID:", cmd.Process.Pid)
//	//pid = strconv.Itoa(cmd.Process.Pid)
//	content = out.String()
//	return
//}
