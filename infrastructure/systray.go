package infrastructure

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"log"
)

func RunSystray() {
	systray.Run(systrayOnReady, nil)
}

func systrayOnReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("CleanRSS")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				log.Println("Quitting")
				systray.Quit()
			}
		}
	}()
}
