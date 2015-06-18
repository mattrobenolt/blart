package main

import (
	"os"
	"syscall"
)

var signals = map[string]os.Signal{
	"ABRT": syscall.SIGABRT,
	"ALRM": syscall.SIGALRM,
	"BUS":  syscall.SIGBUS,
	"FPE":  syscall.SIGFPE,
	"HUP":  syscall.SIGHUP,
	"ILL":  syscall.SIGILL,
	"INT":  syscall.SIGINT,
	"KILL": syscall.SIGKILL,
	"PIPE": syscall.SIGPIPE,
	"QUIT": syscall.SIGQUIT,
	"SEGV": syscall.SIGSEGV,
	"TERM": syscall.SIGTERM,
	"TRAP": syscall.SIGTRAP,
}
