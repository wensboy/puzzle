package util

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

const (
	tag      = "check"
	ruleSep  = ","
	valueSep = "="
)

var (
	vaOnce          sync.Once
	globalValidator *validator
)

type (
	vErr struct {
		Field string
		Type  string
		Value interface{}
		Rule  string
	}

	vErrs    []vErr
	RuleFunc func(v, want interface{}) bool

	validator struct {
		e     *echo.Echo
		rules map[string]RuleFunc
		mu    sync.RWMutex
	}
	Validator interface {
		Register(k string, fn RuleFunc)
		UnRegister(k string)
		SetupDefaultRules()
		Check(s interface{}) vErrs
	}

	Rule struct {
		key  string
		want string
	}
)

func (ve vErr) Error() string {
	return fmt.Sprintf("[%s %s %s %v]", ve.Rule, ve.Type, ve.Field, ve.Value)
}

func (ves vErrs) Error() string {
	var buf bytes.Buffer
	for i, ve := range ves {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(ve.Error())
	}
	return buf.String()
}

func GetGlobalValidator() Validator {
	return globalValidator
}

func InitValidator(e *echo.Echo) {
	globalValidator = NewValidator(e)
	globalValidator.SetupDefaultRules()
}

func NewValidator(e *echo.Echo) *validator {
	return &validator{
		e:     e,
		rules: make(map[string]RuleFunc),
	}
}

func (va *validator) SetupDefaultRules() {
	ruleName := []string{"request", "min", "max", "email", "uuid", "fixed"}
	rulefunc := []RuleFunc{Request, Min, Max, Email, UUID, Fixed}
	for i := range ruleName {
		va.Register(ruleName[i], rulefunc[i])
	}
}

func Request(v, want interface{}) bool {
	vr := reflect.ValueOf(v)
	if !vr.IsValid() {
		return false
	}
	switch vr.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return !vr.IsNil()
	}
	return true
}

func isNumber(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isInteger(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func Min(v, want interface{}) bool {
	vv := reflect.ValueOf(v)
	ww := reflect.ValueOf(want)

	if !vv.IsValid() || !ww.IsValid() {
		return false
	}
	if ww.Kind() != reflect.String {
		return false
	}
	wv, err := strconv.ParseUint(ww.String(), 10, 64)
	if err != nil {
		return false
	}
	switch vv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vv.Int() >= int64(wv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vv.Uint() >= wv
	case reflect.String:
		return len(vv.String()) >= int(wv)
	default:
		return false
	}
}

func Max(v, want interface{}) bool {
	vv := reflect.ValueOf(v)
	ww := reflect.ValueOf(want)

	if !vv.IsValid() || !ww.IsValid() {
		return false
	}
	if ww.Kind() != reflect.String {
		return false
	}
	wv, err := strconv.ParseUint(ww.String(), 10, 64)
	if err != nil {
		return false
	}
	switch vv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vv.Int() <= int64(wv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vv.Uint() <= wv
	case reflect.String:
		return len(vv.String()) <= int(wv)
	default:
		return false
	}
}

func Email(v, want interface{}) bool {
	vv := reflect.ValueOf(v)
	if !vv.IsValid() {
		return false
	}
	if vv.Kind() != reflect.String {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z]{2,})+$`)
	return re.MatchString(vv.String())
}

func UUID(v, want interface{}) bool {
	vv := reflect.ValueOf(v)
	if !vv.IsValid() {
		return false
	}
	if vv.Kind() != reflect.String {
		return false
	}
	re := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return re.MatchString(vv.String())
}

func Fixed(v, want interface{}) bool {
	vv := reflect.ValueOf(v)
	ww := reflect.ValueOf(want)
	if ww.Kind() != reflect.String {
		return false
	}
	wv, _ := strconv.ParseInt(ww.String(), 10, 64)
	if !vv.IsValid() {
		return false
	}
	if vv.Kind() != reflect.String {
		return false
	}
	return int64(len(vv.String())) == wv
}

func (va *validator) apply(k string) RuleFunc {
	va.mu.RLock()
	defer va.mu.RUnlock()
	if fn, found := va.rules[k]; !found {
		panic("[Validator] - Non-existent verification processing")
	} else {
		return fn
	}
}

func (va *validator) Register(k string, fn RuleFunc) {
	va.mu.Lock()
	defer va.mu.Unlock()
	va.rules[k] = fn
}

func (va *validator) UnRegister(k string) {
	va.mu.Lock()
	defer va.mu.Unlock()
	delete(va.rules, k)
}

func (va *validator) parseRules(ruleStr string) []Rule {
	ruleParts := strings.SplitN(ruleStr, ruleSep, -1)
	var rules []Rule
	for i := range ruleParts {
		part := strings.SplitN(ruleParts[i], valueSep, 2)
		k, w := part[0], ""
		if len(part) == 2 {
			w = part[1]
		}
		rules = append(rules, Rule{
			key:  k,
			want: w,
		})
	}
	return rules
}

func (va *validator) Check(s interface{}) vErrs {
	var errs vErrs
	// TODO: check tag from rules to v -> struct
	t, v := reflect.TypeOf(s), reflect.ValueOf(s)
	if t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("[validator] - unexpected validation type")
	}
	for i := 0; i < t.NumField(); i += 1 {
		field := t.Field(i)
		value := v.Field(i)
		tName := value.Kind().String()
		ruleStr := field.Tag.Get(tag)
		if ruleStr == "" {
			// skip no check tag field
			continue
		}
		rules := va.parseRules(ruleStr)
		for j := range rules {
			if !va.apply(rules[j].key)(value.Interface(), rules[j].want) {
				errs = append(errs, vErr{
					Field: field.Name,
					Type:  tName,
					Value: value.Interface(),
					Rule:  rules[j].key,
				})
			}
		}
	}
	return errs
}
