package color

import "fmt"

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

// 字符串带中括号的
func BlackB(str string) string {
	return textColorB(textBlack, str)
}

func RedB(str string) string {
	return textColorB(textRed, str)
}
func YellowB(str string) string {
	return textColorB(textYellow, str)
}
func GreenB(str string) string {
	return textColorB(textGreen, str)
}
func CyanB(str string) string {
	return textColorB(textCyan, str)
}
func BlueB(str string) string {
	return textColorB(textBlue, str)
}
func PurpleB(str string) string {
	return textColorB(textPurple, str)
}
func WhiteB(str string) string {
	return textColorB(textWhite, str)
}
func textColorB(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm[%s]\x1b[0m", color, str)
}
func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
