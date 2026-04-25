package server

import (
	"context"
	"errors"
	"hash/fnv"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core"
)

var serverStatus int32 = constant.ServerStatusIdle

type Server struct {
	worker     []*core.Worker
	ioHandlers []*IOHandler
	listeners  []net.Listener

	numWorker    int // for dispatching tasks to workers
	numIOHandler int // for round-robin assignment of new connection to IO Handler
	// For round-robin assignment of new connection to IO Handler
	nextIOHandler int

	mu       sync.Mutex
	wg       sync.WaitGroup
	once     sync.Once
	draining atomic.Bool

	listenerMu sync.Mutex
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

func (s *Server) isDraining() bool {
	return s.draining.Load()
}

func (s *Server) nextHandler() *IOHandler {
	s.mu.Lock()
	defer s.mu.Unlock()

	handler := s.ioHandlers[s.nextIOHandler]
	s.nextIOHandler = (s.nextIOHandler + 1) % len(s.ioHandlers)
	return handler
}

func (s *Server) addListener(listener net.Listener) {
	s.listenerMu.Lock()
	defer s.listenerMu.Unlock()

	s.listeners = append(s.listeners, listener)
}

func (s *Server) closeListeners() {
	s.listenerMu.Lock()
	defer s.listenerMu.Unlock()

	for _, listener := range s.listeners {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("Failed to close listener: %v", err)
		}
	}
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
	if s.isDraining() {
		close(task.ReplyChan)
		return
	}

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

func (s *Server) Shutdown(ctx context.Context) error {
	s.once.Do(func() {
		log.Println("Shutting down server")
		s.draining.Store(true)
		s.closeListeners()

		for _, handler := range s.ioHandlers {
			handler.CloseMultiplexer()
		}

		done := make(chan struct{})
		go func() {
			s.wg.Wait()
			for _, worker := range s.worker {
				worker.Stop()
			}
			for _, handler := range s.ioHandlers {
				handler.CloseConnections()
			}
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
			log.Printf("Graceful shutdown timed out: %v", ctx.Err())
			for _, handler := range s.ioHandlers {
				handler.CloseConnections()
			}
		}
	})

	return ctx.Err()
}

func (s *Server) WaitForSignal(signals chan os.Signal) {
	// Wait for signal in channel, it not available then wait
	<-signals

	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Shutdown finished with error: %v", err)
	}
}
