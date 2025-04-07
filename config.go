package deepl

import (
	"encoding/json"
	"time"
)

type JSONMarshal func(v any) ([]byte, error)
type JSONUnmarshaler func(data []byte, v any) error

type Config struct {
	AuthKey     string        // deepl api authKey
	Timeout     time.Duration // request timeout
	AccountType int           // deepl account type free|pro
	JSONEncode  JSONMarshal
	JSONDecode  JSONUnmarshaler
}

var DefaultConfig = Config{
	Timeout:     10 * time.Second,
	AccountType: FreeAccount,
	JSONEncode:  json.Marshal,
	JSONDecode:  json.Unmarshal,
}
