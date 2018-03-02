/*
Package display handles all UI operations in sta.
*/
package display

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"regexp"
	"strings"
)

// Display is a type for handling i/o in sta.
type Display struct {
	sess            ssh.Session
	term            *terminal.Terminal
	completeOptions []string
}

// New returns a pointer to a *Display.
func New(sess ssh.Session) *Display {
	d := &Display{
		sess: sess,
		term: terminal.NewTerminal(sess, ""),
	}

	// autocompletion set with d.completeOptions
	d.term.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// exit on empty line or not at end of line
		if line == "" || pos < len(line) {
			return
		}

		// only autocomplete on TAB
		if key != '\t' {
			return
		}

		// first find all possible options
		var matching []string
		var minLen = 1 << 16
		for _, cmd := range d.completeOptions {
			if matched, err := regexp.Match("^"+line+".*", []byte(cmd)); err == nil && matched {
				matching = append(matching, cmd)
				if len(cmd) < minLen {
					minLen = len(cmd)
				}
			}
		}

		// if nothing was found, return
		if len(matching) == 0 {
			return
		}

		// if something was found, look for the common part of all matches
		newLine = line
		newPos = pos
		ok = true
		// for each character in matching commands
	charLoop:
		for ; newPos < minLen; newPos++ {
			// if they are not all equal, stop here
			c := matching[0][newPos]
			for _, m := range matching {
				if m[newPos] != c {
					break charLoop
				}
			}
			// if they are all equal, append the character to newLine
			newLine += string(c)
		}

		// print matching options to screen and return
		if len(matching) > 1 {
			d.term.Write([]byte("\n" + escape["clearline"] + strings.Join(matching, ", ") + escape["up"] + escape["bol"]))
		} else /* if !d.isCursorAtBottomOfScreen() */ {
			d.term.Write([]byte("\n" + escape["clearline"] + escape["up"] + escape["bol"]))
		}
		return
	}
	return d
}

// AppendComplete appends cmds to completeOptions field of d.
// The function checks for duplicates before appending.
func (d *Display) AppendComplete(cmds []string) {
cmdsLoop:
	for _, cmd := range cmds {
		for _, option := range d.completeOptions {
			if option == cmd {
				continue cmdsLoop
			}
		}
		d.completeOptions = append(d.completeOptions, cmd)
	}
}

// ResetComplete resets completeOptions field of d.
func (d *Display) ResetComplete() {
	d.completeOptions = nil
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
	d.term.SetPrompt(styles["reverse"] + styles["bright"] + prompt + styles["yellow"] + ">" + styles["reset"] + " ")
	return d.term.ReadLine()
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
