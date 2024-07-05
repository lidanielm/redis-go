package main

import (
	"bufio"
	"io"
	"fmt"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	ERROR         = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

type Value struct {
	typ string
	str string
	num int
	bulk string
	array []Value
}

type Response struct {
	reader *bufio.Reader
}

func NewResponse(ioReader io.Reader) *Response {
	reader := bufio.NewReader(ioReader)
	return &Response{reader: reader}
}

func (r *Response) ReadLine() ([]byte, error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil || len(line) == 0 {
		return nil, err
	}
	if len(line) > 0 && line[len(line)-1] == '\n' && line[len(line)-2] == '\r' {
		line = line[:len(line)-2] // remove \r\n
		return line, nil
	}
	return line, err
}

func (r *Response) ReadByte() (byte, error) {
	return r.reader.ReadByte()
}

func (r *Response) ReadInteger() (int, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

func (r *Response) Read() (Value, error) {
	respType, err := r.ReadByte()
	if err != nil {
		return Value{}, err
	}
	
	switch respType {
	case SIMPLE_STRING:
		return r.ReadSimpleString()
	case BULK_STRING:
		return r.ReadBulkString()
	case ARRAY:
		return r.ReadArray()
	default:
		fmt.Println("Unknown type: ", respType)
		return Value{}, nil
	}
}

func (r *Response) ReadSimpleString() (Value, error) {
	line, err := r.ReadLine()
	if err != nil {
		return Value{}, err
	}
	return Value{typ: "simple_string", str: string(line)}, nil
}

func (r *Response) ReadBulkString() (Value, error) {
	size, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	} else if size == -1 {
		return Value{}, nil // Null bulk string
	}
	data, err := r.ReadLine()
	if err != nil {
		return Value{}, err
	}
	return Value{typ: "bulk_string", bulk: string(data)}, nil
}

func (r *Response) ReadArray() (Value, error) {
	length, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	arr := make([]Value, length) // Initialize array
	for i := 0; i < length; i++ {
		// Read each line with Read because each line is a bulk string
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}

		arr = append(arr, val)
	}

	return Value{typ: "array", array: arr}, nil
}

// func ReadCRLF(data []byte, atEOF bool) (token []byte, err error) {
// 	if atEOF && len(data) == 0 {
// 		return nil, nil
// 	}
// 	if i := bytes.Index(data, []byte{'\r','\n'}); i >= 0 {
// 		// CRLF found
// 		return DropCR(data[:i]), nil
// 	}

// 	// atEOF with no trailing \r\n
// 	if atEOF {
// 		return DropCR(data), nil
// 	}
// 	return nil, nil
// }

// func DropCR(data []byte) []byte {
// 	// Drop the \r from the end of the data
// 	if len(data) > 0 && data[len(data) - 1] == '\r' {
// 		return data[:len(data) - 1]
// 	}
// }