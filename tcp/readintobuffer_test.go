package main

import (
	"crypto/rand"
	"io"
	"math"
	"net"
	"testing"
	"time"
)

func TestReadIntoBuffer(t *testing.T) {
	payload := make([]byte, 1<<24)
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1<<19)
	total := 0
	count := 0
	start := time.Now()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		total += n
		t.Logf("%d: read %.1f kb", count, math.Ceil(float64(n)/1024))
		count += 1
	}

	conn.Close()
	end := time.Now()
	elapsedMs := end.Sub(start) / time.Millisecond
	totalMb := math.Ceil(float64(total) / 1024 / 1024)
	downloadSpeed := totalMb / float64(elapsedMs)

	t.Logf("TOTAL: read %.1f mb", totalMb)
	t.Logf("Time Elapsed: %dms", elapsedMs)
	t.Logf("Download speed AVG: %.2f mb/ms which is %.f mb/s", downloadSpeed, downloadSpeed*1000)
}
