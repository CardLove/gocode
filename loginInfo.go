package comtools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WtmpUserLogin struct {
	Time     string
	User     string
	Type     string
	Exit     int16
	Terminal string
	Host     string
}

const (
	Empty        int = iota // Record does not contain valid info (formerly known as UT_UNKNOWN on Linux)
	RunLevel         = iota // Change in system run-level (see init(8))
	BootTime         = iota // Time of system boot (in ut_tv)
	NewTime          = iota // Time after system clock change (in ut_tv)
	OldTime          = iota // Time before system clock change (in ut_tv)
	InitProcess      = iota // Process spawned by init(8)
	LoginProcess     = iota // Session leader process for user login
	UserProcess      = iota // Normal process
	DeadProcess      = iota // Terminated process
	Accounting       = iota // Not implemented
)

func MarshalJSON(u int) string {
	switch u {
	case Empty:
		return "Empty"
	case RunLevel:
		return "RunLevel"
	case BootTime:
		return "BootTime"
	case NewTime:
		return "NewTime"
	case OldTime:
		return "OldTime"
	case InitProcess:
		return "InitProcess"
	case LoginProcess:
		return "LoginProcess"
	case UserProcess:
		return "UserProcess"
	case DeadProcess:
		return "DeadProcess"
	case Accounting:
		return "Accounting"
	default:
		return ""
	}
}

//[1] [00000] [~~  ] [shutdown] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-01T08:47:33,240737+0000]
//[2] [00000] [~~  ] [reboot  ] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-01T08:47:37, 579227+0000]
//[1] [00053] [~~  ] [runlevel] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-01T08:47:41, 356903+0000]
//[7] [02271] [    ] [alvin   ] [:0          ] [:0                  ] [0.0.0.0        ] [2020-10-01T08:49:01, 844802+0000]
//[2] [00000] [~~  ] [reboot  ] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-03T01:44:03, 355734+0000]
//[1] [00053] [~~  ] [runlevel] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-03T01:44:08, 809741+0000]
//[7] [04639] [    ] [alvin   ] [:0          ] [:0                  ] [0.0.0.0        ] [2020-10-03T01:44:44, 732214+0000]
//[2] [00000] [~~  ] [reboot  ] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-09T01:10:29, 014040+0000]
//[1] [00053] [~~  ] [runlevel] [~           ] [5.4.0-48-generic    ] [0.0.0.0        ] [2020-10-09T01:10:34, 083325+0000]
//[7] [04239] [    ] [alvin   ] [:0          ] [:0                  ] [0.0.0.0        ] [2020-10-09T01:38:52, 453475+0000]
//[5] [95446] [tty4] [        ] [tty4        ] [                    ] [0.0.0.0        ] [2020-10-12T09:11:08, 277591+0000]
//[6] [95446] [tty4] [LOGIN   ] [tty4        ] [                    ] [0.0.0.0        ] [2020-10-12T09:11:08, 277591+0000]
//[8] [95446] [tty4] [        ] [tty4        ] [                    ] [0.0.0.0        ] [2020-10-12T09:12:23, 811991+0000]
//

//[01505] [tty1] [LOGIN   ] [tty1        ] [                    ] [0.0.0.0        ] [Wed Oct 14 11:10:04 2020 CST]
//[7] [01730] [:0  ] [superred] [tty7        ] [:0                  ] [0.0.0.0        ] [Wed Oct 14 11:10:18 2020 CST]
//[7] [02365] [/0  ] [superred] [pts/0       ] [:0.0                ] [0.0.0.0        ] [Wed Oct 14 11:11:17 2020 CST]
//[7] [02365] [/1  ] [superred] [pts/1       ] [:0.0                ] [0.0.0.0        ] [Wed Oct 14 11:11:58 2020 CST]

//[6] [01515] [tty1] [LOGIN   ] [tty1        ] [                    ] [0.0.0.0        ] [Wed Oct 14 18:08:22 2020 CST]
//[7] [01594] [:0  ] [superred] [tty7        ] [:0                  ] [0.0.0.0        ] [Wed Oct 14 18:08:58 2020 CST]
//[7] [02169] [/0  ] [superred] [pts/0       ] [:0.0                ] [0.0.0.0        ] [Wed Oct 14 18:09:14 2020 CST]

//			"2006-01-02T15:04:05Z07:00"2020-10-12T09:12:23, 811991+0000
//{"address":"0.0.0.0", "device":"~", "exit":{"termination":0, "exit":0}, "host":"5.4.0-48-generic", "id":"~~", "pid":53, "session":0, "time":"Mon, 12 Oct 2020 02:12:38 -0700", "type":"RunLevel", "user":"runlevel"}
//{"address":"0.0.0.0", "device":":0", "exit":{"termination":0, "exit":0}, "host":":0", "id":"", "pid":2277, "session":0, "time":"Mon, 12 Oct 2020 02:13:12 -0700", "type":"UserProcess", "user":"alvin"}
//{"address":"0.0.0.0", "device":"~", "exit":{"termination":0, "exit":0}, "host":"5.4.0-48-generic", "id":"~~", "pid":0, "session":0, "time":"Mon, 12 Oct 2020 02:41:58 -0700", "type":"BootTime", "user":"reboot"}
//{"address":"0.0.0.0", "device":":0", "exit":{"termination":0, "exit":0}, "host":":0", "id":"", "pid":2160, "session":0, "time":"Mon, 12 Oct 2020 02:42:11 -0700", "type":"UserProcess", "user":"alvin"}
//{"address":"0.0.0.0", "device":"~", "exit":{"termination":0, "exit":0}, "host":"5.4.0-48-generic", "id":"~~", "pid":53, "session":0, "time":"Mon, 12 Oct 2020 02:44:58 -0700", "type":"RunLevel", "user":"runlevel"}

func Add8hour(tm string) string {
	h, _ := time.ParseDuration("1h")
	timeTemplate1 := "2006-01-02 15:04:05"
	stamp, _ := time.ParseInLocation(timeTemplate1, tm, time.Local)
	ssssa := stamp.Add(8 * h)
	timeNow := ssssa.String()[0:19]
	return timeNow
}
func JsonToSlice() []WtmpUserLogin {

	wtmpArr := make([]WtmpUserLogin, 0)
	wtmpFiles, _ := ListDir("/var/log", "wtmp")
	fmt.Println("wtmpFiles:", wtmpFiles)
	for _, wtmpFile := range wtmpFiles {
		ret, _ := RunCmd(fmt.Sprintf("./suRun.sh utmpdump %s", wtmpFile))
		retSlice := GetAllMustCompileValue(`\s*\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]`, ret)
		for _, value := range retSlice {
			if len(value) < 9 {
				continue
			}
			timeNew := ""
			if strings.Contains(value[8], "CST") {
				timeNew = TimeToForMat("Mon Jan 02 15:04:05 2006 CST", value[8])

			} else {
				timeNew = TimeToForMat("2006-01-02T15:04:05", strings.Split(value[8], ",")[0])
			}
			//  add  8hour
			//if  Exists("/tmp/.utmpdump"){
			timeNew = Add8hour(timeNew)
			fmt.Println("timeNew**************************************:", timeNew)
			//}
			//fmt.Println(TimeToForMat("2006-01-02T15:04:05","2020-10-14T05:45:23"))
			//fmt.Println(TimeToForMat("Mon Jan 02 15:04:05 2006 CST","Wed Oct 14 11:10:18 2020 CST"))
			//tmpTime, _ := time.ParseInLocation("2006-01-02T15:04:05", strings.Split(value[8], ",")[0], time.Local)
			typeValue, _ := strconv.Atoi(value[1])
			wtmpArr = append(wtmpArr, WtmpUserLogin{
				Time:     timeNew,
				User:     value[4],
				Type:     MarshalJSON(typeValue),
				Exit:     1,
				Terminal: value[3],
				Host:     value[6],
			})
		}

	}

	return wtmpArr
}
