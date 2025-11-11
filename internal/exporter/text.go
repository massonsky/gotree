package exporter

import (
	"fmt"
	"io"
	"path/filepath"

	_types "github.com/massonsky/gotree/internal/types"
)

type TextExporter struct{}

func (e *TextExporter) Export(w io.Writer, entries []_types.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	// Находим максимальную глубину
	maxDepth := 0
	for _, entry := range entries {
		if entry.Depth > maxDepth {
			maxDepth = entry.Depth
		}
	}

	// Генерируем строки
	for i, entry := range entries {
		isLast := (i == len(entries)-1)
		line := formatTextEntry(entry, isLast, maxDepth)
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			return err
		}
	}

	return nil
}

func formatTextEntry(entry _types.Entry, isLast bool, maxDepth int) string {
	prefix := ""
	if entry.Depth > 0 {
		for d := 1; d < entry.Depth; d++ {
			prefix += "│   "
		}
		if isLast {
			prefix += "└── "
		} else {
			prefix += "├── "
		}
	}

	icon := "📄"
	if entry.Info.IsDir() {
		icon = "📁"
	}

	name := filepath.Base(entry.Path)
	if entry.Depth == 0 {
		name = entry.Path
	}

	line := fmt.Sprintf("%s%s %s", prefix, icon, name)

	if !entry.Info.IsDir() {
		size := formatSize(entry.Info.Size())
		line += fmt.Sprintf(" (%s)", size)
	}

	return line
}

func formatSize(bytes int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
	)
	switch {
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
