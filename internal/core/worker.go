package core

import (
	"log"
	"sync"
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

type Task struct {
	Command   *Command
	ReplyChan chan []byte
}

type Worker struct {
	id       int
	redisDB  *RedisDB
	TaskChan chan *Task
	once     sync.Once
	wg       sync.WaitGroup
}

func NewWorker(id int, bufferSize int) *Worker {
	worker := &Worker{
		id:       id,
		redisDB:  NewRedisDB(),
		TaskChan: make(chan *Task, bufferSize),
	}
	worker.wg.Add(1)
	go worker.run()
	return worker
}

func (w *Worker) run() {
	defer w.wg.Done()
	log.Printf("Worker %d started\n", w.id)
	// Not like single-threaded, active expire is triggered before executing the command, so we need to check expire before executing the command
	// We can also use a ticker to trigger active expire periodically
	ticker := time.NewTicker(constant.ActiveExpireFrequency)
	defer ticker.Stop()

	for {
		select {
		case task, ok := <-w.TaskChan:
			if !ok {
				log.Printf("Worker %d stopped\n", w.id)
				return
			}
			log.Printf("Worker %d handling the task", w.id)
			w.ExecuteAndRespond(task)
		case <-ticker.C:
			ActiveDeleteExpiredKeys(w.redisDB)
		}
	}
}

func (w *Worker) Stop() {
	w.once.Do(func() {
		close(w.TaskChan)
	})
	w.wg.Wait()
}

func (w *Worker) ExecuteAndRespond(task *Task) {
	res := ExecuteCommand(w.redisDB, task.Command)

	task.ReplyChan <- res
}
