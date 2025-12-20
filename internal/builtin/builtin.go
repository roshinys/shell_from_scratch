package builtin

// Builtin commands lookup
var Builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

var ExternalBuiltins = make(map[string]bool)


func IsBuiltin(cmd string) bool {
	return Builtins[cmd]
}