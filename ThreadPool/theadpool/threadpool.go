package threadpool

import (
	"io"
	"log"
	"net"
	"sync"
)

// elemenet in the queue
type Job struct {
	conn net.Conn
}

// represent thread in the pool
type Worker struct {
	id      int
	jobChan chan Job
	// stopChan from pool
	stopChan chan struct{}
	// wg from pool
	wg *sync.WaitGroup
}

func NewWorker(id int, jobChan chan Job, stopChan chan struct{}, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:       id,
		jobChan:  jobChan,
		stopChan: stopChan,
		wg:       wg,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer log.Printf("Worker %d closed", w.id)

		for {
			select {
			case job, ok := <-w.jobChan:
				if !ok {
					return
				}
				handleConnection(job.conn)
			case <-w.stopChan:
				return
			}
		}
	}()
}

// thread pool
type Pool struct {
	jobQueue chan Job
	workers  []*Worker
	stopChan chan struct{}   // channel to gracefully stop the pool
	wg       *sync.WaitGroup // wg to wait for all workers to finish
}

func NewPool(numOfWorkers int) *Pool {
	return &Pool{
		workers:  make([]*Worker, numOfWorkers),
		jobQueue: make(chan Job),
		stopChan: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue, p.stopChan, p.wg)
		p.workers[i] = worker
		p.workers[i].Start()
	}
}

// stop gracefully shutdown pool
func (p *Pool) Close() {
	log.Println("Closing pool")
	close(p.stopChan)

	// wait for all workers to finish
	log.Println("Waiting for all workers to finish")
	p.wg.Wait()

	// close job queue
	close(p.jobQueue)
	log.Println("Pool closed successfully")
}

// check if the pool is closed
func (p *Pool) IsClosed() bool {
	select {
	case <-p.stopChan:
		return true
	default:
		return false
	}
}

// push connection to the queue
func (p *Pool) AddJob(conn net.Conn) {
	if p.IsClosed() {
		conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Type: text/plain\r\n\r\nPool is closed\r\n"))
		conn.Close()
		return
	}
	select {
	case p.jobQueue <- Job{conn: conn}:
		log.Println("Adding job from: ", conn.RemoteAddr())
		return
	case <-p.stopChan:
		conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Type: text/plain\r\n\r\nPool is closed\r\n"))
		conn.Close()
		return
	}
}

func handleConnection(conn net.Conn) {
	log.Println("Handle conn from ", conn.RemoteAddr())
	defer conn.Close()
	// Read data from client
	for {
		// conn.SetReadDeadline(time.Now().Add(time.Minute))
		cmd, err := readCommand(conn)
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
		// time.Sleep(time.Second * 10)
		if err := respond(cmd, conn); err != nil {
			log.Println("err write: ", err)
		}
	}
}

func readCommand(conn net.Conn) (string, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}

	return nil
}
