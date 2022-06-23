package log

import (
    "fmt"
    golog "log"
)

/*
    The main purpose of these methods is to log all GoStackInABox messages with
    the same prefix so that the logs are easy to identify and help users diagnose
    what is going on in their tests.
 */

/***********************
 *** Support versions ***
 ***********************/

const (
    logName = "GoStackInABox"
    logFmtString = "%s: %s"
)

type formattedFn func(string, ...interface{})
type standardFn func(...interface{})

var (
    doFatalf formattedFn = golog.Fatalf
    doFatal standardFn = golog.Fatal
    doFatalln standardFn = golog.Fatalln

    doPanicf formattedFn = golog.Panicf
    doPanic standardFn = golog.Panic
    doPanicln standardFn = golog.Panicln

    doPrintf formattedFn = golog.Printf
    doPrint standardFn = golog.Print
    doPrintln standardFn = golog.Println
)

func makeLogStringf(format string, v ...interface{}) string {
    coreLogString := fmt.Sprintf(format, v...)
    return fmt.Sprintf(logFmtString, logName, coreLogString)
}

func makeLogString(v ...interface{}) string {
    coreLogString := fmt.Sprint(v...)
    return fmt.Sprintf(logFmtString, logName, coreLogString)
}

func makeLogStringln(v ...interface{}) string {
    coreLogString := fmt.Sprintln(v...)
    return fmt.Sprintf(logFmtString, logName, coreLogString)
}

/***********************
 *** Printf versions ***
 ***********************/

func Fatalf(format string, v ...interface{}) {
    doFatalf(makeLogStringf(format, v...))
}

func Panicf(format string, v ...interface{}) {
    doPanicf(makeLogStringf(format, v...))
}

func Printf(format string, v ...interface{}) {
    doPrintf(makeLogStringf(format, v...))
}

/**********************
 *** Print versions ***
 **********************/

func Fatal(v ...interface{}) {
    doFatal(makeLogString(v...))
}

func Panic(v ...interface{}) {
    doPanic(makeLogString(v...))
}

func Print(v ...interface{}) {
    doPrint(makeLogString(v...))
}

/************************
 *** Println versions ***
 ************************/

func Fatalln(v ...interface{}) {
    doFatalln(makeLogStringln(v...))
}

func Panicln(v ...interface{}) {
    doPanicln(makeLogStringln(v...))
}

func Println(v ...interface{}) {
    doPrintln(makeLogStringln(v...))
}
