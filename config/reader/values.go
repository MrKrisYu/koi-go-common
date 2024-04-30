package reader

import "time"

// Values is returned by the reader
type Values interface {
	Bytes() []byte
	Get(path ...string) Value
	Set(val interface{}, path ...string)
	Del(path ...string)
	Map() map[string]interface{}
	Scan(v interface{}) error
}

// Value represents a value of any type
type Value interface {
	Bool(defaultValue bool) bool
	Int(defaultValue int) int
	String(defaultValue string) string
	Float64(defaultValue float64) float64
	Duration(defaultValue time.Duration) time.Duration
	StringSlice(defaultValue []string) []string
	StringMap(defaultValue map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}
