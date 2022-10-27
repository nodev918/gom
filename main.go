package main

import (
        "io"
        "log"
        "os"
        "os/exec"
        "os/signal"
        "syscall"

        "github.com/creack/pty"
        "golang.org/x/term"
)

func test() error {
        // Create arbitrary command.
        c := exec.Command("bash")

        // Start the command with a pty.
        ptmx, err := pty.Start(c)
        if err != nil {
                return err
        }
        // Make sure to close the pty at the end.
        defer func() { _ = ptmx.Close() }() // Best effort.

        // Handle pty size.
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGWINCH)
        go func() {
                for range ch {
                        if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
                                log.Printf("error resizing pty: %s", err)
                        }
                }
        }()
        ch <- syscall.SIGWINCH // Initial resize.
        defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

        // Set stdin in raw mode.
        oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
                panic(err)
        }
        defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

        // Copy stdin to the pty and the pty to stdout.
        // NOTE: The goroutine will keep reading until the next keystroke before returning.
        go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
        _, _ = io.Copy(os.Stdout, ptmx)

        return nil
}

func main() {
        if err := test(); err != nil {
                log.Fatal(err)
        }
}


// package main

// import (
// 	"io"
// 	"os"
// 	"os/exec"

// 	"github.com/creack/pty"
// )

// func main() {
// 	c := exec.Command("grep", "--color=auto", "bar")
// 	f, err := pty.Start(c)
// 	if err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		f.Write([]byte("foo\n"))
// 		f.Write([]byte("bar\n"))
// 		f.Write([]byte("baz\n"))
// 		f.Write([]byte{4}) // EOT
// 	}()
// 	io.Copy(os.Stdout, f)
// }



// package main

// import (
// 	"os"
// 	"os/exec"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/layout"
// 	"fyne.io/fyne/v2/widget"
// 	"github.com/creack/pty"
// )

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("germ")

// 	ui := widget.NewTextGrid()       // Create a new TextGrid
// 	ui.SetText("I'm on a terminal!") // Set text to display

// 	c := exec.Command("/bin/bash")
// 	p, err := pty.Start(c)

// 	if err != nil {
// 		fyne.LogError("Failed to open pty", err)
// 		os.Exit(1)
// 	}

// 	defer c.Process.Kill()

// 	p.Write([]byte("ls\r"))
// 	time.Sleep(1 * time.Second)
// 	b := make([]byte, 1024)
// 	_, err = p.Read(b)
// 	if err != nil {
// 		fyne.LogError("Failed to read pty", err)
// 	}
// 	// s := fmt.Sprintf("read bytes from pty.\nContent:%s",  string(b))
// 	ui.SetText(string(b))
// 	// Create a new container with a wrapped layout
// 	// set the layout width to 420, height to 200
// 	w.SetContent(
// 		fyne.NewContainerWithLayout(
// 			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
// 			ui,
// 		),
// 	)

// 	w.ShowAndRun()

// }



// package main

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/layout"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	a := app.New()
// 	w := a.NewWindow("germ")
// 	ui := widget.NewTextGrid()
// 	ui.SetText("I'm on a terminal")

// 	w.SetContent(
// 		fyne.NewContainerWithLayout(
// 			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
// 			ui,
// 		),
// 	)

// 	w.ShowAndRun()
// }
