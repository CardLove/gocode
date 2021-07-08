package comtools

import (
	"bufio"
	"checkManager/data"
	"checkManager/database"
	"checkManager/logger"
	"compress/gzip"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const LOG_PATH = "/var/log/"

const OS_CONF = "/etc/os-release"

var devInfoSet map[string]string
var UsbInfoSlice []data.UsbInfo

var indexInt int = 1

var MonthsMap = map[string]string{"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06", "Jul": "07", "Aug": "08", "Sept": "09", "Oct": "10", "Nov": "11", "Dec": "12"}
var FacturerNameMap = map[string]string{"HUAWEI": "华为", "APPLE": "苹果", "KINGSTON": "金士顿", "SANDISK": "闪迪", "SAMSUNG": "三星", "OPPO": "OPPO", "XIAOMI": "小米", "VIVO": "VIVO", "REALTEK": "瑞昱"}


func GetUosUsbInfoTest(usersStr string) {
	timeGen := ""
	idVendor := ""
	idProduct := ""
	Product := ""
	Manufacturer := ""
	SerialNumber := ""
	VendorName := ""
	ProductName := ""
	DeviceType := "其它设备"
	FacturerName := "其它厂商"
	re := regexp.MustCompile(`(?P<month>[a-zA-Z]{3,4})\s+(?P<day>\d+)\s+(?P<time>\d+:\d+:\d+)`)
	match := re.FindStringSubmatch(usersStr)
	month := MonthsMap[match[1]]
	day := match[2]
	if len(day) < 2 {
		day = "0" + day
	}
	aaaa := usersStr
	timeGen = time.Now().Format("2006") + "-" + month + "-" + day + " " + match[3]

	vidpidReg := regexp.MustCompile(`.*Device.*\(VID=(?P<vid>[a-zA-Z0-9]+)\s+and\s+PID=(?P<pid>[a-zA-Z0-9]+)\).*`)
	idVendormatch := vidpidReg.FindStringSubmatch(aaaa)
	if len(idVendormatch) == 3 {
		idVendor = idVendormatch[1]
		idProduct = idVendormatch[2]
	}

	ProductReg := regexp.MustCompile(`.*Product:\s(?P<Product>.*)\n`)
	Productmatch := ProductReg.FindStringSubmatch(aaaa)
	if len(Productmatch) != 0 {
		Product = Productmatch[1]
	}

	ManufacturerReg := regexp.MustCompile(`.*Direct-Access\s+(?P<Product>.*)`)
	Manufacturermatch := ManufacturerReg.FindStringSubmatch(aaaa)
	if len(Manufacturermatch) != 0 {
		Manufacturer = Manufacturermatch[1]
	}

	SerialNumberReg := regexp.MustCompile(`.*SerialNumber:\s(?P<SerialNumber>.*)`)
	SerialNumbermatch := SerialNumberReg.FindStringSubmatch(aaaa)
	if len(SerialNumbermatch) != 0 {
		SerialNumber = SerialNumbermatch[1]
	}


	UsbInfoSlice = append(UsbInfoSlice, data.UsbInfo{indexInt, idVendor, idProduct, Product,
		Manufacturer, SerialNumber, timeGen, timeGen, VendorName, ProductName, DeviceType, FacturerName})
	indexInt++

}
func GetUsbInfoTest(usersStr string) {
	timeGen := ""
	idVendor := ""
	idProduct := ""
	Product := ""
	Manufacturer := ""
	SerialNumber := ""
	VendorName := ""
	ProductName := ""
	DeviceType := "其它设备"
	FacturerName := "其它厂商"
	re := regexp.MustCompile(`(?P<month>[a-zA-Z]{3,4})\s+(?P<day>\d+)\s+(?P<time>\d+:\d+:\d+)`)
	match := re.FindStringSubmatch(usersStr)
	month := MonthsMap[match[1]]
	day := match[2]
	if len(day) < 2 {
		day = "0" + day
	}
	aaaa := usersStr
	timeGen = time.Now().Format("2006") + "-" + month + "-" + day + " " + match[3]

	idVendorReg := regexp.MustCompile(`.*idVendor=(?P<idVendor>[a-zA-Z0-9]+).*`)
	idVendormatch := idVendorReg.FindStringSubmatch(aaaa)
	if len(idVendormatch) != 0 {
		idVendor = idVendormatch[1]
	}

	idProductReg := regexp.MustCompile(`.*idProduct=(?P<idProduct>[a-zA-Z0-9]+).*`)
	idProductmatch := idProductReg.FindStringSubmatch(aaaa)
	if len(idProductmatch) != 0 {
		idProduct = idProductmatch[1]
	}

	ProductReg := regexp.MustCompile(`.*Product:\s(?P<Product>.*)\n`)
	Productmatch := ProductReg.FindStringSubmatch(aaaa)
	if len(Productmatch) != 0 {
		Product = Productmatch[1]
	}

	ManufacturerReg := regexp.MustCompile(`.*Manufacturer:\s(?P<Manufacturer>.*)\n`)
	Manufacturermatch := ManufacturerReg.FindStringSubmatch(aaaa)
	if len(Manufacturermatch) != 0 {
		Manufacturer = Manufacturermatch[1]
	}

	SerialNumberReg := regexp.MustCompile(`.*SerialNumber:\s(?P<SerialNumber>.*)`)
	SerialNumbermatch := SerialNumberReg.FindStringSubmatch(aaaa)
	if len(SerialNumbermatch) != 0 {
		SerialNumber = SerialNumbermatch[1]
	}
	if SerialNumber != "" {

			UsbInfoSlice = append(UsbInfoSlice, data.UsbInfo{indexInt, idVendor, idProduct, Product,
				Manufacturer, SerialNumber, timeGen, timeGen, VendorName, ProductName, DeviceType, FacturerName})
			indexInt++

	}

}

func GetUsbInfoVender(idVendor string, idProduct string) (string, string) {
	if idVendor != "" &&  idProduct != ""{
		content, _ := ReadAllIntoMemory("usb.ids.txt")
		arrystr := strings.Split(content, "\n")
		vendorReg := fmt.Sprintf(`^%v\s+(?P<Manufacturer>.*)`, idVendor)
		productReg := fmt.Sprintf(`\s+%v\s+(?P<Manufacturer>.*)`, idProduct)
		r, _ := regexp.Compile(vendorReg)
		r1, _ := regexp.Compile(productReg)
		r2, _ := regexp.Compile(`^\S+.*`)
		idVendName := ""
		idProductName := ""
		for index, val := range arrystr {
			arry := r.FindStringIndex(arrystr[index])
			if len(arry) != 0 {
				idVendormatch := r.FindStringSubmatch(val)
				idVendName = idVendormatch[1]
				for i := index + 1; i < len(arrystr); i++ {
					if r2.MatchString(arrystr[i]) == true {
						break
					}
					idProductmatch := r1.FindStringSubmatch(arrystr[i])
					if len(idProductmatch) != 0 {
						idProductName = idProductmatch[1]
						break
					}
				}
			}
			if idVendName != "" && idProductName != "" {
				return idVendName, idProductName
			}
		}
	}
	return "", ""

}
func GetUsbPhoneDevice() {

	mapt := GetUsbDeviceInfo()
	fmt.Println(mapt)
	for _, val := range mapt {
		if val.UsbInfoDeviceType == "手机" {
			database.OneDataWriteDB(data.UsbMobileInfo{0, val.UsbInfoIdVendor,
				val.UsbInfoIdProduct, val.UsbInfoFacturerName,
				val.UsbInfoProduct, val.UsbInfoSerialNumber,
				val.UsbInfoDeviceType, val.UsbInfoProductName,
				val.UsbInfoOldTime})
		}
	}
}

func isUos() bool {
	if Exists(OS_CONF) {
		cfg, err := ini.Load(OS_CONF)
		if err != nil {
			logger.GetIns().Error("Fail to read file: /etc/os-release", err)
			return false
		}
		sectionNames := cfg.SectionStrings()
		for _, param := range sectionNames {
			if param == "DEFAULT" {
				osName := cfg.Section("DEFAULT").Key("NAME").Value()
				if strings.Contains(strings.ToUpper(osName), "UOS") {
					return true
				}
			}
		}
	} else {
		logger.GetIns().Error("/etc/os-release file does not exist")
		return false
	}
	return false
}
func GetUosUsbDeviceInfo() {
	LogPathArr, _ := ListDir("/var/log", "messages")
	LogPathArr1, _ := ListDir("/var/log", "syslog")
	LogPathArr2, _ := ListDir("/var/log", "kern")
	devInfoSet = make(map[string]string, 0)
	LogPathArr = append(LogPathArr, LogPathArr1...)
	LogPathArr = append(LogPathArr, LogPathArr2...)

	UsbInfoSlice = make([]data.UsbInfo, 0)
	for _, path := range LogPathArr {
		if strings.Contains(path, "gz") {
			RunCmd(fmt.Sprintf("suRoot chmod o+r %s", path))
			gzipFile, err := os.Open(path)
			defer gzipFile.Close()
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
			scanner := bufio.NewScanner(gzipReader)
			scanner.Split(bufio.ScanLines)
			var tmpData []string
			tmpData = make([]string, 0)
			num := 0
			regBegin, _ := regexp.Compile(`.*USB Mass Storage device.*`)
			regEnd, _ := regexp.Compile(`.*udisksd.*Mounted.*`)
			flag := false
			for scanner.Scan() {
				if num == 10 {
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if regBegin.MatchString(scanner.Text()) {
					flag = true
					num++
					tmpData = append(tmpData, scanner.Text())
					continue
				}
				if regEnd.MatchString(scanner.Text()) {
					num++
					tmpData = append(tmpData, scanner.Text())
					aaa := strings.Join(tmpData, "\n")
					GetUosUsbInfoTest(aaa)
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if flag == true {
					num++
					tmpData = append(tmpData, scanner.Text())
				}
			}
		} else {
			RunCmd(fmt.Sprintf("suRoot chmod o+r %s", path))
			file, err := os.Open(path)
			if err != nil {
				log.Println(err)
				continue
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			var tmpData []string
			tmpData = make([]string, 0)
			num := 0
			regBegin, _ := regexp.Compile(`.*USB Mass Storage device.*`)
			regEnd, _ := regexp.Compile(`.*udisksd.*Mounted.*`)
			flag := false
			for scanner.Scan() {
				if num == 25 {
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if regBegin.MatchString(scanner.Text()) {
					flag = true
					num++
					tmpData = append(tmpData, scanner.Text())
					continue
				}
				if regEnd.MatchString(scanner.Text()) {
					num++
					tmpData = append(tmpData, scanner.Text())
					aaa := strings.Join(tmpData, "\n")
					GetUosUsbInfoTest(aaa)
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if flag == true {
					num++
					tmpData = append(tmpData, scanner.Text())
				}
			}
		}
	}
	for index, val := range UsbInfoSlice {
		UsbInfoSlice[index].UsbInfoVendorName, UsbInfoSlice[index].UsbInfoProductName = GetUsbInfoVender(val.UsbInfoIdVendor, val.UsbInfoIdProduct)
		if strings.Contains(UsbInfoSlice[index].UsbInfoProduct, "802.11") || strings.Contains(UsbInfoSlice[index].UsbInfoProductName, "802.11") ||
			strings.Contains(strings.ToUpper(val.UsbInfoProductName), "WIFI") {
			UsbInfoSlice[index].UsbInfoDeviceType = "无线网卡"
		} else if strings.Contains(strings.ToUpper(val.UsbInfoProduct), "DATATRAVELER") || strings.Contains(strings.ToUpper(val.UsbInfoManufacturer), "SANDISK") {
			UsbInfoSlice[index].UsbInfoDeviceType = "优盘"
		} else if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "PHONE") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "REDMI") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "NEXUS") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "PIXEL") {
			UsbInfoSlice[index].UsbInfoDeviceType = "手机"
		} else if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "BLUETOOTH") {
			UsbInfoSlice[index].UsbInfoDeviceType = "蓝牙设备"
		} else {
			UsbInfoSlice[index].UsbInfoDeviceType = "其它设备"
		}
		if UsbInfoSlice[index].UsbInfoManufacturer != "" {
			UsbInfoSlice[index].UsbInfoFacturerName = UsbInfoSlice[index].UsbInfoManufacturer
		} else {
			UsbInfoSlice[index].UsbInfoFacturerName = UsbInfoSlice[index].UsbInfoVendorName
		}
		for key, value := range FacturerNameMap {

			if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoManufacturer), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProduct), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoVendorName), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), key) {
				UsbInfoSlice[index].UsbInfoFacturerName = value
			}
		}
	}

}
func GetUsbDeviceInfo() []data.UsbInfo {
	if isUos() == true {
		GetUosUsbDeviceInfo()
		return UsbInfoSlice
	}
	LogPathArr, _ := ListDir("/var/log", "messages")
	LogPathArr1, _ := ListDir("/var/log", "syslog")
	LogPathArr2, _ := ListDir("/var/log", "kern")
	devInfoSet = make(map[string]string, 0)
	LogPathArr = append(LogPathArr, LogPathArr1...)
	LogPathArr = append(LogPathArr, LogPathArr2...)

	UsbInfoSlice = make([]data.UsbInfo, 0)
	for _, path := range LogPathArr {
		if strings.Contains(path, "gz") {
			RunCmd(fmt.Sprintf("suRoot suRun.sh chmod o+r %s", path))
			gzipFile, err := os.Open(path)
			defer gzipFile.Close()
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
			scanner := bufio.NewScanner(gzipReader)
			scanner.Split(bufio.ScanLines)
			var tmpData []string
			tmpData = make([]string, 0)
			num := 0
			regBegin, _ := regexp.Compile(`.*New USB device found.*`)
			regEnd, _ := regexp.Compile(`.*SerialNumber:.*`)
			flag := false
			for scanner.Scan() {
				if num == 10 {
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if regBegin.MatchString(scanner.Text()) {
					flag = true
					num++
					tmpData = append(tmpData, scanner.Text())
					continue
				}
				if regEnd.MatchString(scanner.Text()) {
					num++
					tmpData = append(tmpData, scanner.Text())
					aaa := strings.Join(tmpData, "\n")
					GetUsbInfoTest(aaa)
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if flag == true {
					num++
					tmpData = append(tmpData, scanner.Text())
				}
			}
		} else {
			RunCmd(fmt.Sprintf("suRoot suRun.sh chmod o+r %s", path))
			file, err := os.Open(path)
			if err != nil {
				log.Println(err)
				continue
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			var tmpData []string
			tmpData = make([]string, 0)
			num := 0
			regBegin, _ := regexp.Compile(`.*New USB device found.*`)
			regEnd, _ := regexp.Compile(`.*SerialNumber:.*`)
			flag := false
			for scanner.Scan() {
				if num == 10 {
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if regBegin.MatchString(scanner.Text()) {
					flag = true
					num++
					tmpData = append(tmpData, scanner.Text())
					continue
				}
				if regEnd.MatchString(scanner.Text()) {
					num++
					tmpData = append(tmpData, scanner.Text())
					aaa := strings.Join(tmpData, "\n")
					GetUsbInfoTest(aaa)
					tmpData = make([]string, 0)
					num = 0
					flag = false
				}
				if flag == true {
					num++
					tmpData = append(tmpData, scanner.Text())
				}
			}
		}
	}
	for index, val := range UsbInfoSlice {
		UsbInfoSlice[index].UsbInfoVendorName, UsbInfoSlice[index].UsbInfoProductName = GetUsbInfoVender(val.UsbInfoIdVendor, val.UsbInfoIdProduct)
		if strings.Contains(UsbInfoSlice[index].UsbInfoProduct, "802.11") || strings.Contains(UsbInfoSlice[index].UsbInfoProductName, "802.11") ||
			strings.Contains(strings.ToUpper(val.UsbInfoProductName), "WIFI") {
			UsbInfoSlice[index].UsbInfoDeviceType = "无线网卡"
		} else if strings.Contains(strings.ToUpper(val.UsbInfoProduct), "DATATRAVELER") || strings.Contains(strings.ToUpper(val.UsbInfoManufacturer), "SANDISK") {
			UsbInfoSlice[index].UsbInfoDeviceType = "优盘"
		} else if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "PHONE") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "REDMI") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "NEXUS") ||
			strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "PIXEL") {
			UsbInfoSlice[index].UsbInfoDeviceType = "手机"
		} else if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), "BLUETOOTH") {
			UsbInfoSlice[index].UsbInfoDeviceType = "蓝牙设备"
		} else {
			UsbInfoSlice[index].UsbInfoDeviceType = "其它设备"
		}
		if UsbInfoSlice[index].UsbInfoManufacturer != "" {
			UsbInfoSlice[index].UsbInfoFacturerName = UsbInfoSlice[index].UsbInfoManufacturer
		} else {
			UsbInfoSlice[index].UsbInfoFacturerName = UsbInfoSlice[index].UsbInfoVendorName
		}
		for key, value := range FacturerNameMap {

			if strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoManufacturer), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProduct), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoVendorName), key) ||
				strings.Contains(strings.ToUpper(UsbInfoSlice[index].UsbInfoProductName), key) {
				UsbInfoSlice[index].UsbInfoFacturerName = value
			}
		}
	}
	return UsbInfoSlice
}
