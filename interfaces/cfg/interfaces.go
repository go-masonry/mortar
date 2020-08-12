package cfg

import "time"

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// Value is a convenient interface to cast value to different types
type Value interface {
	// IsSet will tell if this key really exists in the configuration
	IsSet() bool
	// Raw returns an interface{}. For a specific type value use a corresponding method.
	Raw() interface{}
	// Bool returns the value associated with the key as a boolean.
	Bool() bool
	// Int returns the value associated with the key as an integer.
	Int() int
	// Int32 returns the value associated with the key as an integer.
	Int32() int32
	// Int64 returns the value associated with the key as an integer.
	Int64() int64
	// Uint returns the value associated with the key as an unsigned integer.
	Uint() uint
	// Uint32 returns the value associated with the key as an unsigned integer.
	Uint32() uint32
	// Uint64 returns the value associated with the key as an unsigned integer.
	Uint64() uint64
	// Float64 returns the value associated with the key as a float64.
	Float64() float64
	// Time returns the value associated with the key as time.
	Time() time.Time
	// Duration returns the value associated with the key as a duration.
	Duration() time.Duration
	// String returns the value associated with the key as a string.
	String() string
	// IntSlice returns the value associated with the key as a slice of int values.
	IntSlice() []int
	// StringSlice returns the value associated with the key as a slice of strings.
	StringSlice() []string
	// StringMap returns the value associated with the key as a map of interfaces.
	StringMap() map[string]interface{}
	// StringMapString returns the value associated with the key as a map of strings.
	StringMapString() map[string]string
	// StringMapStringSlice returns the value associated with the key as a map to a slice of strings.
	StringMapStringSlice() map[string][]string
	// Unmarshal tries to unmarshal it to a 'result'. 'result' field must be a pointer.
	//
	// It heavy depends on what library is used to provide Config. For example Viper uses 'mapstructure' for that
	Unmarshal(result interface{}) error
}

// Config defines an interface to obtain configuration values from JSON/YAML/TOML or ENV
type Config interface {
	/*
		Get returns a Value associated with a given key, you can later cast this to a type.

		Examples:

		- Get an Int value

				numberOfPossibilities := config.Get("path.to.key").Int() // if key is absent this will return 0

		- It's possible to check if there is an actual value associated with this key

				initTime := time.Now() // default value
				if value := config.Get("path.to.key"); value.IsSet() {
					initTime = value.Time()
				}
	*/
	Get(key string) Value
	// Set sets the value for the key within the Config map, it will take precedence over any other levels (ENV,flags,etc...)
	Set(key string, value interface{})
	// Map the entire configuration to a... map
	Map() map[string]interface{}
	// Implementation returns the actual lib/struct that is responsible for the above logic
	Implementation() interface{}
}

// Builder defines configuration builder options
type Builder interface {
	// SetConfigFile tells builder where to look for file with the configuration map
	SetConfigFile(path string) Builder
	// AddExtraConfigFile allows to add additional files to be merged into the configuration map.
	// Very useful for test/staging/dev environments where you want to override default values with env specific ones
	AddExtraConfigFile(path string) Builder
	// SetEnvDelimiterReplacer allows to customize Environment variable delimiter replacer.
	//
	// By default ENVIRONMENT Variables look something similar to this: CONFIG_FOLDER or SCHEDULER_DEFAULTS_TIMEOUT
	// Config implementations allows (and encourage) you to treat config as a MAP.
	// To access sub key values it is common to use this syntax:
	//		scheduler.defaults.timeout
	// Now if you want ot override this config value with Environment one SCHEDULER_DEFAULTS_TIMEOUT you need to replace
	// '_' with '.'
	SetEnvDelimiterReplacer(from, to string) Builder
	// Build returns Config implementation
	Build() (Config, error)
}
