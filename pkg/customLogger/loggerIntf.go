package customLogger

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	CallInfoStr() string
	PackageName() string
	FileName() string
	PackageAndFileNames() string
	CallInfo() *CallInfo
}

type CallInfo struct {
	PackageName string
	FileName    string
	FuncName    string
	Line        int
}
