package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func IsFileExist(filePath string) error {
	if len(strings.TrimSpace(filePath)) == 0 {
		return os.ErrNotExist
	}

	p, _ := filepath.Abs(filePath)

	if _, err := os.Stat(p); err != nil {
		return fmt.Errorf(`file "%s" err %v`, filePath, err)
	}

	return nil
}
