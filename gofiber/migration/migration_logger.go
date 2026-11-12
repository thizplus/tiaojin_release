package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MigrationLogger handles logging for migration operations
type MigrationLogger struct {
	file      *os.File
	mu        sync.Mutex
	startTime time.Time
	logPath   string
}

var migrationLog *MigrationLogger

// InitMigrationLogger initializes the migration logger
func InitMigrationLogger() error {
	// Create logs directory if not exists
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("migration_%s.log", timestamp))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	migrationLog = &MigrationLogger{
		file:      file,
		startTime: time.Now(),
		logPath:   logPath,
	}

	// Write header
	migrationLog.writeHeader()

	fmt.Printf("📝 Log file: %s\n", logPath)
	return nil
}

// CloseMigrationLogger closes the log file
func CloseMigrationLogger() {
	if migrationLog != nil && migrationLog.file != nil {
		migrationLog.writeSummary()
		migrationLog.file.Close()
	}
}

func (l *MigrationLogger) writeHeader() {
	l.mu.Lock()
	defer l.mu.Unlock()

	header := fmt.Sprintf(`
================================================================================
  MIGRATION LOG
================================================================================
  Started:  %s
  Host:     %s
================================================================================

`, l.startTime.Format("2006-01-02 15:04:05 MST"), getHostname())

	l.file.WriteString(header)
}

func (l *MigrationLogger) writeSummary() {
	l.mu.Lock()
	defer l.mu.Unlock()

	duration := time.Since(l.startTime)
	summary := fmt.Sprintf(`

================================================================================
  MIGRATION COMPLETED
================================================================================
  Finished: %s
  Duration: %v
================================================================================
`, time.Now().Format("2006-01-02 15:04:05 MST"), duration.Round(time.Second))

	l.file.WriteString(summary)
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	logMessage("INFO", format, args...)
}

// LogSuccess logs a success message
func LogSuccess(format string, args ...interface{}) {
	logMessage("SUCCESS", format, args...)
}

// LogWarning logs a warning message
func LogWarning(format string, args ...interface{}) {
	logMessage("WARNING", format, args...)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	logMessage("ERROR", format, args...)
}

// LogSection logs a section header
func LogSection(title string) {
	if migrationLog == nil {
		return
	}

	migrationLog.mu.Lock()
	defer migrationLog.mu.Unlock()

	section := fmt.Sprintf("\n--- %s ---\n", title)
	migrationLog.file.WriteString(fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), section))
}

// LogStats logs statistics
func LogStats(label string, count int) {
	LogInfo("%s: %d", label, count)
}

func logMessage(level, format string, args ...interface{}) {
	if migrationLog == nil {
		return
	}

	migrationLog.mu.Lock()
	defer migrationLog.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	migrationLog.file.WriteString(line)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// LogSyncStart logs the start of a sync operation
func LogSyncStart(operation string, since time.Time) {
	LogSection(fmt.Sprintf("START: %s", operation))
	LogInfo("Syncing data since: %s", since.Format("2006-01-02 15:04:05"))
}

// LogSyncEnd logs the end of a sync operation with stats
func LogSyncEnd(operation string, stats map[string]int) {
	LogSection(fmt.Sprintf("END: %s", operation))
	for key, value := range stats {
		LogStats(key, value)
	}
}

// LogFixStart logs the start of a fix operation
func LogFixStart(operation string) {
	LogSection(fmt.Sprintf("FIX: %s", operation))
}

// LogFixEnd logs the end of a fix operation
func LogFixEnd(operation string, updated, skipped, errors int) {
	LogSuccess("%s completed - Updated: %d, Skipped: %d, Errors: %d", operation, updated, skipped, errors)
}
