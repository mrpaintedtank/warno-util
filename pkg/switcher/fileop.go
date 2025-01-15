package switcher

import (
	"fmt"
	"os"
)

type fileOp struct {
	src, dst string
}

func performOperations(operations []fileOp) error {
	for _, op := range operations {
		fmt.Println("Renaming", op.src, "to", op.dst)
		if err := os.Rename(op.src, op.dst); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", op.src, op.dst, err)
		}
	}
	return nil
}
