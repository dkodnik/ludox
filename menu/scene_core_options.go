package menu

import (
	"strings"

	"github.com/libretro/ludo/core"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneCoreOptions struct {
	entry
}

func buildCoreOptions() Scene {
	var list sceneCoreOptions

	tCoreOptions := l10n.T9(&i18n.Message{ID: "CoreOptions", Other: "Core Options"})

	list.label = tCoreOptions //"Core Options"

	if core.Options == nil {
		tNoOptions := l10n.T9(&i18n.Message{ID: "NoOptions", Other: "No options"})

		list.children = append(list.children, entry{
			label: tNoOptions, //"No options",
			icon:  "subsetting",
		})
		list.segueMount()
		return &list
	}

	for _, v := range core.Options.Vars {
		v := v
		list.children = append(list.children, entry{
			label: strings.Replace(v.Desc, "%", "%%", -1),
			icon:  "subsetting",
			stringValue: func() string {
				val := v.Choices[v.Choice]
				return strings.Replace(val, "%", "%%", -1)
			},
			incr: func(direction int) {
				v.Choice += direction
				if v.Choice < 0 {
					v.Choice = len(v.Choices) - 1
				} else if v.Choice > len(v.Choices)-1 {
					v.Choice = 0
				}
				core.Options.Updated = true
				err := core.Options.Save()
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Core", "Error saving core options: %v", err.Error())
				}
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *sceneCoreOptions) Entry() *entry {
	return &s.entry
}

func (s *sceneCoreOptions) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneCoreOptions) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneCoreOptions) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneCoreOptions) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneCoreOptions) render() {
	genericRender(&s.entry)
}

func (s *sceneCoreOptions) drawHintBar() {
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
