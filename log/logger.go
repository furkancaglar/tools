package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Log_type uint8

var file__error, file__info, file__log io.Writer

const opts = os.O_CREATE | os.O_WRONLY | os.O_APPEND
const perm = 0644
const (
	ERR Log_type = iota
	INFO
	LOG
)

var _files map[Log_type]io.Writer

func init() {
	var err error
	file__error, err = os.OpenFile("error.log", opts, perm)
	__panic__err(err)
	file__info, err = os.OpenFile("info.log", opts, perm)
	__panic__err(err)
	file__log, err = os.OpenFile("log.log", opts, perm)
	__panic__err(err)
	_files = map[Log_type]io.Writer{
		ERR:  file__error,
		INFO: file__info,
		LOG:  file__log,
	}
}

func log__format(data ...interface{}) []byte {
	var _data = []interface{}{time.Now().Format("2006-01-02 03:04:05"), ">>\t"}
	return []byte(fmt.Sprint(append(_data, data...)...))
}

func Log_file_by_name(fname string, data ...interface{}) {
	f, err := os.OpenFile(fname, opts, perm)
	if nil != err {
		fmt.Fprintln(os.Stderr, log__format(err.Error()))
		return
	}
	defer f.Close()
	Log_file(f, data)
}

func Log_file(f io.Writer, data ...interface{}) {
	f.Write(log__format(data))
	f.Write([]byte{'\n'})
}

func Log(t Log_type, data ...interface{}) {
	f, ok := _files[t]
	if !ok {
		fmt.Fprintln(os.Stderr, log__format("Unknown logger type"))
		return
	}
	Log_file(f, data)
}

func __panic__err(err error) {
	if nil != err {
		panic(err)
	}
}
