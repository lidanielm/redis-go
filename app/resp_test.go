package main

import (
	"testing"
	"bufio"
	"strings"
)

func TestReadLine(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("test\r\n"))
	resp := NewResponse(reader)
	line, err := resp.ReadLine()
	if err != nil {
		t.Errorf("Failed to read line: %s", err.Error())
	}
	if string(line) != "test" {
		t.Errorf("Expected 'test', got '%s'", string(line))
	}
}

// func TestReadMap(t *testing.T) {
// 	reader := bufio.NewReader(strings.NewReader("%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n"))
// 	resp := NewResponse(reader)
// 	val, err := resp.Read()
// 	if err != nil {
// 		t.Errorf("Failed to read line: %s", err.Error())
// 	}
// 	fmt.Println(val)
// }