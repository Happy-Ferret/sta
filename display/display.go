/*
Package display handles all IO operations in sta.
*/
package display

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"regexp"
)

// terminal escape sequences
var (
	clr   = "\033[0;0H\033[2J"
	rev   = "\033[1;7m"
	hl    = "\033[3m"
	norm  = "\033[0m"
	up    = "\033A"
	down  = "\033B"
	right = "\033C"
	left  = "\033D"
)

// Display is a type for handling i/o in sta.
type Display struct {
	sess ssh.Session
	term *terminal.Terminal
}

// New returns a pointer to a *Display.
func New(sess ssh.Session) *Display {
	d := &Display{sess, terminal.NewTerminal(sess, "")}
	if _, chW, ok := d.sess.Pty(); ok {
		go func() {
			for {
				select {
				case <-chW:
					d.Clear()
				}
			}
		}()
	} else {
		return nil
	}
	return d
}

// Clear erases the screen and prints a frame around it
func (d *Display) Clear() error {
	d.sess.Write([]byte(clr))
	if rows, cols, ok := d.GetSize(); ok {
		err := d.frame(0, 0, rows, cols)
		if err != nil {
			return err
		}
		d.gotoyx(rows/2, 0)
		return nil
	} else {
		return errors.New("Impossible to initialize pty.")
	}
}

// GetSize returns the display’s current height and width
func (d *Display) GetSize() (rows, cols int, ok bool) {
	if pty, _, ok := d.sess.Pty(); ok {
		return pty.Window.Height, pty.Window.Width, true
	}
	return 0, 0, false
}

// WriteLine writes a string to the screen
func (d *Display) WriteLine(str string) error {
	if !d.gotoyx(1, 1) {
		return errors.New("Impossible to write at given position.")
	}
	_, err := d.sess.Write([]byte(str + "\n"))
	return err
}

// ReadLine reads a string from stdin.
func (d *Display) ReadLine(prompt string) (string, error) {
	d.term.SetPrompt(rev + prompt + " " + hl + ">" + norm + " ")
	return d.term.ReadLine()
}

// CompleteWith sets autocompletion of line with the given slice of strings.
func (d *Display) CompleteWith(cmds []string) {
	d.term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// only autocomplete on TAB
		if key != '\t' {
			return
		}

		found := false
		for _, cmd := range cmds {
			if matched, err := regexp.Match("^"+line+".*", []byte(cmd)); err == nil && matched {
				// if two commands at least match, do not complete
				if found {
					return "", 0, false
				}
				found = true
				newLine = cmd + " "
				newPos = len(cmd) + 1
				ok = true
			}
		}
		// If we get here and found is true, we found something and
		// set the output variables; if nothing was found, the output
		// variables are still nil.
		// In both cases there is only one thing left to do:
		return
	}
}

// gotoyx sends cursor to specified location.
func (d *Display) gotoyx(y, x int) bool {
	if !d.yxOK(y, x) {
		return false
	}
	d.sess.Write([]byte(fmt.Sprintf("\033[%v;%vH", y, x)))
	return true
}

// yxOK indicates if given coordinates are inside the window.
func (d *Display) yxOK(y, x int) bool {
	if rows, cols, ok := d.GetSize(); ok && y >= 0 && y < rows && x >= 0 && x <= cols {
		return true
	}
	return false
}

// frame draws a frame around a given region of the screen.
func (d *Display) frame(y, x, h, w int) error {
	// exit if out of bounds
	if !d.yxOK(y, x) || !d.yxOK(y+h-1, x+w-1) {
		return errors.New("Impossible to draw frame.")
	}

	// set reverse
	_, err := d.sess.Write([]byte(rev))
	if err != nil {
		return err
	}

	// draw top line
	d.gotoyx(y, x)
	for i := x; i < x+w; i++ {
		_, err = d.sess.Write([]byte("#"))
		if err != nil {
			return err
		}
	}

	// draw bottom line
	d.gotoyx(y+h-1, x)
	for i := x; i < x+w; i++ {
		_, err = d.sess.Write([]byte("#"))
		if err != nil {
			return err
		}
	}

	// draw left line
	for i := y; i < y+h; i++ {
		d.gotoyx(i, x)
		_, err = d.sess.Write([]byte("#"))
		if err != nil {
			return err
		}
	}

	// draw right line
	for i := y; i < y+h; i++ {
		d.gotoyx(i, x+w-1)
		_, err = d.sess.Write([]byte("#"))
		if err != nil {
			return err
		}
	}

	_, err = d.sess.Write([]byte(norm))
	if err != nil {
		return err
	}
	d.gotoyx(1, 1)
	return nil
}
