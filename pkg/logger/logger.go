// // pkg/logger/logger.go
// package logger

// import (
//     "fmt"
//     "os"
//     "runtime"
//     "strings"

//     "github.com/sirupsen/logrus"
// )

// // Logger interface defines the logging methods
// type Logger interface {
//     Info(msg string, keysAndValues ...interface{})
//     Error(msg string, keysAndValues ...interface{})
//     Debug(msg string, keysAndValues ...interface{})
//     Warn(msg string, keysAndValues ...interface{})
//     Fatal(msg string, keysAndValues ...interface{})
//     WithField(key string, value interface{}) Logger
//     WithFields(fields map[string]interface{}) Logger
// }

// type logrusLogger struct {
//     log *logrus.Logger
//     entry *logrus.Entry
// }

// var defaultLogger *logrusLogger

// // Init initializes the logger with the given configuration
// func Init(level, format string) error {
//     log := logrus.New()

//     // Configure log output
//     log.SetOutput(os.Stdout)

//     // Set log level
//     logLevel, err := logrus.ParseLevel(level)
//     if err != nil {
//         return fmt.Errorf("invalid log level: %v", err)
//     }
//     log.SetLevel(logLevel)

//     // Set log format
//     if format == "json" {
//         log.SetFormatter(&logrus.JSONFormatter{
//             CallerPrettyfier: callerPrettyfier,
//         })
//     } else {
//         log.SetFormatter(&logrus.TextFormatter{
//             FullTimestamp:    true,
//             CallerPrettyfier: callerPrettyfier,
//         })
//     }

//     // Enable reporting of the calling function
//     log.SetReportCaller(true)

//     defaultLogger = &logrusLogger{
//         log:   log,
//         entry: log.WithFields(logrus.Fields{}),
//     }
//     return nil
// }

// // GetLogger returns the configured logger instance
// func GetLogger() Logger {
//     if defaultLogger == nil {
//         // Create a default logger if not initialized
//         log := logrus.New()
//         log.SetFormatter(&logrus.TextFormatter{
//             FullTimestamp: true,
//         })
//         log.SetLevel(logrus.InfoLevel)
        
//         defaultLogger = &logrusLogger{
//             log:   log,
//             entry: log.WithFields(logrus.Fields{}),
//         }
//     }
//     return defaultLogger
// }

// // NewLoggerFromLogrus creates a new Logger from a logrus Logger
// func NewLoggerFromLogrus(log *logrus.Logger) Logger {
//     return &logrusLogger{
//         log:   log,
//         entry: log.WithFields(logrus.Fields{}),
//     }
// }

// // Implement the Logger interface methods for logrusLogger

// func (l *logrusLogger) Info(msg string, keysAndValues ...interface{}) {
//     fields := createFields(keysAndValues...)
//     l.entry.WithFields(fields).Info(msg)
// }

// func (l *logrusLogger) Error(msg string, keysAndValues ...interface{}) {
//     fields := createFields(keysAndValues...)
//     l.entry.WithFields(fields).Error(msg)
// }

// func (l *logrusLogger) Debug(msg string, keysAndValues ...interface{}) {
//     fields := createFields(keysAndValues...)
//     l.entry.WithFields(fields).Debug(msg)
// }

// func (l *logrusLogger) Warn(msg string, keysAndValues ...interface{}) {
//     fields := createFields(keysAndValues...)
//     l.entry.WithFields(fields).Warn(msg)
// }

// func (l *logrusLogger) Fatal(msg string, keysAndValues ...interface{}) {
//     fields := createFields(keysAndValues...)
//     l.entry.WithFields(fields).Fatal(msg)
// }

// func (l *logrusLogger) WithField(key string, value interface{}) Logger {
//     return &logrusLogger{
//         log:   l.log,
//         entry: l.entry.WithField(key, value),
//     }
// }

// func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
//     return &logrusLogger{
//         log:   l.log,
//         entry: l.entry.WithFields(logrus.Fields(fields)),
//     }
// }

// // Helper function to create fields from key-value pairs
// func createFields(keysAndValues ...interface{}) logrus.Fields {
//     fields := logrus.Fields{}
//     for i := 0; i < len(keysAndValues); i += 2 {
//         if i+1 < len(keysAndValues) {
//             key, ok := keysAndValues[i].(string)
//             if !ok {
//                 continue
//             }
//             fields[key] = keysAndValues[i+1]
//         }
//     }
//     return fields
// }

// // callerPrettyfier customizes the caller information
// func callerPrettyfier(frame *runtime.Frame) (function string, file string) {
//     // Extract only the function name without the package path
//     function = frameToFunction(frame)
    
//     // Shorten file path
//     file = fmt.Sprintf("%s:%d", shortenFilePath(frame.File), frame.Line)
//     return function, file
// }

// // frameToFunction extracts just the function name
// func frameToFunction(frame *runtime.Frame) string {
//     // Split the full function name and return the last part
//     parts := strings.Split(frame.Function, ".")
//     return parts[len(parts)-1]
// }

// // shortenFilePath reduces the file path to a more readable format
// func shortenFilePath(path string) string {
//     // Keep only the last two parts of the path
//     parts := strings.Split(path, "/")
//     if len(parts) > 2 {
//         return strings.Join(parts[len(parts)-2:], "/")
//     }
//     return path
// }


package logger

import (
    "fmt"
    "os"
    "runtime"
    "strings"

    "github.com/sirupsen/logrus"
)

// Logger interface defines the logging methods
type Logger interface {
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Debug(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Fatal(msg string, keysAndValues ...interface{})
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
}

type logrusLogger struct {
    log   *logrus.Logger
    entry *logrus.Entry
}

var defaultLogger *logrusLogger

// Init initializes the logger with the given configuration
func Init(level, format string) error {
    log := logrus.New()
    log.SetOutput(os.Stdout)

    logLevel, err := logrus.ParseLevel(level)
    if err != nil {
        return fmt.Errorf("invalid log level: %v", err)
    }
    log.SetLevel(logLevel)

    if format == "json" {
        log.SetFormatter(&logrus.JSONFormatter{
            CallerPrettyfier: callerPrettyfier,
        })
    } else {
        log.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:    true,
            CallerPrettyfier: callerPrettyfier,
        })
    }

    log.SetReportCaller(true)

    defaultLogger = &logrusLogger{
        log:   log,
        entry: log.WithFields(logrus.Fields{}),
    }
    return nil
}

// GetLogger returns the configured logger instance
func GetLogger() Logger {
    if defaultLogger == nil {
        log := logrus.New()
        log.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
        log.SetLevel(logrus.InfoLevel)

        defaultLogger = &logrusLogger{
            log:   log,
            entry: log.WithFields(logrus.Fields{}),
        }
    }
    return defaultLogger
}

// NewLoggerFromLogrus creates a new Logger from a logrus Logger
func NewLoggerFromLogrus(log *logrus.Logger) Logger {
    return &logrusLogger{
        log:   log,
        entry: log.WithFields(logrus.Fields{}),
    }
}

// Implement the Logger interface methods for logrusLogger

func (l *logrusLogger) Info(msg string, keysAndValues ...interface{}) {
    fields := createFields(keysAndValues...)
    l.entry.WithFields(fields).Info(msg)
}

func (l *logrusLogger) Error(msg string, keysAndValues ...interface{}) {
    fields := createFields(keysAndValues...)
    l.entry.WithFields(fields).Error(msg)
}

func (l *logrusLogger) Debug(msg string, keysAndValues ...interface{}) {
    fields := createFields(keysAndValues...)
    l.entry.WithFields(fields).Debug(msg)
}

func (l *logrusLogger) Warn(msg string, keysAndValues ...interface{}) {
    fields := createFields(keysAndValues...)
    l.entry.WithFields(fields).Warn(msg)
}

func (l *logrusLogger) Fatal(msg string, keysAndValues ...interface{}) {
    fields := createFields(keysAndValues...)
    l.entry.WithFields(fields).Fatal(msg)
}

func (l *logrusLogger) WithField(key string, value interface{}) Logger {
    return &logrusLogger{
        log:   l.log,
        entry: l.entry.WithField(key, value),
    }
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
    return &logrusLogger{
        log:   l.log,
        entry: l.entry.WithFields(logrus.Fields(fields)),
    }
}

// Helper function to create fields from key-value pairs
func createFields(keysAndValues ...interface{}) logrus.Fields {
    fields := logrus.Fields{}
    for i := 0; i < len(keysAndValues); i += 2 {
        if i+1 < len(keysAndValues) {
            key, ok := keysAndValues[i].(string)
            if !ok {
                continue
            }
            fields[key] = keysAndValues[i+1]
        }
    }
    return fields
}

// callerPrettyfier customizes the caller information
func callerPrettyfier(frame *runtime.Frame) (function string, file string) {
    function = frameToFunction(frame)
    file = fmt.Sprintf("%s:%d", shortenFilePath(frame.File), frame.Line)
    return function, file
}

// frameToFunction extracts just the function name
func frameToFunction(frame *runtime.Frame) string {
    parts := strings.Split(frame.Function, ".")
    return parts[len(parts)-1]
}

// shortenFilePath reduces the file path to a more readable format
func shortenFilePath(path string) string {
    parts := strings.Split(path, "/")
    if len(parts) > 2 {
        return strings.Join(parts[len(parts)-2:], "/")
    }
    return path
}