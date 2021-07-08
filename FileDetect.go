package comtools

import (
	"bufio"
	"checkManager/logger"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !isDir(path)
}

//得到指定目录中所有文件夹
func GetAllFolders(pathname string) ([]string, error) {
	var dir []string
	_, err := os.Stat(pathname)
	if err != nil {
		return dir, err
	}

	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			dir = append(dir, fi.Name())
		} else {
		}
	}
	return dir, err
}

// 读取文件的所有内容
func ReadAllIntoMemory(filename string) (content string, err error) {
	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return "", err
	}
	defer fp.Close()

	fileInfo, err := fp.Stat()
	if err != nil {
		return "", err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = fp.Read(buffer) // 文件内容读取到buffer中
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

// 获取指定目录中的所有文件
func GetDirAllFile(dir string) []string {
	filePath := make([]string,0)
	err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				//fmt.Println("dir:", path)
				return nil
			}
			filePath = append(filePath, path)
			return nil
		})
	if err != nil {
		logger.GetIns().Debug("filepath.Walk  err : ", err)
	}
	return filePath
}

// 逐行读取, 一行是一个[]byte, 多行就是[][]byte
func ReadByLine(filename string) (lines [][]byte, err error) {
	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	bufReader := bufio.NewReader(fp)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		} else {
			lines = append(lines, line)
		}
	}

	return
}
var wg sync.WaitGroup
//获取指定目录中文件的个数
func GetDirFileNum(dir string,lock *sync.RWMutex, ct *int)    error {
	defer wg.Done()
	num := 0
	err:= filepath.Walk(dir, func(dir string, f os.FileInfo, err error) error {
		if f == nil {
			lock.Lock()
			*ct  = *ct + num
			lock.Unlock()
			return err
		}
		if !f.IsDir() {
			num++
		}
		return nil
	})
	if err != nil {
		logger.GetIns().Debug("err:",err)
	}
	lock.Lock()
	*ct  = *ct + num
	lock.Unlock()
	return nil

}
//得到文件列表中所有文件个数
func GetPathFileSum(fileList [] string ,lock *sync.RWMutex,sum *int )   {
	var  sumTemp  int
	wg.Add(len(fileList))
	for _ ,value :=  range fileList {
		go GetDirFileNum(value,lock,&sumTemp )
	}
	wg.Wait()
	lock.Lock()
	*sum = sumTemp
	lock.Unlock()
}
//类型转化获取百分比
func  PercentageCon(pos ,sum  int ) string {
	return  strconv.Itoa(int (float32(pos )/ float32(sum )  *10000 ))
}

//得到文件名字
func GetFileName(filePath string) string {
	var filenameWithSuffix, fileSuffix, filenameOnly string
	filenameWithSuffix = path.Base(filePath)
	fileSuffix = path.Ext(filenameWithSuffix)
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly

}

//得到文件扩展名
func GetFileSuffix(filePath string) string {
	var filenameWithSuffix, fileSuffix string
	filenameWithSuffix = path.Base(filePath)
	fileSuffix = path.Ext(filenameWithSuffix)
	return fileSuffix
}

//得到文件扩展名  myself
func GetFileExtenName(filePath string) string {
	filenameWithSuffix := path.Base(filePath)
	return filenameWithSuffix[strings.Index(filenameWithSuffix, "."):]

}

//获取文件修改时间 时间字符串
func GetFileModTime(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
	}
	return fi.ModTime().In(time.Local).Format("2006-01-02 15:04:05")
}

//获取文件创建时间 时间字符串
func GetFileCreateTime(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
	}
	stat_t := fi.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat_t.Ctim.Sec), int64(stat_t.Ctim.Nsec)).In(time.Local).Format("2006-01-02 15:04:05")
}

//获取文件大小
func GetFileSize(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return "0"
	}
	return strconv.FormatInt(fi.Size(), 10)
}

//获取文件所有者
func GetFileOwner(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return "root"
	}
	stat := fi.Sys().(*syscall.Stat_t)
	li, err := user.LookupId(strconv.FormatInt(int64(stat.Uid), 10))
	if err != nil {
		return "root"
	}
	return li.Username
}

func GetSoftwareInstallTime(software string) string {
	LogPathArr, _ := ListDir("/var/log", "dpkg.log")
	data := "0000-00-00 00:00:00"
	timeFormat := `(?P<time>\d+-\d+-\d+\s+\d+:\d+:\d+)\s+(install|upgrade)\s+%s:`
	timeStr := fmt.Sprintf(timeFormat, software)
	reg := regexp.MustCompile(timeStr)
	for _, path := range LogPathArr {
		if strings.Contains(path, "gz") {
			gzipFile, err := os.Open(path)
			if err != nil {
				log.Println(err)
				continue
			}
			gzipReader, err := gzip.NewReader(gzipFile)
			if err != nil {
				log.Println(err)
				continue
			}
			defer gzipReader.Close()
			rd, err := ioutil.ReadAll(gzipReader)
			if err != nil {
				log.Println(err)
				continue
			}
			content := string(rd)
			regStringList := reg.FindAllStringSubmatch(content, -1)
			if len(regStringList) != 0 {
				for _, val := range regStringList {
					if data < val[1] {
						data = val[1]
					}
				}
			}
		} else {
			content, _ := ReadAllIntoMemory(path)
			regStringList := reg.FindAllStringSubmatch(content, -1)
			if len(regStringList) != 0 {
				for _, val := range regStringList {
					if data < val[1] {
						data = val[1]
					}
				}
			}
		}
	}
	return data
}

func GetFileModifyTime(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return time.Now().String()[0:19]
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().String()[0:19]
	}
	dataTime := fi.ModTime().String()[0:19]
	return dataTime
}

//获取文件的创建时间
func GetFileCrtime(filePath string) string {
	//var monthsMap = map[string]string{"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06", "Jul": "07", "Aug": "08", "Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12"}

	ret, _ := RunCmd(fmt.Sprintf(`suRoot suRun.sh stat %s`, filePath))
	crtime := GetMustCompileValue(`\s*Change:\s*(.*)\.\s*`, ret)
	//fmt.Println(crtime)
	//fmt.Println(len(crtime))

	//Tue Oct 20 00:42:04 2020
	fileCrtime := ""
	if len(crtime) == 2 {
		fileCrtime = crtime[1]
	}
	return fileCrtime
}
