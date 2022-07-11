package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AI1411/golang-admin-api/util"
	"github.com/AI1411/golang-admin-api/util/errors"
	"strconv"
	"time"
)

type parser struct {
	err error
}

func convertToUint64(value interface{}) (uint64, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}

func (p *parser) parseString(value interface{}) *string {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		return util.StringToPtr(v)
	case json.Number:
		return util.StringToPtr(v.String())
	default:
		return util.StringToPtr(fmt.Sprint(v))
	}
}

func (p *parser) parseUint8(value interface{}) *uint8 {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	u64, err := convertToUint64(value)
	if err != nil {
		p.err = err
		return nil
	}
	return util.Uint8ToPtr(uint8(u64))
}

func (p *parser) parseUint32(value interface{}) *uint32 {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	u64, err := convertToUint64(value)
	if err != nil {
		p.err = err
		return nil
	}
	return util.Uint32ToPtr(uint32(u64))
}

func (p *parser) parseUint64(value interface{}) *uint64 {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	u64, err := convertToUint64(value)
	if err != nil {
		p.err = err
		return nil
	}
	return util.Uint64ToPtr(u64)
}

func (p *parser) parseInt(value interface{}) *int {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	u64, err := convertToUint64(value)
	if err != nil {
		p.err = err
		return nil
	}
	i := int(u64)
	return &i
}

func (p *parser) parseBool(value interface{}) *bool {
	if value == "" || value == nil || p.err != nil {
		return nil
	}
	strToBoolPtr := func(s string) *bool {
		b, err := strconv.ParseBool(s)
		if err != nil {
			p.err = err
			return nil
		}
		return util.BoolToPtr(b)
	}
	switch v := value.(type) {
	case bool:
		return util.BoolToPtr(v)
	case string:
		return strToBoolPtr(v)
	case json.Number:
		return strToBoolPtr(v.String())
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		return util.BoolToPtr(u64 != 0)
	default:
		msg := fmt.Sprintf("%v cannot convert to bool", v)
		p.err = errors.New(msg)
		return nil
	}
}

func (p *parser) parseTime(value string) time.Time {
	if value == "" || p.err != nil {
		return time.Time{}
	}
	v, err := time.Parse(time.RFC3339, value)
	if err != nil {
		p.err = err
		return v
	}
	return v
}

func (p *parser) parsePtrTime(value string) *time.Time {
	if value == "" || p.err != nil {
		return &time.Time{}
	}
	v, err := time.Parse(time.RFC3339, value)
	if err != nil {
		p.err = err
		return &v
	}
	return &v
}
