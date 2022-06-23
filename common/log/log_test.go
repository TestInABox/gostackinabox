package log
// not `log_test` so that we can access the internal methods for better testing

import (
    "fmt"
    "testing"

    //"github.com/TestInABox/gostackinabox/common/log"
)

func Test_MakeLogString(t *testing.T) {
    t.Run(
        "format",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            result := makeLogStringf("%s %s", "hello", "world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "standard",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            result := makeLogString("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "newline",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world\n", logName)
            result := makeLogStringln("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
}

func Test_FormattedLoggers(t *testing.T) {
    t.Run(
        "fatal",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doFatalf = func(formatted string, v ...interface{}) {
                result = fmt.Sprintf(formatted, v...)
            }
            Fatalf("%s %s", "hello", "world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "panic",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doPanicf = func(formatted string, v ...interface{}) {
                result = fmt.Sprintf(formatted, v...)
            }
            Panicf("%s %s", "hello", "world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "print",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doPrintf = func(formatted string, v ...interface{}) {
                result = fmt.Sprintf(formatted, v...)
            }
            Printf("%s %s", "hello", "world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
}

func Test_StandardLoggers(t *testing.T) {
    t.Run(
        "fatal",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doFatal = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Fatal("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "panic",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doPanic = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Panic("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "print",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world", logName)
            var result string
            doPrint = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Print("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
}

func Test_NewLineLoggers(t *testing.T) {
    t.Run(
        "fatal",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world\n", logName)
            var result string
            doFatalln = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Fatalln("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "panic",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world\n", logName)
            var result string
            doPanicln = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Panicln("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
    t.Run(
        "print",
        func(t *testing.T) {
            expectedMsg := fmt.Sprintf("%s: hello world\n", logName)
            var result string
            doPrintln = func(v ...interface{}) {
                result = fmt.Sprint(v...)
            }
            Println("hello world")
            if result != expectedMsg {
                t.Errorf("Log Messages do not match: \"%s\" != \"%s\"", result, expectedMsg)
            }
        },
    )
}
