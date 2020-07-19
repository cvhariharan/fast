package main

import (
	"fmt"
	"io"
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
	FileURL = "http://212.183.159.230/5MB.zip"
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
	out, err := os.Create(TmpFile)
	if err != nil {
		log.Println(err)
	}

	wg.Add(1)
	go func(p *ProgressCounter) {
		var values []float64
		for !p.Closed {
			p.Progress()
			values = append(values, p.CurrentSpeed)
			time.Sleep(1000 * time.Millisecond)
		}
		var sum float64
		for _, v := range values {
			sum += v
		}
		fmt.Printf("\nAverage speed: %s/s\n", humanize.Bytes(uint64(p.Total/int64(time.Since(p.InitTime).Seconds()))))
		fmt.Printf("Started %s\n", humanize.Time(p.InitTime))
		wg.Done()
	}(p)

	if _, err = io.Copy(out, io.TeeReader(resp.Body, p)); err != nil {
		out.Close()
		log.Println(err)
	}
	os.Remove(TmpFile)
	p.Close()

	wg.Wait()
}
