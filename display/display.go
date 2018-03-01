/*
Package display handles all IO operations in sta.
*/
package display

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"regexp"
)

// Display is a type for handling i/o in sta.
type Display struct {
	sess ssh.Session
	term *terminal.Terminal
}

// New returns a pointer to a *Display.
func New(sess ssh.Session) *Display {
	d := &Display{sess, terminal.NewTerminal(sess, "")}
	// action when window size changes
	/*if _, chW, ok := d.sess.Pty(); ok {
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
	}*/
	return d
}

/*/ Clear erases the screen.
func (d *Display) Clear() error {
	_, err := d.term.Write([]byte(clr))
	return err
}*/

// gotoyx sends cursor to specified location.
func (d *Display) gotoyx(y, x int) bool {
	if !d.okyx(y, x) {
		return false
	}
	d.term.Write([]byte(fmt.Sprintf("\033[%v;%vH", y, x)))
	return true
}

// okyx indicates if given coordinates are inside the window.
func (d *Display) okyx(y, x int) bool {
	if rows, cols, ok := d.GetSize(); ok && y >= 0 && y < rows && x >= 0 && x <= cols {
		return true
	}
	return false
}

// GetSize returns the display’s current height and width.
func (d *Display) GetSize() (rows, cols int, ok bool) {
	if pty, _, ok := d.sess.Pty(); ok {
		return pty.Window.Height, pty.Window.Width, true
	}
	return 0, 0, false
}

// WriteLine writes a string to the screen and adds linebreak.
func (d *Display) WriteLine(str string) error {
	_, err := d.term.Write([]byte(str + "\n"))
	return err
}

// ReadLine reads a string from stdin.
func (d *Display) ReadLine(prompt string) (string, error) {
	d.term.SetPrompt(styles["reverse"] + styles["bright"] + prompt + " " + styles["yellow"] + ">" + styles["reset"] + " ")
	return d.term.ReadLine()
}

// CompleteWith sets autocompletion of line with the given slice of strings.
func (d *Display) CompleteWith(cmds []string) {
	d.term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// exit on empty line or not at end of line
		if line == "" || pos < len(line) {
			return
		}

		// only autocomplete on TAB
		if key != '\t' {
			return
		}

		found := false
		for _, cmd := range cmds {
			if matched, err := regexp.Match("^"+line+".*", []byte(cmd)); err == nil && matched {
				// if two commands at least match, return the common part
				if found {
					i := 0
					for i < len(newLine) && i < len(cmd) {
						if cmd[i] != newLine[i] {
							return cmd[:i], i, true
						}
						i++
					}
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

// WriteParsed parses a string and extracts commands from it and calls CompleteWith to set autocompletion to those commands.
// It is typically called with a context description.
//
//     Syntax      Command
//     **word**    word
//     *words*     look words
//     /words/cmd/ cmd
func (d *Display) WriteParsed(str string) (cmds []string, err error) {
	// **word** --> word
	reWord := regexp.MustCompile("\\*\\*([^\\s\\*]+)\\*\\*")
	str = reWord.ReplaceAllStringFunc(str, func(src string) string {
		word := reWord.ReplaceAllString(src, "$1")
		alreadyIn := false
		for _, cmd := range cmds {
			if cmd == word {
				alreadyIn = true
				break
			}
		}
		if !alreadyIn {
			cmds = append(cmds, word)
		}
		return styles["bright"] + word + styles["reset"]
	})

	// *words* --> look words
	reLook := regexp.MustCompile("\\*([^\\s\\*][^\\*]+[^\\s\\*])\\*")
	str = reLook.ReplaceAllStringFunc(str, func(src string) string {
		words := reLook.ReplaceAllString(src, "$1")
		alreadyIn := false
		for _, cmd := range cmds {
			if cmd == "look "+words {
				alreadyIn = true
				break
			}
		}
		if !alreadyIn {
			cmds = append(cmds, "look "+words)
		}
		return styles["italic"] + words + styles["reset"]
	})

	// /words/cmd/ --> cmd
	reCmd := regexp.MustCompile("/([^\\s/][^/]+[^\\s/])/([^\\s/][^/]+[^\\s/])/")
	str = reCmd.ReplaceAllStringFunc(str, func(src string) string {
		words := reCmd.ReplaceAllString(src, "$1")
		cmd := reCmd.ReplaceAllString(src, "$2")
		alreadyIn := false
		for _, cmdtest := range cmds {
			if cmdtest == cmd {
				alreadyIn = true
				break
			}
		}
		if !alreadyIn {
			cmds = append(cmds, cmd)
		}
		return styles["underscore"] + words + styles["reset"]
	})

	// write output to screen
	return cmds, d.WriteLine(str)
}
