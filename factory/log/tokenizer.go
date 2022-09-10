package log

import "log"

type LogTokeniezr struct {
}

func (l *LogTokeniezr) Tokenize() {
	log.Println("Log Tokenizer Tokenize()")
}
