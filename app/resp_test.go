package main

import (
	"testing"
	"bufio"
	"strings"
)

func TestReadLine(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("test\r\n"))
	resp := NewResponse(reader)
	line, _, err := resp.ReadLine()
	if err != nil {
		t.Errorf("Failed to read line: %s", err.Error())
	}
	if string(line) != "test" {
		t.Errorf("Expected 'test', got '%s'", string(line))
	}
}