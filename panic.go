package ego

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	stderrfile *os.File
)

func RedirectPanic(appname string) (err error) {
	{
		// 移除以往空白 panic 日志文件
		des, err := os.ReadDir("/log")
		if err != nil {
			log.Printf("failed to read /log[%s]", err)
			return err
		}

		for _, v := range des {
			if strings.HasPrefix(v.Name(), "panic_"+appname) {
				vinfo, err := v.Info()
				if err != nil {
					log.Printf("failed to get <%s> info[%s]", v.Name(), err)
					return err
				}
				if vinfo.Size() == 0 {
					err = os.Remove("/log/" + v.Name())
					if err != nil {
						log.Printf("failed to get %s info[%s]", v.Name(), err)
					} else {
						log.Printf("remove /log/%s[%+v]", v.Name(), vinfo)
					}
				}
			}
		}
	}

	path := fmt.Sprintf("/log/panic_%s_%d", appname, time.Now().Unix())
	stderrfile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("failed to open panic file<%s>", err)
		return err
	}

	err = syscall.Dup2(int(stderrfile.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Printf("failed to redirect stderr<%s>", err)
		return err
	}

	return err
}
