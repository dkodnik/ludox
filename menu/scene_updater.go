package menu

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/ludos"
	ntf "github.com/libretro/ludo/notifications"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneUpdater struct {
	entry
}

func buildUpdater() Scene {
	var list sceneUpdater

	tUpdaterMenu := l10n.T9(&i18n.Message{ID: "UpdaterMenu", Other: "Updater Menu"})

	list.label = tUpdaterMenu //"Updater Menu"

	tCheckingUpdates := l10n.T9(&i18n.Message{ID: "CheckingUpdates", Other: "Checking updates"})

	list.children = append(list.children, entry{
		label: tCheckingUpdates, //"Checking updates",
		icon:  "reload",
	})

	list.segueMount()

	go func() {
		rels, err := ludos.GetReleases()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}

		if len(*rels) > 0 {
			rel := (*rels)[0]

			if rel.Name[1:] == ludos.Version {
				tUp2Date := l10n.T9(&i18n.Message{ID: "Up2Date", Other: "Up to date"})

				list.children[0].label = tUp2Date //"Up to date"
				list.children[0].icon = "subsetting"
				return
			}

			tUpgrade2 := l10n.T9(&i18n.Message{ID: "Upgrade2", Other: "Upgrade to "})

			list.children[0].label = tUpgrade2 + rel.Name //"Upgrade to " + rel.Name
			list.children[0].icon = "menu_saving"
			list.children[0].callbackOK = func() {
				asset := ludos.FilterAssets(rel.Assets)
				if asset == nil {
					txtI18n := l10n.T9(&i18n.Message{ID: "NoMatchAsset", Other: "No matching asset"})
					ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n)
					return
				}
				go func() {
					path := filepath.Join(ludos.UpdatesDir, asset.Name)
					ludos.DownloadRelease(path, asset.BrowserDownloadURL)
				}()
			}
		} else {
			tNoUpdatesFound := l10n.T9(&i18n.Message{
				ID:    "NoUpdatesFound",
				Other: "No updates found",
			})

			list.children[0].label = tNoUpdatesFound //"No updates found"
			list.children[0].icon = "menu_exit"
		}
	}()

	return &list
}

func (s *sceneUpdater) Entry() *entry {
	return &s.entry
}

func (s *sceneUpdater) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneUpdater) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneUpdater) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneUpdater) update(dt float32) {
	if ludos.IsDownloading() {
		tLudosDownload := l10n.T9(&i18n.Message{ID: "LudosDownloadUpdate", Other: "Downloading update %.0f%%%%"})

		s.children[0].label = fmt.Sprintf(tLudosDownload, ludos.GetProgress()*100)
		s.children[0].icon = "reload"
		s.children[0].callbackOK = nil
	} else if ludos.IsDone() {
		tRebootAndUpgrade := l10n.T9(&i18n.Message{ID: "RebootAndUpgrade", Other: "Reboot and upgrade"})

		s.children[0].label = tRebootAndUpgrade //"Reboot and upgrade"
		s.children[0].icon = "reload"
		s.children[0].callbackOK = func() {
			cmd := exec.Command("/usr/sbin/shutdown", "-r", "now")
			core.UnloadGame()
			err := cmd.Run()
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			}
		}
	}

	genericInput(&s.entry, dt)
}

func (s *sceneUpdater) render() {
	genericRender(&s.entry)
}

func (s *sceneUpdater) drawHintBar() {
	genericDrawHintBar()
}
