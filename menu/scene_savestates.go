package menu

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/savestates"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneSavestates struct {
	entry
}

func buildSavestates() Scene {
	var list sceneSavestates

	tSavestates := l10n.T9(&i18n.Message{ID: "Savestates", Other: "Savestates"})

	list.label = tSavestates //"Savestates"

	tSaveState := l10n.T9(&i18n.Message{ID: "SaveState", Other: "Save State"})

	list.children = append(list.children, entry{
		label: tSaveState, //"Save State",
		icon:  "savestate",
		callbackOK: func() {
			name := utils.DatedName(state.GamePath)
			err := menu.TakeScreenshot(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			}
			err = savestates.Save(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				menu.stack[len(menu.stack)-1] = buildSavestates()
				menu.tweens.FastForward()
				txtI18n := l10n.T9(&i18n.Message{ID: "StateSaved", Other: "State saved."})
				ntf.DisplayAndLog(ntf.Success, "Menu", txtI18n)
			}
		},
	})

	gameName := utils.FileName(state.GamePath)
	gameName = strings.Replace(gameName, "[", "\\[", -1)
	gameName = strings.Replace(gameName, "]", "\\]", -1)
	paths, _ := filepath.Glob(settings.Current.SavestatesDirectory + "/" + gameName + "@*.state")
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, path := range paths {
		path := path
		date := strings.Replace(utils.FileName(path), gameName+"@", "", 1)
		list.children = append(list.children, entry{
			label: "Load " + date, // TODO: !Локализовать!
			icon:  "loadstate",
			path:  path,
			callbackOK: func() {
				err := savestates.Load(path)
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				} else {
					state.MenuActive = false

					txtI18n := l10n.T9(&i18n.Message{ID: "StateLoaded", Other: "State loaded."})
					ntf.DisplayAndLog(ntf.Success, "Menu", txtI18n)
				}
			},
			callbackX: func() { askDeleteSavestateConfirmation(func() { deleteSavestateEntry(&list, path) }) },
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneSavestates) Entry() *entry {
	return &s.entry
}

func (s *sceneSavestates) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneSavestates) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneSavestates) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneSavestates) update(dt float32) {
	genericInput(&s.entry, dt)
}

func removeSavestateEntry(s []entry, path string) []entry {
	l := []entry{}
	for _, e := range s {
		if e.path != path {
			l = append(l, e)
		}
	}

	return l
}

func deleteSavestateEntry(list *sceneSavestates, path string) {
	err := os.Remove(path)
	if err != nil {
		txtI18n := l10n.T9(&i18n.Message{ID: "CouldNotDelSavState", Other: "Could not delete savestate: %s"})
		ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n, err.Error())
		return
	}
	list.children = removeSavestateEntry(list.children, path)
	if list.ptr >= len(list.children) {
		list.ptr = len(list.children) - 1
	}
	genericAnimate(&list.entry)
}

// Override rendering
func (s *sceneSavestates) render() {
	list := &s.entry

	_, h := menu.GetFramebufferSize()

	thumbnailDrawCursor(list)

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		if e.labelAlpha > 0 {
			drawSavestateThumbnail(
				list, i,
				filepath.Join(settings.Current.ScreenshotsDirectory, utils.FileName(e.path)+".png"),
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, textColor.Alpha(e.iconAlpha),
			)
			menu.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				textColor.Alpha(e.iconAlpha))
			if i == 0 {
				menu.DrawImage(menu.icons["savestate"],
					680*menu.ratio-25*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-25*e.scale*menu.ratio+fontOffset,
					50*menu.ratio, 50*menu.ratio,
					e.scale, textColor.Alpha(e.iconAlpha))
			}

			menu.Font.SetColor(textColor.Alpha(e.labelAlpha))
			menu.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.5*menu.ratio, e.label)
		}
	}
}

func (s *sceneSavestates) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	ptr := menu.stack[len(menu.stack)-1].Entry().ptr

	_, upDown, _, a, b, x, _, _, _, guide := hintIcons()

	tHBarResume := l10n.T9(&i18n.Message{ID: "HBarResume", Other: "RESUME"})
	tHBarNavigate := l10n.T9(&i18n.Message{ID: "HBarNavigate", Other: "NAVIGATE"})
	tHBarBack := l10n.T9(&i18n.Message{ID: "HBarBack", Other: "BACK"})
	tHBarSave := l10n.T9(&i18n.Message{ID: "HBarSave", Other: "SAVE"})
	tHBarLoad := l10n.T9(&i18n.Message{ID: "HBarLoad", Other: "LOAD"})
	tHBarDelete := l10n.T9(&i18n.Message{ID: "HBarDelete", Other: "DELETE"})

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, tHBarResume, h)
	}
	stackHint(&stack, upDown, tHBarNavigate, h)
	stackHint(&stack, b, tHBarBack, h)
	if ptr == 0 {
		stackHint(&stack, a, tHBarSave, h)
	} else {
		stackHint(&stack, a, tHBarLoad, h)
	}

	list := menu.stack[len(menu.stack)-1].Entry()
	if list.children[list.ptr].callbackX != nil {
		stackHint(&stack, x, tHBarDelete, h)
	}
}
