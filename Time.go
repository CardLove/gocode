package comtools

import "time"

func TimeToForMat(timeFormat string ,timeSrc  string ) string  {
	timeLayout := "2006-01-02 15:04:05"
	timeTemp, _ := time.ParseInLocation(timeFormat, timeSrc, time.Local)
	return  timeTemp.Format(timeLayout)
}