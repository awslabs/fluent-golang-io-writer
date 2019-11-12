package logger

import (
	"encoding/json"
	"fmt"

	"github.com/fluent/fluent-logger-golang/fluent"
)

// FluentWriter implements io.Writer (and io.Closer)
type FluentWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
	AddSecondaryJSONConverter(f JSONConverter)
}

type fluentWriter struct {
	// Fluent Client
	fluentSender *fluent.Fluent
	// Prefix for that tag, this serves as a fallback to ensure that there is always a tag
	tagPrefix string
	// keys in the JSON whose values can be used to construct the tag. Keys are processed in order
	tagKeys []string
	// function to be used to convert the msg to map[string]string.
	converter JSONConverter
	// function to be used to construct a tag from the raw message
	tagConstructor TagConstructor
}

// JSONConverter allows users to provide a function to run on each log message to convert it to map[string]interface{}
type JSONConverter func(p []byte) (map[string]interface{}, error)

// TagConstructor allows users to provide a function to run on the raw message to construct a tag
// MUST return a non-empty value for the tag no matter what
type TagConstructor func(p []byte, tagPrefix string) string

// NewFluentWriter returns an io.Writer which will write to a
// local Fluent instance listening at address present in the config
// tagKeys are keys in the JSON log message whose values will be appended to tagPrefix to construct the tag
func NewFluentWriter(config fluent.Config, tagPrefix string, tagKeys []string) (FluentWriter, error) {
	logger, err := fluent.New(config)
	if err != nil {
		return nil, err
	}

	return &fluentWriter{
		fluentSender: logger,
		tagPrefix:    tagPrefix,
		tagKeys:      tagKeys,
	}, nil
}

// NewFluentWriterWithTagConstructor returns an io.Writer which will write to a
// local Fluent instance listening at address present in the config
// TagConstructor - a function which can be run on the raw bytes
// of each write to construct a tag
func NewFluentWriterWithTagConstructor(config fluent.Config, tagPrefix string, f TagConstructor) (FluentWriter, error) {
	logger, err := fluent.New(config)
	if err != nil {
		return nil, err
	}

	return &fluentWriter{
		fluentSender:   logger,
		tagPrefix:      tagPrefix,
		tagConstructor: f,
	}, nil
}

// AddJSONConverter adds a JSONConverter - a function
// which can convert each write to JSON
// This will be called if the log is not already JSON
// (it is called second, after json.Unmarshal fails)
func (w *fluentWriter) AddSecondaryJSONConverter(f JSONConverter) {
	w.converter = f
}

// Write function means this struct implements IO.Writer
func (w *fluentWriter) Write(p []byte) (n int, err error) {
	// try to unmarshal byte into map[string]interface{}, this means its JSON
	var msg map[string]interface{}
	err = json.Unmarshal(p, &msg)
	if err == nil {
		tag := w.constructTag(msg, p)
		return w.send(tag, msg, len(p))
	}

	// try the converter, if it was provided
	if w.converter != nil {
		msg, err := w.converter(p)
		if err == nil {
			tag := w.constructTag(msg, p)
			return w.send(tag, msg, len(p))
		}
	}

	// if its not JSON, construct a JSON string with the msg like:
	// {
	// 		"log": "original msg"
	// }
	msg["log"] = string(p)
	tag := w.tagPrefix
	if w.tagConstructor != nil {
		tag = w.tagConstructor(p, w.tagPrefix)
	}
	return w.send(tag, msg, len(p))
}

// Close closes the connection to Fluentd/Fluent Bit
func (w *fluentWriter) Close() error {
	return w.fluentSender.Close()
}

// Construct the tag based on the provided tagKeys
// Ex: tagKeys = ["level", "requestType"]
// Message = {"level": "info", "requestType": "post", "path": "/enroll"}
// tag will be "{tagPrefix}.info.post"
func (w *fluentWriter) constructTag(msg map[string]interface{}, p []byte) string {
	if w.tagConstructor != nil {
		return w.tagConstructor(p, w.tagPrefix)
	}

	tag := w.tagPrefix
	for _, key := range w.tagKeys {
		if val, ok := msg[key]; ok {
			tag = tag + "." + fmt.Sprintf("%v", val)
		}
	}
	return tag
}

func (w *fluentWriter) send(tag string, msg map[string]interface{}, length int) (n int, err error) {
	err = w.fluentSender.Post(tag, msg)
	if err != nil {
		// print to stdout as a fallback
		fmt.Printf("Failed to send: %v\n", msg)
		return length, nil
	}
	return 0, err
}
