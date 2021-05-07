package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/quoeamaster/golang_blogs/ggsingleton/logger"
)

func TestMultiRoutineAccess_loggerImpl(t *testing.T) {
	// storing address of *Logger instance for validation.
	msgArray := make([]string, 0)

	// wait-group to make sure a certain number of go-routine(s)
	// has finished its task.
	var wgroup sync.WaitGroup

	for i := 0; i < 10; i++ {
		// updates the wait-group counter.
		wgroup.Add(1)

		go func(idx int) {
			// decreses the wait-group counter by 1.
			// When the counter returns to 0, the wait-group will end the "wait".
			defer wgroup.Done()

			log := logger.GetLoggerImpl()
			// append the address value of instance "log"
			msgArray = append(msgArray, fmt.Sprintf("%p", log))
			log.Log(fmt.Sprintf("{loggerImpl} this is a log entry from [%v]\n", idx))
		}(i)
	}
	wgroup.Wait()

	// verification
	if len(msgArray) == 0 {
		t.Fatalf("expect to have a least one message")
	}
	addrLine := msgArray[0]
	for i := 1; i < len(msgArray); i++ {
		line := msgArray[i]
		if addrLine != line {
			t.Errorf("expect both lines (addresses of Logger) should be identical, [%v] vs [%v]\n", addrLine, line)
		}
	}
}

func TestMultiRoutineAccessWithDelay_loggerImpl(t *testing.T) {
	// storing address of *Logger instance for validation.
	msgArray := make([]string, 0)

	// wait-group to make sure a certain number of go-routine(s)
	// has finished its task.
	var wgroup sync.WaitGroup

	for i := 0; i < 10; i++ {
		// updates the wait-group counter.
		wgroup.Add(1)

		go func(idx int) {
			// decreses the wait-group counter by 1.
			// When the counter returns to 0, the wait-group will end the "wait".
			defer wgroup.Done()

			// add a random delay to simulate multi access.
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)))

			log := logger.GetLoggerImpl()
			// append the address value of instance "log"
			msgArray = append(msgArray, fmt.Sprintf("%p", log))
			log.Log(fmt.Sprintf("[with delay] {loggerImpl} this is a log entry from [%v]\n", idx))
		}(i)
	}
	wgroup.Wait()

	// verification
	if len(msgArray) == 0 {
		t.Fatalf("expect to have a least one message")
	}
	addrLine := msgArray[0]
	for i := 1; i < len(msgArray); i++ {
		line := msgArray[i]
		if addrLine != line {
			t.Errorf("expect both lines (addresses of Logger) should be identical, [%v] vs [%v]\n", addrLine, line)
		}
	}
}
