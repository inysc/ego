package ego

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/inysc/ego/config"
)

func RedirectPanic() error {

	path := fmt.Sprintf("%s_panic_%d", config.SrvName(), time.Now().Unix())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return err
	}

	return nil
}
