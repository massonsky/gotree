package tree

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tree/internal/config"
	"tree/internal/logger"
	"tree/internal/metrics"
	_type "tree/internal/types"
	"tree/internal/ui"

	"github.com/schollz/progressbar/v3"
	// Новый импорт
)

func WalkDirWithContext(
	ctx context.Context,
	root string,
	cfg *config.Config,
	progressEnabled bool,
) (WalkResult, error) {
	startTime := time.Now() // ← Теперь используется!
	logger.Debugf("Starting directory walk with progress: %t", progressEnabled)

	root, err := filepath.Abs(root)
	if err != nil {
		return WalkResult{}, err
	}

	// Считаем общее количество файлов для прогресс-бара
	totalFiles := 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if path != root {
			totalFiles++
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to count files: %v", err)
		return WalkResult{}, err
	}

	logger.Debugf("Total files to process: %d", totalFiles)

	var bar *progressbar.ProgressBar
	if progressEnabled {
		bar = ui.NewProgressBar(
			int64(totalFiles),
			"Scanning files",
			ui.DefaultProgressBarConfig(),
		)
		defer bar.Finish()
		ctx = ui.WithCancel(ctx, bar)
	}

	var entries []_type.Entry

	// Добавляем корневой элемент
	rootInfo, err := os.Stat(root)
	if err != nil {
		return WalkResult{}, err
	}
	entries = append(entries, _type.Entry{
		Path:  filepath.Base(root),
		Info:  rootInfo,
		Depth: 0,
	})

	// Основной обход
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			logger.Warn("Directory walk cancelled by user")
			return ctx.Err()
		default:
		}

		if path == root {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !cfg.ShowHiddenFiles && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		depth := len(strings.Split(relPath, string(filepath.Separator)))

		if cfg.MaxDepth > 0 && depth > cfg.MaxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		entries = append(entries, _type.Entry{
			Path:  relPath,
			Info:  info,
			Depth: depth,
		})

		if progressEnabled {
			bar.Add(1)
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Directory walk failed: %v", err)
		return WalkResult{}, err
	}

	// Собираем метрики
	mets := metrics.Collect(entries, startTime)
	logger.Infof("Found %d entries in %s", len(entries)-1, root)

	return WalkResult{
		Entries: entries,
		Metrics: mets,
	}, nil
}

// WalkDir совместимость (возвращает только entries)
func WalkDir(root string, cfg *config.Config) ([]_type.Entry, error) {
	result, err := WalkDirWithContext(context.Background(), root, cfg, true)
	if err != nil {
		return nil, err
	}
	return result.Entries, nil
}
