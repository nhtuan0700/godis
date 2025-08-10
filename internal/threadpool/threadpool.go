package threadpool

import (
	"io"
	"log"
	"net"
	"time"
)

// elemenet in the queue
type Job struct {
	conn net.Conn
}

// represent thread in the pool
type Worker struct {
	id      int
	jobChan chan Job
}

func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			handleConnection(job.conn)
		}
	}()
}

// thread pool
type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

func NewPool(numOfWorkers int) *Pool {
	return &Pool{
		workers:  make([]*Worker, numOfWorkers),
		jobQueue: make(chan Job),
	}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		p.workers[i].Start()
	}
}

// push connection to the queue
func (p *Pool) AddJob(conn net.Conn) {
	log.Println("Adding job from: ", conn.RemoteAddr())
	p.jobQueue <- Job{conn: conn}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Read data from client
	buf := make([]byte, 1000)
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			netErr, ok := err.(net.Error)
			switch {
			case ok && netErr.Timeout():
				log.Println("Read timeout")
			case err == io.EOF:
				log.Printf("client %s closed connection", conn.RemoteAddr())
			default:
				log.Printf("read error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		// process
		time.Sleep(time.Second * 10)
		log.Printf("Request from %s\n", conn.RemoteAddr())
		conn.Write([]byte("HTTP/1.1 200 OK \r\n\r\nWelcome to Godis!\r\n"))
	}
}
