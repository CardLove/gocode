package comtools

import (
	"checkManager/logger"
	"sync"
	"time"
)

type State int

const (
	Start State = iota
	Stop
	Pause
	Continue
)

type ThreadPool struct {
	tasks      chan func() error
	numOfWorks int
	total      int
	results    chan error

	state    State
	IsFinish bool
	StopFlag bool
}

// 初始化
func (t *ThreadPool) Init(number int, total int) {
	t.tasks = make(chan func() error, total)
	t.results = make(chan error, total)
	t.numOfWorks = number
	t.total = total
	t.state = Start
}

func (t *ThreadPool) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	//for {
	// 	 task, ok := <-t.tasks
	//	 if !ok {
	//		 break
	//	 }

	for task := range t.tasks {
		switch t.state {
		case Pause:
			time.Sleep(time.Second)
			continue
		case Stop:
			break
		case Continue, Start:
			result := task()
			t.results <- result
		default:
			logger.GetIns().Debug("error state:%d", t.state)
		}
	}
}

// 开门接客
func (t *ThreadPool) Start() {
	// 开启Number个goroutine
	var wg sync.WaitGroup
	for i := 0; i < t.numOfWorks; i++ {
		wg.Add(1)
		go t.worker(&wg)
	}
	wg.Wait()
	close(t.results)
}

// 关门送客
func (t *ThreadPool) Stop() {
	close(t.tasks)
	close(t.results)
}

// 添加任务
func (t *ThreadPool) AddTask(task func() error) {
	t.tasks <- task
}

func (t *ThreadPool) ChangeState(state State) {
	t.state = state
}
