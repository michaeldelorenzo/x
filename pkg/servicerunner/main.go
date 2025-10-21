package servicerunner

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Semaphore is a custom no-memory-usage channel type for signaling
type Semaphore chan struct{}

var (
	// QuitSignal is a used to signal the quit chanel to unblock
	QuitSignal = struct{}{}
)

// Runner is responsible for blocking functionality
type Runner struct {
	quitChannel Semaphore
	cleanup     func()
	blocking    bool
}

// NewServiceRunner will execute the func it receives in a non-blocking manner
// and returns a Runner type
func NewServiceRunner(routine func()) *Runner {
	if routine != nil {
		go routine()
	}

	b := &Runner{
		quitChannel: NewSemaphore(),
	}

	return b
}

// Start will block on the calling goroutine until Stop() is called.
// An error will be printed if block has already been run without unblocking
func (b *Runner) Start() {
	if b.isCurrentlyBlocking() {
		err := errors.New("error: unable to block because it's already blocking")
		log.Println(err)
		return
	}
	b.blocking = true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)

	go func() {
		<-c
		if b.cleanup != nil {
			b.cleanup()
		}
		b.quitChannel <- QuitSignal
		return
	}()

	<-b.quitChannel
}

// Stop will continue execution on the calling goroutine.
// If Start() was never called previously, error will be printed
func (b *Runner) Stop() {
	if b.quitChannel == nil || !b.blocking {
		err := errors.New("error: unable to unblock because it's not blocking")
		log.Println(err)
		return
	}

	defer func() {
		b.quitChannel = nil
		b.blocking = false
	}()

	if b.cleanup != nil {
		b.cleanup()
	}

	b.quitChannel <- QuitSignal
}

// NewSemaphore returns an initialized chan struct{}{}
func NewSemaphore() Semaphore {
	return make(Semaphore)
}

// SetCleanup injects the closure to be run if OS term & kill
// signals are detected.
func (b *Runner) SetCleanup(routine func()) *Runner {
	b.cleanup = routine

	return b
}

func (b *Runner) isCurrentlyBlocking() bool {
	var isBlocking bool
	if b.blocking {
		return true
	}

	return isBlocking
}
