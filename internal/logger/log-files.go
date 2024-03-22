package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DEFAULT_LOG_FOLDER_NAME = "logs"
)

func WriteLogFile(contents []byte, workdir string, logFile string) error {
	path := filepath.Join(workdir, DEFAULT_LOG_FOLDER_NAME)

	if err := ensureLogDir(path); err != nil {
		return fmt.Errorf("unable to ensure the log directory exists: %w", err)
	}

	fpath := filepath.Join(path, logFile)
	if err := os.WriteFile(fpath, contents, 0644); err != nil {
		return fmt.Errorf("unable to write contents to log file: %w", err)
	}

	l := Get()
	l.Infof("wrote log contents to %s", fpath)

	return nil
}

func ensureLogDir(logdir string) error {
	if _, err := os.Stat(logdir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(logdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
