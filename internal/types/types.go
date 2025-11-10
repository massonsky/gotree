package types

import (
	"io"
	"os"
)

// Entry представляет элемент файловой системы
type Entry struct {
	Path  string
	Info  os.FileInfo
	Depth int // Глубина вложенности для форматирования вывода
}

// Exporter интерфейс для всех форматов экспорта
type Exporter interface {
	Export(w io.Writer, entries []interface{}) error
}

// Format поддерживаемые форматы
type Format string

const (
	FormatPNG  Format = "png"
	FormatTXT  Format = "txt"
	FormatJSON Format = "json"
)

type TextExporter struct{}

type PNGExporter struct {
	fontPath string
}
