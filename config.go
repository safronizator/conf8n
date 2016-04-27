// Package conf8n is here to simplify the way you work with your config files.
// It helps to avoid type-casting-hell, which occurs when you trying to parse your config files manually.
// Now you can load, read and traverse even complicated hierarchical configs a clear & simple way.
package conf8n

import (
	"fmt"
	"strings"
)

const (
	SEP = "."
)

// Base struct of the package. Represents loaded configuration.
type Config struct {
	data map[string]interface{}
}

// Represents value, got from config by given key or through iteration.
// Has methods to cast underlying interface value to concrete type.
type ConfigValue struct {
	v interface{}
}

type Iterator interface {
	Next()
	Finished() bool
	Index() int
	Key() string
	Value() *ConfigValue
}

type ListIterator struct {
	a []interface{}
	i int
}

type MapIterator struct {
	*ListIterator
	m map[string]interface{}
}

type EmptyIterator struct{}

// Base constructor for Config struct. Requires config data to be prepared as map[string]interface{}.
// For most cases you can use more high-level constructors (see docs for NewConfigFromYaml(),
// NewConfigFromJson() and NewConfigFromFile())
func NewConfig(fromData map[string]interface{}) *Config {
	return &Config{data: fromData}
}

// Get value by given key.
// Supports nested keys: for example, key "db.user" could be interpreted as is, if set;
// if not - system will lookup for value with key "user" in section with key "db"
func (c *Config) Get(key string) *ConfigValue {
	if v, ok := c.data[key]; ok {
		return &ConfigValue{v: v}
	}
	return &ConfigValue{v: getValueWithCompositeKey(c.data, strings.Split(key, SEP), 0)}
}

// Returns true if key was set and we has not nil value
func (v *ConfigValue) IsSet() bool {
	return v.v != nil
}

// Returns true if value can be interpreted as slice
func (v *ConfigValue) IsSlice() bool {
	_, is := v.v.([]interface{})
	return is
}

// Returns true if value can be interpreted as map
func (v *ConfigValue) IsMap() bool {
	if _, is := v.v.(map[string]interface{}); is {
		return true
	}
	_, is := v.v.(map[interface{}]interface{})
	return is
}

// Silently converts value to int - even if key was not set in config
func (v *ConfigValue) Int() int {
	i, _ := v.v.(int)
	return i
}

// Silently converts value to string
func (v *ConfigValue) String() string {
	s, _ := v.v.(string)
	return s
}

// Silently converts value to float
func (v *ConfigValue) Float() float64 {
	f, _ := v.v.(float64)
	return f
}

// Silently converts value to bool
func (v *ConfigValue) Bool() bool {
	b, _ := v.v.(bool)
	return b
}

// Returns underlying value withou casting (as interface{})
func (v *ConfigValue) Raw() interface{} {
	return v.v
}

// Get new Config instance from value
func (v *ConfigValue) Config() *Config {
	return NewConfig(toStrMap(v.v))
}

// Tries to cast value to int; reports error if key was not set or it was non int
func (v *ConfigValue) MustInt() (int, error) {
	if !v.IsSet() {
		return 0, fmt.Errorf("Value is not set")
	}
	if i, ok := v.v.(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("Value is not int")
}

// Tries to cast value to string; reports error if key was not set or it was non string
func (v *ConfigValue) MustString() (string, error) {
	if !v.IsSet() {
		return "", fmt.Errorf("Value is not set")
	}
	if s, ok := v.v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("Value is not string")
}

// Tries to cast value to float; reports error if key was not set or it was non float
func (v *ConfigValue) MustFloat() (float64, error) {
	if !v.IsSet() {
		return .0, fmt.Errorf("Value is not set")
	}
	if f, ok := v.v.(float64); ok {
		return f, nil
	}
	return .0, fmt.Errorf("Value is not float")
}

// Tries to cast value to bool; reports error if key was not set or it was non bool
func (v *ConfigValue) MustBool() (bool, error) {
	if !v.IsSet() {
		return false, fmt.Errorf("Value is not set")
	}
	if b, ok := v.v.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("Value is not bool")
}

// Tries to cast value to int. If it was not set, or can't be casted, returns given default value
func (v *ConfigValue) DefInt(def int) int {
	if i, ok := v.v.(int); ok {
		return i
	}
	return def
}

// Tries to cast value to string. If it was not set, or can't be casted, returns given default value
func (v *ConfigValue) DefString(def string) string {
	if s, ok := v.v.(string); ok {
		return s
	}
	return def
}

// Tries to cast value to float. If it was not set, or can't be casted, returns given default value
func (v *ConfigValue) DefFloat(def float64) float64 {
	if f, ok := v.v.(float64); ok {
		return f
	}
	return def
}

// Tries to cast value to bool. If it was not set, or can't be casted, returns given default value
func (v *ConfigValue) DefBool(def bool) bool {
	if b, ok := v.v.(bool); ok {
		return b
	}
	return def
}

// Tries to cast value to slice and return count of its elements. Returns 0 on failure
func (v *ConfigValue) Count() int {
	if a, ok := v.v.([]interface{}); ok {
		return len(a)
	}
	return 0
}

// Returns iterator for the value (if it was set as array).
//
// Example 1 (array key iteration):
// 	for i := config.Get("myArrayValue").Iterate(); !i.Finished(); i.Next() {
// 		fmt.Println(i.Value())
// 	}
//
// Example 2 (map key iteration):
// 	for i := config.Get("myMapValue").Iterate(); !i.Finished(); i.Next() {
// 		fmt.Println(i.Key(), ":", i.Value())
// 	}
func (v *ConfigValue) Iterate() Iterator {
	if a, ok := v.v.([]interface{}); ok {
		return &ListIterator{a, 0}
	}
	if m := toStrMap(v.v); m != nil {
		return &MapIterator{&ListIterator{mapGetKeys(m), 0}, m}
	}
	return &EmptyIterator{}
}

// See doc for ConfigValue.Iterate()
func (i *ListIterator) Next() {
	if !i.Finished() {
		i.i++
	}
}

// See doc for ConfigValue.Iterate()
func (i *ListIterator) Finished() bool {
	return i.i == len(i.a)
}

// See doc for ConfigValue.Iterate()
func (i *ListIterator) Value() *ConfigValue {
	return &ConfigValue{i.a[i.i]}
}

// Returns current iteration index
func (i *ListIterator) Index() int {
	return i.i
}

// Always return empty string
func (i *ListIterator) Key() string {
	return ""
}

// See doc for ConfigValue.Iterate()
func (i *MapIterator) Value() *ConfigValue {
	return &ConfigValue{i.m[i.Key()]}
}

// Return current key
func (i *MapIterator) Key() string {
	return i.ListIterator.Value().String()
}

// No action here
func (i *EmptyIterator) Next() {
	//pass
}

// Always return empty value
func (i *EmptyIterator) Value() *ConfigValue {
	return &ConfigValue{nil}
}

// Always return true
func (i *EmptyIterator) Finished() bool {
	return true
}

// Always return 0
func (i *EmptyIterator) Index() int {
	return 0
}

// Always return empty string
func (i *EmptyIterator) Key() string {
	return ""
}

func mapGetKeys(m map[string]interface{}) []interface{} {
	a := make([]interface{}, 0, len(m))
	for k, _ := range m {
		a = append(a, k)
	}
	return a
}
