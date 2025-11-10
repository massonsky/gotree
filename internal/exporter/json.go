package exporter

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"time"

	_types "tree/internal/types"
)

type JSONExporter struct{}

// JSONEntry структура для сериализации
type JSONEntry struct {
	Path     string    `json:"path"`
	Type     string    `json:"type"` // "file" or "directory"
	Size     int64     `json:"size"` // 0 for directories
	Depth    int       `json:"depth"`
	ModTime  time.Time `json:"mod_time"`
	IsHidden bool      `json:"is_hidden"`
}

func (e *JSONExporter) Export(w io.Writer, entries []_types.Entry) error {
	jsonEntries := make([]JSONEntry, len(entries))

	for i, entry := range entries {
		jsonEntries[i] = JSONEntry{
			Path:     entry.Path,
			Type:     map[bool]string{true: "directory", false: "file"}[entry.Info.IsDir()],
			Size:     entry.Info.Size(),
			Depth:    entry.Depth,
			ModTime:  entry.Info.ModTime(),
			IsHidden: strings.HasPrefix(filepath.Base(entry.Path), "."),
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonEntries)
}
