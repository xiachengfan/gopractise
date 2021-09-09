package xlog

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type rollingFile struct {
	mu       sync.Mutex
	file     *os.File
	basePath string
	filePath string
	fileFrag string

	rolling RollingFormat
}

type RollingFormat string

const (
	MonthlyRolling  RollingFormat = "2006-01"
	DailyRolling                  = "2006-01-02"
	HourlyRolling                 = "2006-01-02-15"
	MinutelyRolling               = "2006-01-02-15-04"
	SecondlyRolling               = "2006-01-02-15-04-05"
)

func (r *rollingFile) roll() error {
	suffix := time.Now().Format(string(r.rolling))
	if r.file != nil {
		if suffix == r.fileFrag {
			return nil
		}
		r.file.Close()
		r.file = nil
	}
	r.fileFrag = suffix
	r.filePath = fmt.Sprintf("%s.%s", r.basePath, r.fileFrag)

	if dir, _ := filepath.Split(r.basePath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return errors.Wrap(err, " MkdirAll ERROR")
		}
	}

	f, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return errors.Wrap(err, "openFile error")
	} else {
		r.file = f
		return nil
	}
}

func (r *rollingFile) Write(b []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.roll(); err != nil {
		return 0, err
	}

	n, err := r.file.Write(b)
	if err != nil {
		return n, errors.Wrap(err, "Write error")
	} else {
		return n, nil
	}
}

func NewRollingFile(basePath string, rolling RollingFormat) (io.Writer, error) {
	if _, file := filepath.Split(basePath); file == "" {
		return nil, errors.Errorf("invalid base-path = %s, file name is required", basePath)
	}
	return &rollingFile{basePath: basePath, rolling: rolling}, nil
}
