package tree

import (
	"github.com/massonsky/tree/internal/metrics"
	_type "github.com/massonsky/tree/internal/types"
)

// WalkResult содержит результат обхода директории
type WalkResult struct {
	Entries []_type.Entry
	Metrics metrics.Metrics
}
