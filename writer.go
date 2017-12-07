package structured

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
	"time"
)

// StructedWriter wraps a Writer's data in JSON. This is used for passing a
// writer to the standard library logger and outputting JSON encoded logs.
type StructuredWriter struct {
	w     io.Writer
	funcs map[string]func([]byte) (interface{}, error)
}

// New returns a new StructredWriter.
func New(w io.Writer, opts ...StructuredWriterOption) *StructuredWriter {
	wr := &StructuredWriter{
		w:     w,
		funcs: make(map[string]func([]byte) (interface{}, error)),
	}

	for _, o := range opts {
		o.configure(wr)
	}

	return wr
}

// StructuredWriterOption configures a StrucutredWriter.
type StructuredWriterOption interface {
	configure(w *StructuredWriter)
}

// optionFunc enables normal function to implement StructuredWriterOption.
type optionFunc func(w *StructuredWriter)

// configure implements StructuredWriterOption.
func (f optionFunc) configure(w *StructuredWriter) {
	f(w)
}

// WithFieldFunc returns a StructuredWriterOption that adds a field of the
// given name with a value of the result of the funciton. Each write will
// invoke the given function with the corresponding data. If the function
// returns an error, the writer will return the error.
func WithFieldFunc(name string, f func([]byte) (interface{}, error)) StructuredWriterOption {
	return optionFunc(func(w *StructuredWriter) {
		w.funcs[name] = f
	})
}

// WithTimestamp returns a StructuredWriterOption that adds a "timestamp"
// field with the current timestamp in nano-seconds via time.Now().UnixNano().
func WithTimestamp() StructuredWriterOption {
	return optionFunc(func(w *StructuredWriter) {
		w.funcs["timestamp"] = func([]byte) (interface{}, error) {
			return time.Now().UnixNano(), nil
		}
	})
}

// WithCallSite returns a StructuredWriterOption that adds a "callsite" field
// that will give the filename and linenumber of the log.
func WithCallSite() StructuredWriterOption {
	return optionFunc(func(w *StructuredWriter) {
		w.funcs["callsite"] = func([]byte) (interface{}, error) {
			_, fileName, lineNumber, _ := runtime.Caller(1)
			return fmt.Sprintf("%s:%d", path.Base(fileName), lineNumber), nil
		}
	})
}

// Write implements io.Writer. It takes the given data and wraps it into JSON
// marshalled data. The given data will be written to the "msg" key. The
// output is written to a single line. It trims the message of whitespace.
func (w *StructuredWriter) Write(data []byte) (int, error) {
	m := map[string]interface{}{
		"msg": w.parseData(data),
	}

	for k, f := range w.funcs {
		value, err := f(data)
		if err != nil {
			return 0, err
		}
		m[k] = value
	}

	data, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}

	return w.w.Write(data)
}

// parseData will attempt to unmarshal the data in JSON. If it does not, it
// will return a string.
func (w *StructuredWriter) parseData(data []byte) interface{} {
	var u map[string]interface{}
	if err := json.Unmarshal(data, &u); err != nil {
		// User did not provide (valid) JSON
		return strings.TrimSpace(string(data))
	}

	return u
}
