package comtools

import (
	"checkManager/data"
	"checkManager/pubConst"
)

func GetState(statue data.ETaskStatue ) string {
	//生成当前程序运行的状态
	taskStatue := "0"
	switch statue {
	case data.EStateStop:
		taskStatue = pubConst.STOP
	case data.EStateRunning:
		taskStatue = pubConst.RUNNING
	case data.EStatePause:
		taskStatue = pubConst.PAUSE
	default:
		taskStatue = "null"
	}
	return taskStatue
}
