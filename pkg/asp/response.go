package asp

import (
	"strconv"
	"strings"
	"time"
)

const (
	responseTypeOK    = "O"
	responseTypeError = "E"

	lineTypeHeader = "H"
	lineTypeData   = "D"

	delimiter = "\t"
	linebreak = "\n"
)

type Response struct {
	lines [][]string
}

func NewOKResponse() *Response {
	r := &Response{}
	return r.write(responseTypeOK)
}

func NewErrorResponse(code int) *Response {
	r := &Response{}
	return r.write(responseTypeError, strconv.Itoa(code))
}

func (r *Response) WriteHeader(elems ...string) *Response {
	return r.write(lineTypeHeader, elems...)
}

func (r *Response) WriteData(elems ...string) *Response {
	return r.write(lineTypeData, elems...)
}

func (r *Response) write(lineType string, elems ...string) *Response {
	r.lines = append(r.lines, append([]string{lineType}, elems...))
	return r
}

func (r *Response) AppendHeader(elems ...string) *Response {
	return r.append(lineTypeHeader, elems...)
}

func (r *Response) AppendData(elems ...string) *Response {
	return r.append(lineTypeData, elems...)
}

// append Appends elems to the last line of the given lineType or adds a new line if no line of lineType is found
func (r *Response) append(lineType string, elems ...string) *Response {
	// Iterate lines in reverse order, append element to first line with given type
	for i := len(r.lines) - 1; i >= 0; i-- {
		if len(r.lines[i]) > 0 && r.lines[i][0] == lineType {
			r.lines[i] = append(r.lines[i], elems...)
			return r
		}
	}

	// No line with given type was found, write a new one
	return r.write(lineType, elems...)
}

func (r *Response) Serialize() string {
	serialized := ""
	size := 0
	for _, line := range r.lines {
		for j, elem := range line {
			serialized += elem
			if j+1 < len(line) {
				serialized += delimiter
			}
			// Cannot use len(elem) here, since it counts bytes not characters
			size += len([]rune(elem))
		}

		// Still need to append the line indicating the size, so add linebreak for every line
		serialized += linebreak
	}

	serialized += strings.Join([]string{"$", strconv.Itoa(size), "$"}, delimiter)

	return serialized
}

func NewErrorResponseWithMessage(code int, message string) *Response {
	return NewErrorResponse(code).
		WriteHeader("asof", "err").
		WriteData(Timestamp(), message)
}

func NewSyntaxErrorResponse() *Response {
	return NewErrorResponseWithMessage(107, "Invalid Syntax!")
}

func Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
