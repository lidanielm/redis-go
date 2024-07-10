package main

import (
	"bufio"
	"io"
	"fmt"
	"strconv"
	"reflect"
	"errors"
)

type DeserializationError struct

const (
	SIMPLE_STRING	= '+'
	ERROR        	= '-'
	INTEGER      	= ':'
	NULL 		 	= '_'
	BULK_STRING  	= '$'
	ARRAY        	= '*'
	BOOLEAN	     	= '#'
	MAP				= '%'
)

type Value struct {
	typ string
	str string
	num int
	bulk string
	array []Value
	isMap bool
	truth bool
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
	case INTEGER:
		x, err := r.ReadInteger()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "integer", num: x}, nil
	case ERROR:
		return r.ReadError()
	case BULK_STRING:
		return r.ReadBulkString()
	case ARRAY:
		return r.ReadArray()
	case NULL:
		return Value{typ: "null"}, nil
	case BOOLEAN:
		return r.ReadBoolean()
	case MAP:
		// Will return an array with the isMap flag set to true
		return r.ReadMap()
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
		// Read each line with Read because each line can be any type
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}

		arr = append(arr, val)
	}

	return Value{typ: "array", array: arr}, nil
}

func (r *Response) ReadError() (Value, error) {
	line, err := r.ReadLine()
	if err != nil {
		return Value{}, err
	}
	return Value{typ: "error", str: string(line)}, nil
}

func (r *Response) ReadBoolean() (Value, error) {
	line, err := r.ReadLine()
	if err != nil {
		return Value{}, err
	}
	if string(line) == "t" {
		return Value{typ: "boolean", truth: true}, nil
	} else {
		return Value{typ: "boolean", truth: false}, nil
	}
}

func (r *Response) ReadMap() (Value, error) {
	length, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	arr := make([]Value, length) // Initialize array
	for i := 0; i < length; i++ {
		// Read each line with Read because each line can be any type
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}

		arr = append(arr, val)
	}

	return Value{typ: "map", array: arr}, nil

	// // Read first key-value pair to get types
	// key, err := r.Read()
	// if err != nil {
	// 	return Value{}, err
	// }
	// val, err := r.Read()
	// if err != nil {
	// 	return Value{}, err
	// }

	// keyType := key.typ
	// valType := val.typ

	// // Make map
	// m := make(map[string]string)

	// if keyType == "simple_string" || keyType == "bulk_string" || keyType == "error" {
	// 	if valType == "simple_string" || valType == "bulk_string" || valType == "error" {
	// 		m := make(map[string]string)
	// 	} else if valType == "integer" {
	// 		m := make(map[string]int)
	// 	} else if valType == "boolean" {
	// 		m := make(map[string]int)
	// 	} else {
	// 		return Value{}, errors.New("Map keys and values can only be primitive data types.")
	// 	}
	// } else if keyType == "integer" {
	// 	if valType == "simple_string" || valType == "bulk_string" || valType == "error" {
	// 		m := make(map[int]string)
	// 	} else if valType == "integer" {
	// 		m := make(map[int]int)
	// 	} else if valType == "boolean" {
	// 		m := make(map[int]int)
	// 	} else {
	// 		return Value{}, errors.New("Map keys and values can only be primitive data types.")
	// 	}	
	// } else if keyType == "boolean" {
	// 	if valType == "simple_string" || valType == "bulk_string" || valType == "error" {
	// 		m := make(map[bool]string)
	// 	} else if valType == "integer" {
	// 		m := make(map[bool]int)
	// 	} else if valType == "boolean" {
	// 		m := make(map[bool]int)
	// 	} else {
	// 		return Value{}, errors.New("Map keys and values can only be primitive data types.")
	// 	}		
	// } else {
	// 	return Value{}, errors.New("Map keys and values can only be primitive data types.")
	// }

	// // Read array
	// arr := r.ReadArray()
	// for i := 0; i < len(arr); i += 2 {
	// 	k := arr[i]
	// 	v := arr[i+1]
	// 	// Assume that all keys are distinct
	// 	mp[k] = v
	// }

	// // TODO: change this to map
	// return Value{typ: "map", array: arr}, nil
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "simple_string":
		return v.marshalSimpleString()
	case "integer":
		return v.marshalInteger()
	case "error":
		// address error
	case "bulk_string":
		return v.marshalBulkString()
	case "null":
		return v.marshalNull()
	case "array":
		return v.marshalArray()
	case "boolean":
		return v.marshalBoolean()
	case "map":
		return v.marshalMap()
	default:
		return []byte
	}
}

// Helper marshal functions
func (v Value) marshalSimpleString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, v.num)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulkString() []byte {
	len := len(v.str)

	var bytes []byte
	bytes = append(bytes, BULK_STRING)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.array)

	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')
	
	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	
	return bytes
}

func (v Value) marshalBoolean() []byte {
	var bytes []byte
	bytes = append(bytes, BOOLEAN)
	if v.boolean {
		bytes = append(bytes, 't')
	} else {
		bytes = append(bytes, 'f')
	}
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	bytes := []byte{'_', '\r', '\n'}
	return bytes
}

func (v Value) marshalMap() []byte {
	numPairs := len(v.array) / 2

	var bytes []byte
	bytes = append(bytes, MAP)
	bytes = append(bytes, strconv.Itoa(numPairs)...)
	bytes = append(bytes, '\r', '\n')
	
	for i := 0; i < 2 * numPairs; i++ {
		bytes := append(bytes, v.array[i].Marshal()...)
	}

	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}