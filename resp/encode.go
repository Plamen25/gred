package resp

import (
	"errors"
	"io"
	"strconv"
)

// ErrInvalidValue is returned if the value to encode is invalid.
var ErrInvalidValue = errors.New("resp: invalid value")

// Encode encode the value v and writes the serialized data to w.
func Encode(w io.Writer, v interface{}) error {
	return encodeValue(w, v)
}

// encodeValue encodes the value v and writes the serialized data to w.
func encodeValue(w io.Writer, v interface{}) error {
	switch v := v.(type) {
	case SimpleString:
		return encodeSimpleString(w, v)
	case Error:
		return encodeError(w, v)
	case Integer:
		return encodeInteger(w, v)
	case BulkString:
		return encodeBulkString(w, v)
	case Array:
		return encodeArray(w, v)
	default:
		return ErrInvalidValue
	}
}

// encodeArray encodes an array value to w.
func encodeArray(w io.Writer, v Array) error {
	// Special case for a nil array
	if v == nil {
		err := encodePrefixed(w, '*', []byte("-1"))
		return err
	}

	// First encode the number of elements
	n := len(v)
	err := encodePrefixed(w, '*', []byte(strconv.Itoa(n)))
	if err != nil {
		return err
	}

	// Then encode each value
	for _, el := range v {
		err = encodeValue(w, el)
		if err != nil {
			return err
		}
	}
	return nil
}

// encodeBulkString encodes a bulk string to w.
func encodeBulkString(w io.Writer, v BulkString) error {
	// Special case for a nil bulk string
	if v == nil {
		err := encodePrefixed(w, '$', []byte("-1"))
		return err
	}

	// First encode the length
	n := len(v)
	err := encodePrefixed(w, '$', []byte(strconv.Itoa(n)))
	if err != nil {
		return err
	}
	// Then the string
	_, err = w.Write(append(v, '\r', '\n'))
	return err
}

// encodeInteger encodes an integer value to w.
func encodeInteger(w io.Writer, v Integer) error {
	return encodePrefixed(w, ':', []byte(strconv.FormatInt(int64(v), 10)))
}

// encodeSimpleString encodes a simple string value to w.
func encodeSimpleString(w io.Writer, v SimpleString) error {
	return encodePrefixed(w, '+', v)
}

// encodeError encodes an error value to w.
func encodeError(w io.Writer, v Error) error {
	return encodePrefixed(w, '-', v)
}

// encodePrefixed encodes the data v to w, with the specified prefix.
func encodePrefixed(w io.Writer, prefix byte, v []byte) error {
	buf := make([]byte, len(v)+3)
	buf[0] = prefix
	copy(buf[1:], v)
	copy(buf[len(buf)-2:], "\r\n")
	_, err := w.Write(buf)
	return err
}