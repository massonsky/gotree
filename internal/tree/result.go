package tree

import (
	"tree/internal/metrics"
	_type "tree/internal/types"
)

// WalkResult содержит результат обхода директории
type WalkResult struct {
	Entries []_type.Entry
	Metrics metrics.Metrics
}
