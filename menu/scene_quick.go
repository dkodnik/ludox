package menu

import (
	"github.com/libretro/ludo/favorites"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneQuick struct {
	entry
}

func buildQuickMenu() Scene {
	var list sceneQuick

	tQuickMenu := l10n.T9(&i18n.Message{ID: "QuickMenu", Other: "Quick Menu"})

	list.label = tQuickMenu //"Quick Menu"

	tResume := l10n.T9(&i18n.Message{ID: "Resume", Other: "Resume"})

	list.children = append(list.children, entry{
		label: tResume, //"Resume",
		icon:  "resume",
		callbackOK: func() {
			state.MenuActive = false
			state.FastForward = false
		},
	})

	tReset := l10n.T9(&i18n.Message{ID: "Reset", Other: "Reset"})

	list.children = append(list.children, entry{
		label: tReset, //"Reset",
		icon:  "reset",
		callbackOK: func() {
			state.Core.Reset()
			state.MenuActive = false
			state.FastForward = false
		},
	})

	tToFavorites := l10n.T9(&i18n.Message{ID: "ToFavorites", Other: "To Favorites"})

	list.children = append(list.children, entry{
		label: tToFavorites,
		icon:  "favorites-content",
		callbackOK: func() {
			favorites.Push(favorites.Game{
				Path:     state.GamePath,
				Name:     utils.FileName(state.GamePath),
				CorePath: state.CorePath,
			})
			txtI18n := l10n.T9(&i18n.Message{ID: "AddedToFavorites", Other: "Added to Favorites."})
			ntf.DisplayAndLog(ntf.Success, "Menu", txtI18n)
		},
	})

	tSavestates := l10n.T9(&i18n.Message{ID: "Savestates", Other: "Savestates"})

	list.children = append(list.children, entry{
		label: tSavestates, //"Savestates",
		icon:  "states",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildSavestates())
		},
	})

	tTakeScreenshot := l10n.T9(&i18n.Message{ID: "TakeScreenshot", Other: "Take Screenshot"})

	list.children = append(list.children, entry{
		label: tTakeScreenshot, //"Take Screenshot",
		icon:  "screenshot",
		callbackOK: func() {
			name := utils.DatedName(state.GamePath)
			err := menu.TakeScreenshot(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				txtI18n := l10n.T9(&i18n.Message{ID: "TookScreenshot", Other: "Took a screenshot."})
				ntf.DisplayAndLog(ntf.Success, "Menu", txtI18n)
			}
		},
	})

	tOptions := l10n.T9(&i18n.Message{ID: "Options", Other: "Options"})

	list.children = append(list.children, entry{
		label: tOptions, //"Options",
		icon:  "subsetting",
		callbackOK: func() {
			list.segueNext()
			menu.Push(buildCoreOptions())
		},
	})

	tDiskControl := l10n.T9(&i18n.Message{ID: "DiskControl", Other: "Disk Control"})

	if state.Core != nil && state.Core.DiskControlCallback != nil {
		list.children = append(list.children, entry{
			label: tDiskControl, //"Disk Control",
			icon:  "core-disk-options",
			callbackOK: func() {
				list.segueNext()
				menu.Push(buildCoreDiskControl())
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneQuick) Entry() *entry {
	return &s.entry
}

func (s *sceneQuick) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneQuick) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneQuick) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneQuick) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneQuick) render() {
	genericRender(&s.entry)
}

func (s *sceneQuick) drawHintBar() {
	genericDrawHintBar()
}
