//go:build linux
// +build linux

package main

// NewKeyLogger creates a new keylogger depending on
// the platform we are running on (currently only Windows
// is supported)
func NewKeyLogger() KeyLogger {
	return &KeyLoggerLinux{}
}

// KeyLoggerWindows represents the keylogger
type KeyLoggerLinux struct {
	lastKey int
}

// GetKey gets the current entered key by the user, if there is any
func (kl *KeyLoggerLinux) GetKey() Key {

	return Key{Empty: true}
}
