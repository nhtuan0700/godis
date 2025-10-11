package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Process ID: ", os.Getgid())
	
	// 1. Create a channel to receive system signals.
	// This channel will be the asynchronous endpoint for signal notifications.
	// think it as a queue
	sigs := make(chan os.Signal, 1)

	// 2. Register the signals we want to catch
	// We want to catch SIGINT (interupt from the terminal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 3. Create a channel to block the main goroutine until a shutdown is complete
	done := make(chan bool)

	go func() {
		// block
		// dequeue, pop from the queue
		sig := <-sigs

		fmt.Printf("\n\n[HANDLER] Received signal: %v\n", sig)
		// Signal the main goroutine that the handler has finished and it's safe to exit
		done <- true
	}()

	fmt.Println("[MAIN] Waiting for work or signal...")
	<-done

	fmt.Println("[MAIN] Application shutdown successfully")
}
