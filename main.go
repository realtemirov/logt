package logt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	header string = "┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────"
	bottom string = "│──────────┼──────────────────────────────────────────────────────────────────────────────────────────────────────"
	footer string = "└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────"
)

type logType string

type ctxKey string

const (
	info    logType = "info"
	err     logType = "error"
	msg     logType = "message"
	succes  logType = "succes"
	debug   logType = "debug"
	warning logType = "warning"
	data    logType = "data"
	withCtx logType = "context"
	bot     logType = "bot"

	infoColor    color.Attribute = color.FgHiYellow
	errColor     color.Attribute = color.FgHiRed
	msgColor     color.Attribute = color.FgHiWhite
	succesColor  color.Attribute = color.FgHiGreen
	debugColor   color.Attribute = color.FgHiCyan
	warningColor color.Attribute = color.FgHiYellow
	dataColor    color.Attribute = color.FgHiMagenta
	botColor     color.Attribute = color.FgHiBlue
	contextKey   ctxKey          = "logt-key"
)

// Types ---------------------------------------------------------------------------------

type ILog interface {
	NewWriter(functionName string, saveFile bool) IWriter
	SetContext(ctx context.Context, fields ...any) context.Context
}

type IWriter interface {
	Close()
	Data(str ...any)
	Debug(str ...any)
	Error(str ...any)
	Info(str ...any)
	Succes(str ...any)
	Warning(str ...any)
	Msg(str ...any)
	Write(str ...any)
	Send(str ...any)
	FromContext(ctx context.Context) IWriter
}

// Log is a struct for creating a new log implementation
type Log struct {
	Token     string
	UserID    int64
	Logo      bool
	NameSpace string
}

// Writer is a struct for writing logs
type writer struct {
	funcname  string
	nameSpace string
	userID    int64
	bot       *tgbotapi.BotAPI
	txt       strings.Builder
	save      bool
	ctx       strings.Builder
}

type detail struct {
	namespace string
	userID    int64
	token     string
	bot       *tgbotapi.BotAPI
}

// Public functions and methods  ---------------------------------------------------------------------------------

// NewLog create new log data
// namespace - name of namespace
// intro - show intro
//
// Example:
//
// logt.NewLog(log *Log)
//
//	type Log struct {
//	 // Telegram bot token
//		Token     string = "123456789:ABCDEFGHi-ijklmnop-QRSTUVXYZ"
//	 // Telegram user ID
//		UserID     int64 = 9876543210
//		Logo      bool
//		NameSpace string
//	}
//
// logt.NewLog("namespace", false)
func NewLog(log *Log) ILog {

	if log.Logo {
		show()
	}

	if log.Token != "" && log.UserID != 0 {
		user, tgBot, error := checkBot(log.Token, log.UserID)
		if error != nil {
			print(namespace(log.NameSpace), err, false, error)
			panic(error)
		} else {
			print(namespace(log.NameSpace), debug, false, user)
		}
		return &detail{
			namespace: log.NameSpace,
			userID:    log.UserID,
			token:     log.Token,
			bot:       tgBot,
		}
	}

	return &detail{
		namespace: log.NameSpace,
		userID:    log.UserID,
		token:     log.Token,
		bot:       nil,
	}
}

// NewWriter create new writer
//
// Set true if you want to save function logs to a file.
//
// Attention!
//
// This application may have a slight impact on performance and memory.
//
// I recommend that you keep only the essential functions.
//
// Example:
// w := l.NewWriter("repository.Create()",true)
//
// defer w.Close()
//
// w.Info("some text")
//
// w.Error("some text")
//
// If you set true, it save to a file, "2006-01-02 15:04:05-repository.Create().txt"
func (l *detail) NewWriter(functionName string, saveFile bool) IWriter {
	w := &writer{funcname: functionName, nameSpace: namespace(l.namespace), userID: l.userID, bot: l.bot, save: saveFile}
	w.funcName(true)
	return w
}

// Set value to context and return context
func (l *detail) SetContext(ctx context.Context, fields ...any) context.Context {

	str := strings.Builder{}
	length := len(fields)

	value, ok := ctx.Value(contextKey).(string)
	if ok {
		if length > 0 {
			str.WriteString(value + ", ")
		} else {
			str.WriteString(value)
		}
	}
	for i := 0; i < length; i++ {
		str.WriteString(strManual(fields[i]))
		if i != length-1 {
			str.WriteString(", ")
		}
	}
	return context.WithValue(ctx, contextKey, str.String())
}

// Close close writer
func (w *writer) Close() {
	if w.save {
		file, err := os.Create(fmt.Sprintf("%s-%s.txt", time.Now().Format("2006-01-02 15:04:05"), w.funcname))
		if err != nil {
			w.Error(err)
		}
		defer file.Close()

		file.WriteString("// Generated by logt\n")
		file.WriteString("// If you print this text, it will be print with color\n")
		file.WriteString(w.txt.String())
	}
	w.txt.Reset()
	w.funcName(false)
}

// Data for warning message
//
// Color: magenta
func (w *writer) Data(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, data, w.save, str...)[1])
	} else {
		print(w.nameSpace, data, w.save, str...)
	}
}

// Debug for debugging message
//
// Color: cyan
func (w *writer) Debug(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, debug, w.save, str...)[1])
	} else {
		print(w.nameSpace, debug, w.save, str...)
	}
}

// Error for error message
//
// Color: red
func (w *writer) Error(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	errors.Join()
	if w.save {
		w.txt.WriteString(print(w.nameSpace, err, w.save, str...)[1])
	} else {
		print(w.nameSpace, err, w.save, str...)
	}
}

// Info for info message
//
// Color: yellow
func (w *writer) Info(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, info, w.save, str...)[1])
	} else {
		print(w.nameSpace, info, w.save, str...)
	}
}

// Msg for simple message
//
// Color: standard
func (w *writer) Msg(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, msg, w.save, str...)[1])
	} else {
		print(w.nameSpace, msg, w.save, str...)
	}
}

// Succes for success message
//
// Color: green
func (w *writer) Succes(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, succes, w.save, str...)[1])
	} else {
		print(w.nameSpace, succes, w.save, str...)
	}
}

// Warning for warning message
//
// Color: yellow
func (w *writer) Warning(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, warning, w.save, str...)[1])
	} else {
		print(w.nameSpace, warning, w.save, str...)
	}
}

// Write for just output message
//
// Color: standard
func (w *writer) Write(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	if w.save {
		w.txt.WriteString(print(w.nameSpace, msg, w.save, str...)[1])
	} else {
		print(w.nameSpace, msg, w.save, str...)
	}
}

// Send log to telegram with userID
//
// Integration: telegram
func (w *writer) Send(str ...any) {
	if w.ctx.String() != "" {
		str = append(str, w.ctx.String())
	}
	checker(str)
	txt := print(w.nameSpace, bot, w.save, str...)
	w.txt.WriteString(txt[1])
	for len(txt[0]) > 4090 {
		_, error := w.bot.Send(tgbotapi.NewMessage(w.userID, txt[0][0:4090]))
		txt[0] = txt[0][4090:]
		if error != nil {
			w.Error(error.Error())
		}
	}
	_, error := w.bot.Send(tgbotapi.NewMessage(w.userID, txt[0]))
	if error != nil {
		w.Error(error.Error())
	}
}

// Get value from context
func (w *writer) FromContext(ctx context.Context) IWriter {

	value, ok := ctx.Value(contextKey).(string)
	if ok {
		if value != "" {
			value = string(contextKey) + "-" + value
		} else {
			value = string(contextKey) + "-" + "\"value\":\"not found\""
		}
	} else {
		value = string(contextKey) + "-" + "\"value\":\"not found\""
	}

	w.ctx.WriteString(value)
	return w
}

// Priveta functions and methods ---------------------------------------------------------------------------------
func print(n string, _type logType, save bool, str ...any) []string {
	var (
		c       color.Color
		x       string
		inf     string
		send    string
		txt     strings.Builder
		file    strings.Builder
		need    bool = false
		infBool bool = true
		isError bool = false
		hasCtx  bool = false
	)

	switch _type {
	case info:
		c = *color.New(infoColor)
		inf = replace(string(info), true)
	case err:
		c = *color.New(errColor)
		inf = replace(string(err), true)
		isError = true
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
	case bot:
		c = *color.New(botColor)
		inf = replace(string(bot), true)
		need = true
	}

	file.WriteString(p(n, c.Sprint(header)))

	for i := 0; i < len(str); i++ {
		if isError {
			if fmt.Sprintf("%T", str[i]) == "*errors.errorString" {
				str[i] = str[i].(error).Error()
			}
		}

		var s string

		value, ok := str[i].(string)
		if ok {
			if strings.HasPrefix(value, string(contextKey)) {
				hasCtx = true
				s = value[len(string(contextKey)+"-"):]
			}
		}
		if !hasCtx {
			s = strManual(str[i])
		}
	n:
		s, x = checkStr(s)

		if s != "" {
			if i == 0 {

				if !infBool {
					if hasCtx {
						inf = replace(string(withCtx), true)
						hasCtx = false
					}
				}

				send = c.Sprint(fmt.Sprintf("│%s│%s", inf, x))
				inf = replace(inf, false)
				if need {
					txt.WriteString(x)
				}
				if save {
					file.WriteString(p(n, send))
				} else {
					fmt.Println(1)
					p(n, send)
				}
				if infBool {
					inf = replace(inf, false)
					infBool = false
					hasCtx = false
				}
			} else {
				if infBool {
					inf = replace(inf, false)
					infBool = false
					hasCtx = false
				}
				if !infBool {
					if hasCtx {
						inf = replace(string(withCtx), true)
						hasCtx = false
					}
				}
				send = c.Sprint(fmt.Sprintf("│%s│%s", inf, x))
				inf = replace(inf, false)
				if need {
					txt.WriteString(x)
				}
				if save {
					file.WriteString(p(n, send))
				} else {
					p(n, send)
				}
			}
			goto n
		} else {
			if i != 0 {
				if infBool {
					inf = replace(inf, false)
					infBool = false
					hasCtx = false
				}
			}

			if !infBool {
				if hasCtx {
					inf = replace(string(withCtx), true)
					hasCtx = false
				}
			}
			send = c.Sprint(fmt.Sprintf("│%s│%s", inf, x))
			inf = replace(inf, false)
			if need {
				txt.WriteString(x)
			}
			if save {
				file.WriteString(p(n, send))
			} else {
				p(n, send)
			}
			if infBool {
				inf = replace(inf, false)
				infBool = false
				hasCtx = false
			}
		}
		if i != len(str)-1 {
			send = c.Sprint(bottom)
			if need {
				txt.WriteString(x)
			}
			if save {
				file.WriteString(p(n, send))
			} else {
				p(n, send)
			}
		}
	}
	file.WriteString(p(n, c.Sprint(footer)))
	return []string{txt.String(), file.String()}
}

func replace(s string, first bool) string {
	if first {
		return fmt.Sprintf(" %-9s", s)
	}
	return strings.Repeat(" ", len(s))
}

func strManual(str any) string {
	js, err := json.MarshalIndent(str, "", "    ")
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%v", string(js))
	return s
}

func namespace(s string) string {
	if s == "" {
		return ""
	}
	return color.HiWhiteString("| %s |", s)
}

func checkStr(s string) (string, string) {
	var x string

	if len(s) > 75 {
		x = s[0:75]
		if strings.Contains(x, "\n") {
			index := strings.Index(x, "\n")
			x = x[0:index]
			s = s[index+1:]
		} else {
			s = s[75:]
		}
	} else {

		if strings.Contains(s, "\n") {
			index := strings.Index(s, "\n")
			x = s[0:index]
			s = s[index+1:]
		} else {
			x = s
			s = ""
		}
	}
	x = fmt.Sprintf(" %-75s ", x)
	return s, x
}

func p(namespace, data string) string {
	year, month, day := time.Now().Date()
	hour, min, sec := time.Now().Clock()
	s := fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] %s %s\n", year, month, day, hour, min, sec, namespace, data)
	fmt.Print(s)
	return s
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

func (w *writer) funcName(start bool) {
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
	f_len := (90 - len(w.funcname)) / 2
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

func getBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return bot
}

func checkBot(token string, userID int64) (*tgbotapi.User, *tgbotapi.BotAPI, error) {
	tgBot := getBot(token)
	user, err := tgBot.GetMe()
	if err != nil {
		return nil, nil, err
	}
	_, err = tgBot.Send(tgbotapi.NewMessage(userID, "👋 Hello from logt!"))
	if err != nil {
		return nil, nil, err
	}
	return &user, tgBot, nil
}

type User struct {
	Firstname string  `json:"first_name"`
	Lastname  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Address   Address `json:"address"`
}

type Address struct {
	Country string `json:"country"`
	City    string `json:"city"`
}
