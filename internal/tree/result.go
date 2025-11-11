package tree

import (
	"github.com/massonsky/gotree/internal/metrics"
	_type "github.com/massonsky/gotree/internal/types"
)

// WalkResult содержит результат обхода директории
type WalkResult struct {
	Entries []_type.Entry
	Metrics metrics.Metrics
}
