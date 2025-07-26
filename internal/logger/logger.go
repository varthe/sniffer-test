package logger

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type record struct {
	ts  string
	src string
	url string
	msg string
}

var (
	ch   = make(chan record, 1<<16)
	done = make(chan struct{})
)

func init() {
	file, err := os.OpenFile("/app/sniffer.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	bw := bufio.NewWriterSize(file, 1<<20)
	w := csv.NewWriter(bw)

	if stat, _ := file.Stat(); stat.Size() == 0 {
		_ = w.Write([]string{"Timestamp", "Source", "URL", "Comment"})
	}

	go func() {
		t := time.NewTicker(250 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case rec := <-ch:
				_ = w.Write([]string{rec.ts, rec.src, rec.url, rec.msg})
			case <-t.C:
				w.Flush()
				bw.Flush()
			case <-done:
				close(ch)
				for rec := range ch {
					_ = w.Write([]string{rec.ts, rec.src, rec.url, rec.msg})
				}
				w.Flush()
				bw.Flush()
				file.Close()
				return
			}
		}
	}()

	runtime.SetFinalizer(file, func(*os.File) { close(done) })
}

func File(src, url, msg string) {
	ch <- record{
		ts:  time.Now().Format(time.RFC3339Nano),
		src: src,
		url: url,
		msg: msg,
	}
}

func Console(message string, args ...any) {
	msg := fmt.Sprintf(message, args...)
	log.Println(msg)
}

func Fatal(message string, args ...any) {
	msg := fmt.Sprintf(message, args...)
	log.Println(msg)
	os.Exit(1)
}