package mdparser

import (
	"time"
)

type FrontmatterProcessor struct {
	data map[string]interface{}
}

func NewFrontmatterProcessor(data map[string]interface{}) *FrontmatterProcessor {
	return &FrontmatterProcessor{data: data}
}

func (fp *FrontmatterProcessor) GetString(key string) (string, bool) {
	if val, ok := fp.data[key].(string); ok {
		return val, true
	}
	return "", false
}

func (fp *FrontmatterProcessor) GetBool(key string) (bool, bool) {
	if val, ok := fp.data[key].(bool); ok {
		return val, true
	}
	return false, false
}

func (fp *FrontmatterProcessor) GetTime(key string) (time.Time, bool) {
	if strVal, ok := fp.GetString(key); ok {
		if timeVal, err := time.Parse(time.RFC3339, strVal); err == nil {
			return timeVal, true
		}
	}
	return time.Time{}, false
}
