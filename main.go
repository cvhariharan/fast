package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/dustin/go-humanize"
)

const (
	TmpFile = ".test"
	FileURL = "https://rawcdn.githack.com/cvhariharan/fast/d0df10713b43baf8c62297d07a4c70841a0c9334/assets/globe.tif"

	// SleepTime in milliseconds
	SleepTime = 1000
)

type ProgressCounter struct {
	Length       int64
	Total        int64
	InitTime     time.Time
	PrevTime     time.Time
	CurrentSpeed float64
	Closed       bool
}

func (c *ProgressCounter) Write(p []byte) (int, error) {
	c.Length += int64(len(p))
	c.Total += int64(len(p))
	return int(c.Length), nil
}

func (c *ProgressCounter) Progress() {
	elapsed := time.Since(c.PrevTime)
	c.PrevTime = time.Now()
	c.CurrentSpeed = float64(c.Length) / elapsed.Seconds()
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rCurrent speed: %s/s", humanize.Bytes(uint64(c.CurrentSpeed)))
	c.Length = 0
}

func (c *ProgressCounter) Close() {
	c.Closed = true
}

func NewProgressCounter() *ProgressCounter {
	return &ProgressCounter{
		PrevTime: time.Now(),
		InitTime: time.Now(),
	}
}

func main() {
	var wg sync.WaitGroup

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	resp, err := http.Get(FileURL)
	if err != nil {
		log.Println(err)
	}

	p := NewProgressCounter()
	tmpFile, err := ioutil.TempFile(os.TempDir(), "fast-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	defer p.Close()
	defer os.Remove(tmpFile.Name())

	wg.Add(1)
	go func(p *ProgressCounter) {
		for !p.Closed {
			p.Progress()
			time.Sleep(SleepTime * time.Millisecond)
		}
		fmt.Printf("\nAverage speed: %s/s\n", humanize.Bytes(uint64(p.Total/int64(time.Since(p.InitTime).Seconds()))))
		fmt.Printf("Started %s\n", humanize.Time(p.InitTime))
		wg.Done()
	}(p)

	if _, err = io.Copy(tmpFile, io.TeeReader(resp.Body, p)); err != nil {
		tmpFile.Close()
		log.Println(err)
	}

	wg.Wait()
}
