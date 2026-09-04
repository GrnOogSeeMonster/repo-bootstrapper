package logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var (
	debugMode bool
	logger    *slog.Logger
)

// Init initializes the logger with debug mode setting
func Init(debug bool) {
	debugMode = debug

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if debugMode {
		opts.Level = slog.LevelDebug
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}

// Info logs an info message
func Info(msg string, args ...any) {
	if logger == nil {
		Init(false)
	}
	logger.Info(msg, args...)
}

// Debug logs a debug message (only if debug mode enabled)
func Debug(msg string, args ...any) {
	if logger == nil {
		Init(false)
	}
	logger.Debug(msg, args...)
}

// Error logs an error message
func Error(msg string, err error, args ...any) {
	if logger == nil {
		Init(false)
	}
	allArgs := append([]any{"error", err}, args...)
	logger.Error(msg, allArgs...)
}

// OperationType represents different operation types
type OperationType int

const (
	OpCreateDir OperationType = iota
	OpCreateFile
	OpWriteFile
	OpGitInit
	OpGitAdd
	OpGitCommit
	OpGitRemoteAdd
	OpGitPush
	OpGHRepoCreate
	OpGHLabelCreate
	OpSpecKitInit
)

// Operation represents a single operation
type Operation struct {
	Type      OperationType
	Target    string
	Details   string
	DryRun    bool
	Timestamp time.Time
	Error     error
}

// OperationLog tracks all operations
type OperationLog struct {
	Operations []Operation
}

// Add adds an operation to the log
func (ol *OperationLog) Add(op Operation) {
	op.Timestamp = time.Time{}
	ol.Operations = append(ol.Operations, op)
}

// PrintSummary prints a summary of all operations
func (ol *OperationLog) PrintSummary(dryRun bool) {
	var created, modified, errors int

	for _, op := range ol.Operations {
		if op.Error != nil {
			errors++
		}
		switch op.Type {
		case OpCreateFile, OpCreateDir:
			created++
		case OpWriteFile:
			modified++
		}
	}

	prefix := ""
	if dryRun {
		prefix = "[DRY RUN] Would "
		fmt.Printf("\n%s create %d files/directories\n", prefix, created)
	} else {
		fmt.Printf("\nSummary: Created %d items", created)
		if modified > 0 {
			fmt.Printf(", modified %d items", modified)
		}
		if errors > 0 {
			fmt.Printf(", %d errors", errors)
		}
		fmt.Println()
	}
}

// ToJSON converts the log to JSON
func (ol *OperationLog) ToJSON() (string, error) {
	data, err := json.MarshalIndent(ol.Operations, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
