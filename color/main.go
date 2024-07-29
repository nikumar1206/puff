package color

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Inverse   = "\033[7m"
	Hidden    = "\033[8m"
	Strikeout = "\033[9m"

	FgBlack         = "\033[30m"
	FgRed           = "\033[31m"
	FgGreen         = "\033[32m"
	FgYellow        = "\033[33m"
	FgBlue          = "\033[34m"
	FgMagenta       = "\033[35m"
	FgCyan          = "\033[36m"
	FgWhite         = "\033[37m"
	FgBrightBlack   = "\033[90m"
	FgBrightRed     = "\033[91m"
	FgBrightGreen   = "\033[92m"
	FgBrightYellow  = "\033[93m"
	FgBrightBlue    = "\033[94m"
	FgBrightMagenta = "\033[95m"
	FgBrightCyan    = "\033[96m"
	FgBrightWhite   = "\033[97m"

	BgBlack         = "\033[40m"
	BgRed           = "\033[41m"
	BgGreen         = "\033[42m"
	BgYellow        = "\033[43m"
	BgBlue          = "\033[44m"
	BgMagenta       = "\033[45m"
	BgCyan          = "\033[46m"
	BgWhite         = "\033[47m"
	BgBrightBlack   = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"
)

// Colorize applies the given ANSI attributes to the text.
func Colorize(text string, attributes ...string) string {
	var formatStr string
	for _, attr := range attributes {
		formatStr += attr
	}
	formatStr += text + Reset
	return formatStr
}

// ColorizeBold applies the Bold attribute along with any additional ANSI attributes to the text.
func ColorizeBold(text string, attributes ...string) string {
	return Colorize(text, append([]string{Bold}, attributes...)...)
}
