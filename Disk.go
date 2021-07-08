package comtools

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"log"
	"strings"
	"syscall"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
func GetDiskSize(filePath string) (size uint64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(filePath, &fs)
	if err != nil {
		return 0
	}
	return fs.Blocks * uint64(fs.Bsize)

}

//  得到当前系统磁盘分区
func GetCurrentOsAllPartition() []string {
	var partitionList []string
	fdiskContent, _ := RunCmd(`suRoot suRun.sh fdisk -l`)
	diskRet := GetAllMustCompileValue(`(/dev/sd[a-z][0-9]+)\s+\*?\s*[0-9]+\s*[0-9]+\s*[0-9]+\s*.*[G|M]{1}\s+`, fdiskContent)
	for _, value := range diskRet {
		log.Println(value[1])
		partitionList = append(partitionList, value[1])
	}
	//fmt.Println("partitionList :", partitionList)
	return partitionList
}

//将磁盘转到挂在点
func GetpartirionMounton(partition []string) []string {
	mountonList := make([]string, 0)

	for _, value := range partition {
		ret, _ := RunCmd(fmt.Sprintf("suRoot suRun.sh df  -h |grep %s  |awk '{print $6}'", value))
		if ret != "" {
			if strings.Compare(strings.TrimSpace(ret), "/") == 0 { //根目录只检查这个几个目录
				mountonList = append(mountonList, "/home/")
				//mountonList = append(mountonList, "/mnt/")
				mountonList = append(mountonList, "/opt/")
				//mountonList = append(mountonList, "/home/")
				//mountonList = append(mountonList, "/mnt/")
				//mountonList = append(mountonList, "/opt/")
				//mountonList = append(mountonList, "/root/")
				mountonList = append(mountonList, "/root/")

			} else {
				mountonList = append(mountonList, strings.TrimSpace(ret))
			}
		} else { //检查工具临时挂载的文件不存在创建挂载
			if Exists("/tmp/WLHMount") {
				RunCmd(fmt.Sprintf("suRoot suRun.sh mkdir   /tmp/WLHMount%s  -p", value))
				RunCmd(fmt.Sprintf("suRoot suRun.sh  mount  %s  /tmp/WLHMount%s ", value, value))
				mountonList = append(mountonList, fmt.Sprintf("/tmp/WLHMount%s ", value))

			} else {
				RunCmd("suRoot suRun.sh mkdir   /tmp/WLHMount  -p")
				RunCmd(fmt.Sprintf("suRoot suRun.sh mkdir   /tmp/WLHMount%s  -p", value))
				RunCmd(fmt.Sprintf("suRoot suRun.sh mount  %s  /tmp/WLHMount%s ", value, value))
				mountonList = append(mountonList, fmt.Sprintf("/tmp/WLHMount%s ", value))

			}

		}
	}
	return mountonList
}

//获取系统中挂在的磁盘名字和挂载点
func GetPartitionNameAndPartitionMounton() [][]string {
	partitions := make([][]string, 0)
	infos, _ := disk.Partitions(false)
	for _, info := range infos {
		if info.Fstype == "squashfs" {
			continue
		}
		partitions = append(partitions, []string{info.Device, info.Mountpoint})
	}
	return partitions
}

//磁盘大小格式化
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
