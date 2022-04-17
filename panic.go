package ego

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/inysc/ego/config"
)

func RedirectPanic() error {

	path := fmt.Sprintf("/log/%s_panic_%d", config.SrvName(), time.Now().Unix())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return err
	}

	des, err := os.ReadDir("/log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read /log[%s]", err)
	}
	for _, v := range des {
		if strings.HasPrefix(v.Name(), config.SrvName()) && !strings.HasSuffix(path, v.Name()) {
			vinfo, err := v.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get %s info[%s]", v.Name(), err)
				return nil
			}
			if vinfo.Size() == 0 {
				os.Remove("/log/" + v.Name())
				fmt.Printf("remove /log/%s[%+v]\n", v.Name(), vinfo)
			}
		}
	}

	return nil
}
