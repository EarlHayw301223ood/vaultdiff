package vault

import (
	"errors"
	"fmt"
)

// MergeStrategy controls how conflicting keys are resolved.
type MergeStrategy string

const (
	MergeStrategyOurs   MergeStrategy = "ours"
	MergeStrategyTheirs MergeStrategy = "theirs"
	MergeStrategyUnion  MergeStrategy = "union"
)

// MergeResult holds the merged data and metadata about the operation.
type MergeResult struct {
	Data      map[string]string
	Conflicts []string
	Added     []string
	Overwritten []string
}

// MergeOptions configures the merge behaviour.
type MergeOptions struct {
	Strategy MergeStrategy
	DryRun   bool
}

// Merge combines secret data from sourcePath into destPath using the given
// strategy. When DryRun is true the result is returned without writing.
func Merge(c LogicalClient, sourcePath, destPath string, opts MergeOptions) (*MergeResult, error) {
	if sourcePath == "" {
		return nil, errors.New("source path must not be empty")
	}
	if destPath == "" {
		return nil, errors.New("destination path must not be empty")
	}
	if sourcePath == destPath {
		return nil, errors.New("source and destination paths must differ")
	}

	src, err := FetchAtRef(c, sourcePath, "latest")
	if err != nil {
		return nil, fmt.Errorf("fetch source: %w", err)
	}
	dst, err := FetchAtRef(c, destPath, "latest")
	if err != nil {
		return nil, fmt.Errorf("fetch destination: %w", err)
	}

	result := &MergeResult{
		Data: make(map[string]string),
	}

	// Seed with destination data.
	for k, v := range dst {
		result.Data[k] = v
	}

	for k, sv := range src {
		dv, exists := dst[k]
		if !exists {
			result.Data[k] = sv
			result.Added = append(result.Added, k)
			continue
		}
		if sv == dv {
			continue
		}
		// Conflict: values differ.
		result.Conflicts = append(result.Conflicts, k)
		switch opts.Strategy {
		case MergeStrategyTheirs:
			result.Data[k] = sv
			result.Overwritten = append(result.Overwritten, k)
		case MergeStrategyUnion:
			result.Data[k] = dv + "," + sv
			result.Overwritten = append(result.Overwritten, k)
		default: // ours
			// keep destination value — no change needed
		}
	}

	if opts.DryRun {
		return result, nil
	}

	data := toInterfaceMap(result.Data)
	_, err = c.Write(kvV2WritePath(destPath), map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("write merged secret: %w", err)
	}
	return result, nil
}
