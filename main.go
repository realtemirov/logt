package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	header string = "┌────────────────────────────────────────────────────────────────────┐"
	bottom string = "│──────────┼─────────────────────────────────────────────────────────│"
	footer string = "└────────────────────────────────────────────────────────────────────┘"
)

type logType string

const (
	info    logType = "info"
	err     logType = "error"
	msg     logType = "message"
	succes  logType = "succes"
	debug   logType = "debug"
	warning logType = "warning"
	data    logType = "data"

	infoColor    color.Attribute = color.FgHiYellow
	errColor     color.Attribute = color.FgHiRed
	msgColor     color.Attribute = color.FgHiWhite
	succesColor  color.Attribute = color.FgHiGreen
	debugColor   color.Attribute = color.FgHiCyan
	warningColor color.Attribute = color.FgHiYellow
	dataColor    color.Attribute = color.FgHiMagenta
)

// LogData is a struct for creating a new log
type LogData struct {
	nameSpace string
}

// NewLogData create new log data
// namespace - name of namespace
// intro - show intro
//
// Example:
// logt.NewLogData("namespace", true)
//
// logt.NewLogData("namespace", false)
func NewLogData(namespace string, intro bool) *LogData {
	if intro {
		show()
	}
	return &LogData{
		nameSpace: namespace,
	}
}

// NewWriter create new writer
// We recommend for function name
// name - name of function
// Don't forget to close writer
//
// Example:
// w := logt.NewLogData("namespace", true).NewWriter("repository.Create()")
// defer w.Close()
// w.Info("some text")
// w.Error("some text")
func (l *LogData) NewWriter(name string) *Writer {
	w := &Writer{funcname: name, nameSpace: namespace(l.nameSpace)}
	w.funcName(true)
	return w
}

// Close close writer
func (w *Writer) Close() {
	w.funcName(false)
}

// Writer is a struct for writing logs
// Don't use this struct directly
// Don't forget to close writer
type Writer struct {
	funcname  string
	nameSpace string
}

// Data for warning message
// Color: magenta
func (w *Writer) Data(str ...any) {
	checker(str)
	print(w.nameSpace, data, str...)
}

// Debug for debugging message
// Color: cyan
func (w *Writer) Debug(str ...any) {
	checker(str)
	print(w.nameSpace, debug, str...)
}

// Error for error message
// Color: red
func (w *Writer) Error(str ...any) {
	checker(str)
	print(w.nameSpace, err, str...)
}

// Info for info message
// Color: yellow
func (w *Writer) Info(str ...any) {
	checker(str)
	print(w.nameSpace, info, str...)
}

// Msg for simple message
// Color: standard
func (w *Writer) Msg(str ...any) {
	checker(str)
	print(w.nameSpace, msg, str...)
}

// Succes for success message
// Color: green
func (w *Writer) Succes(str ...any) {
	checker(str)
	print(w.nameSpace, succes, str...)
}

// Warning for warning message
// Color: yellow
func (w *Writer) Warning(str ...any) {
	checker(str)
	print(w.nameSpace, warning, str...)
}

// Write for just output message
// Color: standard
func (w *Writer) Write(str ...any) {
	checker(str)
	print(w.nameSpace, msg, str...)
}

func print(n string, _type logType, str ...any) {
	var (
		c       color.Color
		x       string
		inf     string
		infBool bool = true
	)

	switch _type {
	case info:
		c = *color.New(infoColor)
		inf = replace(string(info), true)
	case err:
		c = *color.New(errColor)
		inf = replace(string(err), true)
	case msg:
		c = *color.New(msgColor)
		inf = replace(string(msg), true)
	case succes:
		c = *color.New(succesColor)
		inf = replace(string(succes), true)
	case debug:
		c = *color.New(debugColor)
		inf = replace(string(debug), true)
	case warning:
		c = *color.New(warningColor)
		inf = replace(string(warning), true)
	case data:
		c = *color.New(dataColor)
		inf = replace(string(data), true)
	}

	p(n, c.Sprint(header))
	defer p(n, c.Sprint(footer))
	for i := 0; i < len(str); i++ {

		s := strManual(str[i])

	n:
		s, x = checkStr(s)

		if s != "" {
			if i == 0 {
				p(n, c.Sprint(fmt.Sprintf("│%s│%s│", inf, x)))
				if infBool {
					inf = replace(inf, false)
					infBool = false
				}
			} else {
				if infBool {
					inf = replace(inf, false)
					infBool = false
				}
				p(n, c.Sprint(fmt.Sprintf("│%s│%s│", inf, x)))
			}
			goto n
		} else {
			if i != 0 {
				if infBool {
					inf = replace(inf, false)
					infBool = false
				}
			}
			p(n, c.Sprint(fmt.Sprintf("│%s│%s│", inf, x)))
			if infBool {
				inf = replace(inf, false)
				infBool = false
			}
		}
		if i != len(str)-1 {
			p(n, c.Sprint(bottom))
		}
	}
}

func replace(s string, first bool) string {
	if first {
		return fmt.Sprintf(" %-9s", s)
	}
	return strings.Repeat(" ", len(s))
}

func strManual(str any) string {
	js, err := json.MarshalIndent(str, "", "  ")
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%v", string(js))
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "    ")
	return s
}

func namespace(s string) string {
	if s == "" {
		return ""
	}
	return color.HiMagentaString("░ %s ░", s)
}

func checkStr(s string) (string, string) {
	var x string
	if len(s) > 55 {
		x = s[0:55]
		s = s[55:]
	} else {
		x = s
		s = ""
	}
	x = fmt.Sprintf(" %-55s ", x)

	return s, x
}

func p(namespace, data string) {
	year, month, day := time.Now().Date()
	hour, min, sec := time.Now().Clock()
	fmt.Printf("[%d-%02d-%02d %02d:%02d:%02d] %s %s\n", year, month, day, hour, min, sec, namespace, data)
}

func show() {
	log_t_blue1 := `
			╔═════════════════════════════════════════════════════════╗
			║                 `
	log_t_reset1 := `_                   _`
	log_t_blue2 := `	                  ║
			║                `
	log_t_reset2 := `| | ___   __ _      | |_`
	log_t_blue3 := `                 ║
			║                `
	log_t_reset3 := `| |/ _ \ / _' |_____| __|`
	log_t_blue4 := `                ║
			║                `
	log_t_reset4 := `| | (_) | (_| |_____| |_`
	log_t_blue5 := `                 ║
			║                `
	log_t_reset5 := `|_|\___/ \__, |      \__|`
	log_t_blue6 := `                ║
			║                         `
	log_t_reset6 := `|___/`
	log_t_blue7 := `                	       	  ║
			║                                                         ║
			╚═════════════════════════════════════════════════════════╝
	`
	desc_yellow_1 := `
				     ╔═════════════════════════════╗
				     ║        `
	desc_reset_2 := ` open source `
	desc_yellow_3 := `        ║
				     ╠═════════════════════════════╣
				     ║         `
	desc_reset_3 := `free logger`
	desc_yellow_4 := `         ║
				     ╚═════════════════════════════╝`

	blue := color.New(color.FgBlue, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	reset := color.New(color.FgHiWhite, color.Bold)
	blue.Print(log_t_blue1)
	reset.Print(log_t_reset1)
	blue.Print(log_t_blue2)
	reset.Print(log_t_reset2)
	blue.Print(log_t_blue3)
	reset.Print(log_t_reset3)
	blue.Print(log_t_blue4)
	reset.Print(log_t_reset4)
	blue.Print(log_t_blue5)
	reset.Print(log_t_reset5)
	blue.Print(log_t_blue6)
	reset.Print(log_t_reset6)
	blue.Print(log_t_blue7)
	reset.Println()
	yellow.Print(desc_yellow_1)
	reset.Print(desc_reset_2)
	yellow.Print(desc_yellow_3)
	reset.Print(desc_reset_3)
	yellow.Print(desc_yellow_4)
	reset.Println()
	reset.Println()
}

func (w *Writer) funcName(start bool) {
	var (
		n      string
		action string
	)
	if w.nameSpace != "" {
		n = w.nameSpace
	}
	if start {
		action = color.HiGreenString(":: START ::")
	} else {
		action = color.HiRedString("::  END  ::")
	}
	f_len := (50 - len(w.funcname)) / 2
	if f_len%2 != 0 {
		f_len++
	}

	txt := strings.Repeat(" ", f_len) + w.funcname + strings.Repeat(" ", f_len)
	p(n, fmt.Sprintf("%s %s", action, color.HiMagentaString(fmt.Sprintf("---> %s <---", txt))))
}

func checker(str []any) {
	if len(str) == 0 {
		panic("logt: empty message")
	}
}

func main() {
	// Create a new logger
	// nameSpace: The name of the logger
	// true: Enable show intro message
	l := NewLogData("repository", true)

	// Initialize writer logger
	w := l.NewWriter("item.Create()")
	defer w.Close()

	// Write a message
	w.Data("Hello world")
	w.Debug("Hello world")
	w.Error("Hello world")
	w.Info("Hello world")
	w.Msg("Hello world")
	w.Succes("Hello world")
	w.Warning("Hello world")
	w.Write("Hello world")

}
