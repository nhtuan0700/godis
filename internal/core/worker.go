package core

import (
	"errors"
	"log"
	"strconv"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

type Task struct {
	Command   *Command
	ReplyChan chan []byte
}

type Worker struct {
	id        int
	dictStore *data_structure.Dict
	TaskChan  chan *Task
}

func NewWorker(id int, bufferSize int) *Worker {
	worker := &Worker{
		id:        id,
		dictStore: data_structure.NewDict(),
		TaskChan:  make(chan *Task, bufferSize),
	}
	go worker.run()
	return worker
}

func (w *Worker) run() {
	log.Printf("Worker %d started\n", w.id)
	for task := range w.TaskChan {
		log.Printf("Worker %d handling the task", w.id)
		w.ExecuteAndRespond(task)
	}
}

func (w *Worker) ExecuteAndRespond(task *Task) {
	var res []byte
	switch task.Command.Cmd {
	case constant.CMD_PING:
		res = w.cmdPING(task.Command.Args)
	case constant.CMD_SET:
		res = w.cmdSET(task.Command.Args)
	case constant.CMD_GET:
		res = w.cmdGET(task.Command.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	task.ReplyChan <- res
}

func (w *Worker) cmdSET(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	var ttlMs uint64 = 0

	key, value = args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = uint64(ttlSec) * 1000
	}

	w.dictStore.Set(key, w.dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func (w *Worker) cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func (w *Worker) cmdGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}
	key := args[0]
	obj := w.dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if w.dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}
