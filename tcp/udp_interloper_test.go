package main

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestListenPacketUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "localhost:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "localhost:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	interloper, err := net.ListenPacket("udp", "localhost:")
	if err != nil {
		t.Fatal(err)
	}
	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal()
	}
	_ = interloper.Close()
	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	_, err = client.WriteTo(ping, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(interrupt, buf[:n]) {
		t.Errorf("expected %q; actual %q", interrupt, buf[:n])
	}
	if addr.String() != interloper.LocalAddr().String() {
		t.Errorf("expected msg from %q; actual sender is %q", interloper.LocalAddr(), addr)
	}

	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected %q; actual %q", ping, buf[:n])
	}
	if addr.String() != serverAddr.String() {
		t.Errorf("expected msg from %q; actual sender is %q", serverAddr, addr)
	}
}
