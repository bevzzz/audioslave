package main

type KeyLogger interface {
	GetKey() Key
}

// Key is a single key entered by the user
type Key struct {
	Empty   bool
	Rune    rune
	Keycode int
}
