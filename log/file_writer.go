package log

import (
	"os"
	"sync"
)

type RotatePolicy string

const (
	ROTATE_POLICY_SIZE RotatePolicy = "SIZE"
	ROTATE_POLICY_HOUR RotatePolicy = "HOUR"
	ROTATE_POLICY_DAY  RotatePolicy = "DAY"
)

type logfileOptions struct {
	rotatePolicy RotatePolicy
	rotateSizeMB int64
}
type FileOption func(*logfileOptions)

func applyFileOptions(loggerOptions *logfileOptions, options ...FileOption) {
	for _, option := range options {
		option(loggerOptions)
	}
}

func defaultFileOptions() *logfileOptions {
	return &logfileOptions{rotatePolicy: ROTATE_POLICY_DAY}
}

func WithRotatePolicy(policy RotatePolicy) FileOption {
	return func(options *logfileOptions) {
		options.rotatePolicy = policy
	}
}
func WithRotateSizeMB(sizeMB int64) FileOption {
	return func(options *logfileOptions) {
		options.rotateSizeMB = sizeMB
	}
}

type FileWriter struct {
	logFile *os.File
	options *logfileOptions
	sync.Mutex
}

func NewFileWriter(logPath string, options ...FileOption) (*FileWriter, error) {
	logFile, err := newFile(logPath)
	if nil != err {
		return nil, err
	}
	opts := defaultFileOptions()
	applyFileOptions(opts, options...)
	writer := &FileWriter{logFile: logFile, options: opts}

	writer.StartRotate()

	return writer, nil
}

func (writer *FileWriter) Write(p []byte) (n int, err error) {
	return writer.logFile.Write(p)
}
