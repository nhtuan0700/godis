package server

import (
	"hash/fnv"
	"log"
	"math/rand"
	"runtime"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core"
)

var serverStatus int32 = constant.ServerStatusIdle

type Server struct {
	worker     []*core.Worker
	ioHandlers []*IOHandler

	numWorker    int // for dispatching tasks to workers
	numIOHandler int // for round-robin assignment of new connection to IO Handler
	// For round-robin assignment of new connection to IO Handler
	nextIOHandler int
}

func NewServer() (*Server, error) {
	numCore := runtime.NumCPU()
	numIOHandler := numCore / 2
	numWorker := numCore / 2

	log.Printf("Initialize server with %d IO Handlers and %d Workers \n", numIOHandler, numWorker)
	server := &Server{
		worker:       make([]*core.Worker, numWorker),
		ioHandlers:   make([]*IOHandler, numIOHandler),
		numWorker:    numWorker,
		numIOHandler: numIOHandler,
	}

	for i := 0; i < numWorker; i++ {
		server.worker[i] = core.NewWorker(i, 1024)
	}

	for i := 0; i < numIOHandler; i++ {
		ioHandler, err := NewIOHandler(i, server)
		if err != nil {
			return nil, err
		}

		server.ioHandlers[i] = ioHandler
	}

	return server, nil
}

func (s *Server) getWorkerID(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % s.numWorker
}

// set k1 123
// k1 -> 1
// get k1
// k1 -> 1
func (s *Server) dispatch(task *core.Task) {
	// For commands like PING etc., dont have a key
	// We can send them to any worker
	// TODO for commands with keys, we will
	var workerID int
	if len(task.Command.Args) > 0 {
		key := task.Command.Args[0]
		workerID = s.getWorkerID(key)
	} else {
		workerID = rand.Intn(s.numWorker)
	}

	s.worker[workerID].TaskChan <- task
}
