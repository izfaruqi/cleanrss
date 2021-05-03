package infrastructure

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/pkg/browser"
	"log"
	"sync"
)

type SystrayInfo struct {
	mListen    *systray.MenuItem
	listenAddr string
	onQuit     func()
}

func RunSystray(info *SystrayInfo, cond *sync.Cond) {
	systray.Run(info.systrayOnReady(cond), nil)
}

func (i *SystrayInfo) SetListeningAddress(addr string) {
	i.listenAddr = addr
	i.mListen.SetTitle("Listening on " + addr)
}

func (i *SystrayInfo) SetOnQuitClicked(onQuit func()) {
	i.onQuit = onQuit
}

func (i *SystrayInfo) systrayOnReady(cond *sync.Cond) func() {
	return func() {
		cond.L.Lock()
		defer cond.Broadcast()
		defer cond.L.Unlock()
		systray.SetIcon(icon.Data)
		systray.SetTitle("CleanRSS")
		i.mListen = systray.AddMenuItem("Listening on ", "")
		mQuit := systray.AddMenuItem("Quit", "Quit the app")
		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					log.Println("Quitting")
					if i.onQuit != nil {
						i.onQuit()
					}
					systray.Quit()
				case <-i.mListen.ClickedCh:
					err := browser.OpenURL("http://" + i.listenAddr)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}()
	}
}
