package menu

import (
	"fmt"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneCoreDiskControl struct {
	entry
}

func buildCoreDiskControl() Scene {
	var list sceneCoreDiskControl

	tCoreDiskControl := l10n.T9(&i18n.Message{ID: "CoreDiskControl", Other: "Core Disk Control"})

	list.label = tCoreDiskControl //"Core Disk Control"

	for i := uint(0); i < state.Core.DiskControlCallback.GetNumImages(); i++ {
		index := i
		list.children = append(list.children, entry{
			label: fmt.Sprintf("Disk %d", index+1),
			icon:  "subsetting",
			stringValue: func() string {
				if index == state.Core.DiskControlCallback.GetImageIndex() {
					return "Active"
				}
				return ""
			},
			callbackOK: func() {
				if index == state.Core.DiskControlCallback.GetImageIndex() {
					return
				}
				state.Core.DiskControlCallback.SetEjectState(true)
				state.Core.DiskControlCallback.SetImageIndex(index)
				state.Core.DiskControlCallback.SetEjectState(false)

				txtI18n := l10n.T9(&i18n.Message{ID: "Switched2Disk", Other: "Switched to disk %d."})
				ntf.DisplayAndLog(ntf.Success, "Menu", txtI18n, index+1)
				state.MenuActive = false
			},
		})
	}

	tNoDisk := l10n.T9(&i18n.Message{ID: "NoDisk", Other: "No disk"})

	if len(list.children) == 0 {
		list.children = append(list.children, entry{
			label: tNoDisk, //"No disk",
			icon:  "subsetting",
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneCoreDiskControl) Entry() *entry {
	return &s.entry
}

func (s *sceneCoreDiskControl) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneCoreDiskControl) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneCoreDiskControl) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneCoreDiskControl) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneCoreDiskControl) render() {
	genericRender(&s.entry)
}

func (s *sceneCoreDiskControl) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, leftRight, _, b, _, _, _, _, guide := hintIcons()

	tHBarResume := l10n.T9(&i18n.Message{ID: "HBarResume", Other: "RESUME"})
	tHBarNavigate := l10n.T9(&i18n.Message{ID: "HBarNavigate", Other: "NAVIGATE"})
	tHBarBack := l10n.T9(&i18n.Message{ID: "HBarBack", Other: "BACK"})
	tHBarSet := l10n.T9(&i18n.Message{ID: "HBarSet", Other: "SET"})

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, tHBarResume, h)
	}
	stackHint(&stack, upDown, tHBarNavigate, h)
	stackHint(&stack, b, tHBarBack, h)
	stackHint(&stack, leftRight, tHBarSet, h)
}
