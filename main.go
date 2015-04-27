package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/fsnotify.v1"
)

var (
	filesFlag = flag.String("f", "", "files and directories to watch, split by ':'")
	sigFlag   = flag.String("s", "HUP", "signal to send on change")
)

var signals = map[string]os.Signal{
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"BUS":    syscall.SIGBUS,
	"CHLD":   syscall.SIGCHLD,
	"CONT":   syscall.SIGCONT,
	"EMT":    syscall.SIGEMT,
	"FPE":    syscall.SIGFPE,
	"HUP":    syscall.SIGHUP,
	"ILL":    syscall.SIGILL,
	"INFO":   syscall.SIGINFO,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGIOT,
	"KILL":   syscall.SIGKILL,
	"PIPE":   syscall.SIGPIPE,
	"PROF":   syscall.SIGPROF,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STOP":   syscall.SIGSTOP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"TRAP":   syscall.SIGTRAP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"WINCH":  syscall.SIGWINCH,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
}

func signalByName(name string) (sig os.Signal, err error) {
	var ok bool
	if sig, ok = signals[strings.ToUpper(name)]; !ok {
		err = errors.New(fmt.Sprintf("unknown signal: %s", name))
	}
	return
}

func usageAndExit(s interface{}) {
	fmt.Println(s)
	flag.Usage()
	os.Exit(1)
}

func init() {
	flag.Parse()
}

func main() {
	sig, err := signalByName(*sigFlag)
	if err != nil {
		usageAndExit(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		usageAndExit(err)
	}
	defer watcher.Close()

	if *filesFlag == "" {
		usageAndExit("no files to watch")
	}

	if len(flag.Args()) == 0 {
		usageAndExit("no command specified")
	}

	// start watching files for changes
	for _, file := range strings.Split(*filesFlag, ":") {
		err = watcher.Add(file)
		// if a file doesn't exist that you're trying to watch at this
		// point, it's likely a config error, and we should bail
		if err != nil {
			usageAndExit(err)
		}
	}

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		usageAndExit(err)
	}

	fmt.Println("==> starting child", strings.Join(flag.Args(), " "))

	done := make(chan struct{})
	go func() {
		cmd.Wait()
		done <- struct{}{}
	}()

	go func() {
		var event fsnotify.Event
		var err error
		for {
			select {
			case event = <-watcher.Events:
				log.Println("==> detected change in", event.Name)
				// TODO: add throttle/debounce behavior to prevent a ton
				// of signals being sent.
				cmd.Process.Signal(sig)
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					// File was renamed, so remove the old watch,
					// and add a new one
					watcher.Remove(event.Name)
					watcher.Add(event.Name)
				}
			case err = <-watcher.Errors:
				log.Println("==> error:", err)
			}
		}
	}()

	// Listen to signals send to parent, and pass along to the child
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		var sig os.Signal
		for {
			sig = <-c
			cmd.Process.Signal(sig)
			switch sig {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				countdown := 5 * time.Second

				fmt.Println("==> attempting to shut down cleanly")
				fmt.Printf("==> waiting up to %s for child to exit\n", countdown)

				// try and wait for the child to shut down before killing
				select {
				case <-done:
				case <-time.After(countdown):
				}

				// it hasn't shut down yet, so attempt to SIGKILL
				fmt.Println("==> attempting to now kill child")
				cmd.Process.Signal(os.Kill)

				select {
				case <-done:
				case <-time.After(time.Second):
				}

				// still hasn't exited, so killing self
				fmt.Println("==> now committing suicide")
				os.Exit(1)
			}
		}
	}()

	<-done
}