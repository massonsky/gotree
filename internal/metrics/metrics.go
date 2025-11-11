package metrics

import (
	"fmt"
	"time"

	_type "github.com/massonsky/gotree/internal/types"
)

// Metrics содержит статистику по директории
type Metrics struct {
	TotalFiles     int
	TotalDirs      int
	TotalSize      int64
	MaxDepth       int
	ScanDuration   time.Duration
	FilesPerSecond float64
}

// Collect собирает метрики из списка записей
func Collect(entries []_type.Entry, startTime time.Time) Metrics {
	var m Metrics
	m.ScanDuration = time.Since(startTime)

	for _, entry := range entries {
		if entry.Info.IsDir() {
			m.TotalDirs++
		} else {
			m.TotalFiles++
			m.TotalSize += entry.Info.Size()
		}

		if entry.Depth > m.MaxDepth {
			m.MaxDepth = entry.Depth
		}
	}

	// Вычисляем производительность
	if m.ScanDuration.Seconds() > 0 {
		m.FilesPerSecond = float64(m.TotalFiles+m.TotalDirs) / m.ScanDuration.Seconds()
	}

	return m
}

// String форматирует метрики для вывода
func (m Metrics) String() string {
	duration := m.ScanDuration.Truncate(time.Millisecond).String()
	perf := fmt.Sprintf("%.1f files/sec", m.FilesPerSecond)

	return fmt.Sprintf(`📊 Scan Metrics:
   Files:       %d
   Directories: %d
   Total Size:  %s
   Max Depth:   %d
   Duration:    %s
   Performance: %s`,
		m.TotalFiles,
		m.TotalDirs,
		FormatSize(m.TotalSize),
		m.MaxDepth,
		duration,
		perf,
	)
}

// formatSize преобразует байты в человекочитаемый формат
func FormatSize(bytes int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
