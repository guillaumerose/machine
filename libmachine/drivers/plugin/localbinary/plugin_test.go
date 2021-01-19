package localbinary

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLocalBinaryPluginAddress(t *testing.T) {
	lbp := &Plugin{}
	expectedAddr := "127.0.0.1:12345"

	lbp.addrCh = make(chan string, 1)
	lbp.addrCh <- expectedAddr

	// Call the first time to read from the channel
	addr, err := lbp.Address()
	if err != nil {
		t.Fatalf("Expected no error, instead got %s", err)
	}
	if addr != expectedAddr {
		t.Fatal("Expected did not match actual address")
	}

	// Call the second time to read the "cached" address value
	addr, err = lbp.Address()
	if err != nil {
		t.Fatalf("Expected no error, instead got %s", err)
	}
	if addr != expectedAddr {
		t.Fatal("Expected did not match actual address")
	}
}

func TestLocalBinaryPluginAddressTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test")
	}

	lbp := &Plugin{
		addrCh:  make(chan string, 1),
		timeout: 1 * time.Second,
	}

	addr, err := lbp.Address()

	assert.Empty(t, addr)
	assert.EqualError(t, err, "Failed to dial the plugin server in 1s")
}

func TestLocalBinaryPluginClose(t *testing.T) {
	lbp := &Plugin{}
	lbp.stopCh = make(chan struct{})
	go lbp.Close()
	_, isOpen := <-lbp.stopCh
	if isOpen {
		t.Fatal("Close did not send a stop message on the proper channel")
	}
}

func TestExecServer(t *testing.T) {
	var buffer bytes.Buffer

	log.SetLevel(log.DebugLevel)
	log.SetOutput(io.MultiWriter(&buffer, os.Stderr))

	executor := &FakeExecutor{}
	lbp := &Plugin{
		MachineName: "test",
		Executor:    executor,
		addrCh:      make(chan string, 1),
		stopCh:      make(chan struct{}),
		timeout:     time.Second,
	}

	done := make(chan error)
	go func() {
		done <- lbp.execServer()
	}()

	addr, err := lbp.Address()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1:12345", addr)

	assert.Eventually(t, func() bool {
		return strings.Contains(buffer.String(), "(test) DBG | Uh oh, something in plugin went wrong...") &&
			strings.Contains(buffer.String(), "(test) Doing some fun plugin stuff...")
	}, time.Second, 100*time.Millisecond) // logs do not appear instantly

	_ = lbp.Close()
	assert.NoError(t, <-done)
}

type FakeExecutor struct {
}

func (fe *FakeExecutor) Start() (*bufio.Scanner, *bufio.Scanner, error) {
	stdout := "127.0.0.1:12345\nDoing some fun plugin stuff...\n"
	stderr := "Uh oh, something in plugin went wrong...\n"
	return bufio.NewScanner(bytes.NewReader([]byte(stdout))), bufio.NewScanner(bytes.NewReader([]byte(stderr))), nil
}

func (fe *FakeExecutor) Close() error {
	return nil
}
