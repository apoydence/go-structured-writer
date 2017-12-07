package structured_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	structured "github.com/apoydence/go-structured-writer"
)

func TestStructuredWriter(t *testing.T) {
	t.Parallel()

	b := bytes.Buffer{}
	w := structured.New(&b)

	// Should trim the \n
	n, err := w.Write([]byte("some log msg\n"))
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	if n != b.Len() {
		t.Fatalf("expected n (%d) to equal %d", n, b.Len())
	}

	var m map[string]interface{}
	err = json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		t.Fatalf("invalid json: %s", err)
	}

	if m["msg"] != "some log msg" {
		t.Fatalf("expected %v to equal 'some log msg'", m["msg"])
	}

	// Must have newline as last rune
	if b.Bytes()[b.Len()-1] != '\n' {
		t.Fatalf("expected the last rune (%v) to be a '\\n'", b.Bytes()[b.Len()-1])
	}
}

func TestWithFieldFunc(t *testing.T) {
	t.Parallel()

	b := bytes.Buffer{}
	var data []byte
	var returnErr error
	w := structured.New(&b, structured.WithFieldFunc("new-field", func(d []byte) (interface{}, error) {
		data = d
		return 99, returnErr
	}))

	_, err := w.Write([]byte("some log msg"))
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	if string(data) != "some log msg" {
		t.Fatalf("expected %s to equal 'some log msg'", data)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		t.Fatalf("invalid json: %s", err)
	}

	if m["new-field"] != 99.0 {
		t.Fatalf("expected %v to equal 99.0", m["new-field"])
	}

	returnErr = errors.New("some-error")
	_, err = w.Write([]byte("some log msg"))
	if err == nil {
		t.Fatal("expected err to not be nil")
	}
}

func TestWithTimestamp(t *testing.T) {
	t.Parallel()

	b := bytes.Buffer{}
	w := structured.New(&b, structured.WithTimestamp())

	_, err := w.Write([]byte("some log msg"))
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	d := json.NewDecoder(&b)
	d.UseNumber()

	var m map[string]interface{}
	err = d.Decode(&m)
	if err != nil {
		t.Fatalf("invalid json: %s", err)
	}

	ts, err := m["timestamp"].(json.Number).Int64()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	now := time.Now().UnixNano()
	if ts > now {
		t.Fatalf("expected %v to be less than the current time %v", ts, now)
	}
}

func TestWithCallSite(t *testing.T) {
	t.Parallel()

	b := bytes.Buffer{}
	w := structured.New(&b, structured.WithCallSite())

	_, err := w.Write([]byte("some log msg"))
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		t.Fatalf("invalid json: %s", err)
	}

	if match, _ := regexp.MatchString(`^[./a-zA-Z0-9]+`, m["callsite"].(string)); !match {
		t.Fatalf("expected callsite to be populated with filename:linenumber '%v'", m["callsite"])
	}
}

func TestUserProvidedJSON(t *testing.T) {
	t.Parallel()

	b := bytes.Buffer{}
	w := structured.New(&b, structured.WithCallSite())

	_, err := w.Write([]byte(`{"name":"metric-name","value":99.9}`))
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		t.Fatalf("invalid json: %s", err)
	}

	if m["msg"].(map[string]interface{})["name"].(string) != "metric-name" {
		t.Fatalf("expected %v to equal 'metric-name'", m["msg"].(map[string]interface{})["name"])
	}
}
