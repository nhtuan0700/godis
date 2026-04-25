package core

import (
	"log"
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

type Task struct {
	Command   *Command
	ReplyChan chan []byte
}

type Worker struct {
	id        int
	redisDB   *RedisDB
	TaskChan  chan *Task
}

func NewWorker(id int, bufferSize int) *Worker {
	worker := &Worker{
		id:        id,
		redisDB:   NewRedisDB(),
		TaskChan:  make(chan *Task, bufferSize),
	}
	go worker.run()
	return worker
}

func (w *Worker) run() {
	log.Printf("Worker %d started\n", w.id)
	// Not like single-threaded, active expire is triggered before executing the command, so we need to check expire before executing the command
	// We can also use a ticker to trigger active expire periodically
	ticker := time.NewTicker(constant.ActiveExpireFrequency)
	
	for {
		select {
		case task := <-w.TaskChan:
			log.Printf("Worker %d handling the task", w.id)
			w.ExecuteAndRespond(task)
		case <-ticker.C:
			ActiveDeleteExpiredKeys(w.redisDB)
		}
	}
}

func (w *Worker) ExecuteAndRespond(task *Task) {
	res := ExecuteCommand(w.redisDB, task.Command)

	task.ReplyChan <- res
}
