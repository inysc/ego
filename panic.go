package ego

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func RedirectPanic(appname string) {
	path := fmt.Sprintf("/log/%s_panic_%d", appname, time.Now().Unix())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open panic file<%s>", err)
	}

	err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to redirect stderr<%s>", err)
	}

	des, err := os.ReadDir("/log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read /log[%s]", err)
	}

	// 移除以往空白 panic 日志文件
	panicPrefix := appname + "_panic_"
	for _, v := range des {
		if strings.HasPrefix(v.Name(), panicPrefix) && path != v.Name() {
			vinfo, err := v.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get <%s> info[%s]", v.Name(), err)
			}
			if vinfo.Size() == 0 {
				err = os.Remove("/log/" + v.Name())
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to get %s info[%s]", v.Name(), err)
				} else {
					fmt.Printf("remove /log/%s[%+v]\n", v.Name(), vinfo)
				}
			}
		}
	}
}
