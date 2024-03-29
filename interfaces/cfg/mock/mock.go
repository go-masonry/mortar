// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_cfg is a generated GoMock package.
package mock_cfg

import (
	reflect "reflect"
	time "time"

	cfg "github.com/go-masonry/mortar/interfaces/cfg"
	gomock "github.com/golang/mock/gomock"
)

// MockValue is a mock of Value interface.
type MockValue struct {
	ctrl     *gomock.Controller
	recorder *MockValueMockRecorder
}

// MockValueMockRecorder is the mock recorder for MockValue.
type MockValueMockRecorder struct {
	mock *MockValue
}

// NewMockValue creates a new mock instance.
func NewMockValue(ctrl *gomock.Controller) *MockValue {
	mock := &MockValue{ctrl: ctrl}
	mock.recorder = &MockValueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValue) EXPECT() *MockValueMockRecorder {
	return m.recorder
}

// Bool mocks base method.
func (m *MockValue) Bool() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bool")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Bool indicates an expected call of Bool.
func (mr *MockValueMockRecorder) Bool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bool", reflect.TypeOf((*MockValue)(nil).Bool))
}

// Duration mocks base method.
func (m *MockValue) Duration() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Duration")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Duration indicates an expected call of Duration.
func (mr *MockValueMockRecorder) Duration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Duration", reflect.TypeOf((*MockValue)(nil).Duration))
}

// Float64 mocks base method.
func (m *MockValue) Float64() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Float64")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Float64 indicates an expected call of Float64.
func (mr *MockValueMockRecorder) Float64() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Float64", reflect.TypeOf((*MockValue)(nil).Float64))
}

// Int mocks base method.
func (m *MockValue) Int() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Int")
	ret0, _ := ret[0].(int)
	return ret0
}

// Int indicates an expected call of Int.
func (mr *MockValueMockRecorder) Int() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int", reflect.TypeOf((*MockValue)(nil).Int))
}

// Int32 mocks base method.
func (m *MockValue) Int32() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Int32")
	ret0, _ := ret[0].(int32)
	return ret0
}

// Int32 indicates an expected call of Int32.
func (mr *MockValueMockRecorder) Int32() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int32", reflect.TypeOf((*MockValue)(nil).Int32))
}

// Int64 mocks base method.
func (m *MockValue) Int64() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Int64")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Int64 indicates an expected call of Int64.
func (mr *MockValueMockRecorder) Int64() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int64", reflect.TypeOf((*MockValue)(nil).Int64))
}

// IntSlice mocks base method.
func (m *MockValue) IntSlice() []int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IntSlice")
	ret0, _ := ret[0].([]int)
	return ret0
}

// IntSlice indicates an expected call of IntSlice.
func (mr *MockValueMockRecorder) IntSlice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IntSlice", reflect.TypeOf((*MockValue)(nil).IntSlice))
}

// IsSet mocks base method.
func (m *MockValue) IsSet() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSet")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSet indicates an expected call of IsSet.
func (mr *MockValueMockRecorder) IsSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSet", reflect.TypeOf((*MockValue)(nil).IsSet))
}

// Raw mocks base method.
func (m *MockValue) Raw() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Raw")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Raw indicates an expected call of Raw.
func (mr *MockValueMockRecorder) Raw() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Raw", reflect.TypeOf((*MockValue)(nil).Raw))
}

// String mocks base method.
func (m *MockValue) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockValueMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockValue)(nil).String))
}

// StringMap mocks base method.
func (m *MockValue) StringMap() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringMap")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// StringMap indicates an expected call of StringMap.
func (mr *MockValueMockRecorder) StringMap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringMap", reflect.TypeOf((*MockValue)(nil).StringMap))
}

// StringMapString mocks base method.
func (m *MockValue) StringMapString() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringMapString")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// StringMapString indicates an expected call of StringMapString.
func (mr *MockValueMockRecorder) StringMapString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringMapString", reflect.TypeOf((*MockValue)(nil).StringMapString))
}

// StringMapStringSlice mocks base method.
func (m *MockValue) StringMapStringSlice() map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringMapStringSlice")
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// StringMapStringSlice indicates an expected call of StringMapStringSlice.
func (mr *MockValueMockRecorder) StringMapStringSlice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringMapStringSlice", reflect.TypeOf((*MockValue)(nil).StringMapStringSlice))
}

// StringSlice mocks base method.
func (m *MockValue) StringSlice() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringSlice")
	ret0, _ := ret[0].([]string)
	return ret0
}

// StringSlice indicates an expected call of StringSlice.
func (mr *MockValueMockRecorder) StringSlice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringSlice", reflect.TypeOf((*MockValue)(nil).StringSlice))
}

// Time mocks base method.
func (m *MockValue) Time() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Time")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Time indicates an expected call of Time.
func (mr *MockValueMockRecorder) Time() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Time", reflect.TypeOf((*MockValue)(nil).Time))
}

// Uint mocks base method.
func (m *MockValue) Uint() uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uint")
	ret0, _ := ret[0].(uint)
	return ret0
}

// Uint indicates an expected call of Uint.
func (mr *MockValueMockRecorder) Uint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uint", reflect.TypeOf((*MockValue)(nil).Uint))
}

// Uint32 mocks base method.
func (m *MockValue) Uint32() uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uint32")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// Uint32 indicates an expected call of Uint32.
func (mr *MockValueMockRecorder) Uint32() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uint32", reflect.TypeOf((*MockValue)(nil).Uint32))
}

// Uint64 mocks base method.
func (m *MockValue) Uint64() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uint64")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Uint64 indicates an expected call of Uint64.
func (mr *MockValueMockRecorder) Uint64() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uint64", reflect.TypeOf((*MockValue)(nil).Uint64))
}

// Unmarshal mocks base method.
func (m *MockValue) Unmarshal(result interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", result)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockValueMockRecorder) Unmarshal(result interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockValue)(nil).Unmarshal), result)
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConfig) Get(key string) cfg.Value {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(cfg.Value)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockConfigMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfig)(nil).Get), key)
}

// Implementation mocks base method.
func (m *MockConfig) Implementation() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Implementation")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Implementation indicates an expected call of Implementation.
func (mr *MockConfigMockRecorder) Implementation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Implementation", reflect.TypeOf((*MockConfig)(nil).Implementation))
}

// Map mocks base method.
func (m *MockConfig) Map() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Map")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Map indicates an expected call of Map.
func (mr *MockConfigMockRecorder) Map() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Map", reflect.TypeOf((*MockConfig)(nil).Map))
}

// Set mocks base method.
func (m *MockConfig) Set(key string, value interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value)
}

// Set indicates an expected call of Set.
func (mr *MockConfigMockRecorder) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockConfig)(nil).Set), key, value)
}

// MockBuilder is a mock of Builder interface.
type MockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderMockRecorder
}

// MockBuilderMockRecorder is the mock recorder for MockBuilder.
type MockBuilderMockRecorder struct {
	mock *MockBuilder
}

// NewMockBuilder creates a new mock instance.
func NewMockBuilder(ctrl *gomock.Controller) *MockBuilder {
	mock := &MockBuilder{ctrl: ctrl}
	mock.recorder = &MockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilder) EXPECT() *MockBuilderMockRecorder {
	return m.recorder
}

// AddExtraConfigFile mocks base method.
func (m *MockBuilder) AddExtraConfigFile(path string) cfg.Builder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddExtraConfigFile", path)
	ret0, _ := ret[0].(cfg.Builder)
	return ret0
}

// AddExtraConfigFile indicates an expected call of AddExtraConfigFile.
func (mr *MockBuilderMockRecorder) AddExtraConfigFile(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddExtraConfigFile", reflect.TypeOf((*MockBuilder)(nil).AddExtraConfigFile), path)
}

// Build mocks base method.
func (m *MockBuilder) Build() (cfg.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build")
	ret0, _ := ret[0].(cfg.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build.
func (mr *MockBuilderMockRecorder) Build() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockBuilder)(nil).Build))
}

// SetConfigFile mocks base method.
func (m *MockBuilder) SetConfigFile(path string) cfg.Builder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConfigFile", path)
	ret0, _ := ret[0].(cfg.Builder)
	return ret0
}

// SetConfigFile indicates an expected call of SetConfigFile.
func (mr *MockBuilderMockRecorder) SetConfigFile(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfigFile", reflect.TypeOf((*MockBuilder)(nil).SetConfigFile), path)
}

// SetEnvDelimiterReplacer mocks base method.
func (m *MockBuilder) SetEnvDelimiterReplacer(from, to string) cfg.Builder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetEnvDelimiterReplacer", from, to)
	ret0, _ := ret[0].(cfg.Builder)
	return ret0
}

// SetEnvDelimiterReplacer indicates an expected call of SetEnvDelimiterReplacer.
func (mr *MockBuilderMockRecorder) SetEnvDelimiterReplacer(from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEnvDelimiterReplacer", reflect.TypeOf((*MockBuilder)(nil).SetEnvDelimiterReplacer), from, to)
}
