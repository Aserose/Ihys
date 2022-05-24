package customLogger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

type logger struct {
	log *logrus.Logger
}

func NewLogger() Logger {
	return &logger{
		log: logrus.New(),
	}
}

func (l logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l logger) Print(args ...interface{}) {
	l.log.Print(args...)
}

func (l logger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

func (l logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l logger) CallInfoStr() string {
	result := l.CallInfo()

	return fmt.Sprintf("{ Package name: %s; File name: %s; Func name: %s; Line: %d } ",
		result.PackageName,
		result.FileName,
		result.FuncName,
		result.Line)
}

func (l logger) PackageName() string {
	result := l.CallInfo()

	return fmt.Sprintf(" {Package: %s }",
		result.PackageName)
}

func (l logger) FileName() string {
	result := l.CallInfo()

	return fmt.Sprintf("{ File: %s }",
		result.FileName)
}

func (l logger) PackageAndFileNames() string {
	result := l.CallInfo()

	return fmt.Sprintf("{ Package: %s; File: %s }",
		result.PackageName, result.FileName)
}

func (l logger) CallInfo() *CallInfo {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &CallInfo{
		PackageName: packageName,
		FileName:    fileName,
		FuncName:    funcName,
		Line:        line,
	}
}
