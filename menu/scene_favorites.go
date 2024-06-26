package menu

import (
	"os"
	"path/filepath"

	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/favorites"
	"github.com/libretro/ludo/history"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/state"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneFavorites struct {
	entry
}

func buildFavorites() Scene {
	var list sceneFavorites

	tFavorites := l10n.T9(&i18n.Message{ID: "Favorites", Other: "Favorites"})
	list.label = tFavorites

	favorites.Load()
	for _, game := range favorites.List {
		game := game // needed for callbackOK
		strippedName, tags := extractTags(game.Name)
		list.children = append(list.children, entry{
			label:      strippedName,
			subLabel:   game.System,
			gameName:   game.Name,
			path:       game.Path,
			system:     game.System,
			tags:       tags,
			callbackOK: func() { loadFavoritesEntry(&list, game) },
			callbackX:  func() { askDeleteGameConfirmation(func() { deleteFavoritesEntry(&list, game) }) },
		})
	}

	if len(favorites.List) == 0 {
		tEmptyFavorites := l10n.T9(&i18n.Message{ID: "EmptyFavorites", Other: "Empty favorites"})

		list.children = append(list.children, entry{
			label: tEmptyFavorites,
			icon:  "subsetting",
		})
	}

	list.segueMount()
	return &list
}

func loadFavoritesEntry(list Scene, game favorites.Game) {
	if _, err := os.Stat(game.Path); os.IsNotExist(err) {
		txtI18n := l10n.T9(&i18n.Message{ID: "GameNotFound", Other: "Game not found."})
		ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n)
		return
	}
	corePath := game.CorePath
	if _, err := os.Stat(corePath); os.IsNotExist(err) {
		txtI18n := l10n.T9(&i18n.Message{ID: "CoreNotFound", Other: "Core not found: %s"})
		ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n, filepath.Base(corePath))
		return
	}
	if state.CorePath != corePath {
		err := core.Load(corePath)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}

		state.SystemName = game.System
	}
	if state.GamePath != game.Path {
		err := core.LoadGame(game.Path)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			return
		}
		history.Push(history.Game{
			Path:     game.Path,
			Name:     game.Name,
			System:   game.System,
			CorePath: corePath,
		})
		list.segueNext()
		menu.Push(buildQuickMenu())
		menu.tweens.FastForward() // position the elements without animating
		state.MenuActive = false
	} else {
		list.segueNext()
		menu.Push(buildQuickMenu())
	}
}

func removeFavoritesGame(s []favorites.Game, game favorites.Game) []favorites.Game {
	l := []favorites.Game{}
	for _, g := range s {
		if g.Path != game.Path {
			l = append(l, g)
		}
	}
	return l
}

func removeFavoritesEntry(s []entry, game favorites.Game) []entry {
	l := []entry{}
	for _, g := range s {
		if g.path != game.Path {
			l = append(l, g)
		}
	}

	return l
}

func deleteFavoritesEntry(list *sceneFavorites, game favorites.Game) {
	favorites.List = removeFavoritesGame(favorites.List, game)
	favorites.Save()
	refreshTabs()
	list.children = removeFavoritesEntry(list.children, game)

	if len(favorites.List) == 0 {
		tEmptyFavorites := l10n.T9(&i18n.Message{ID: "EmptyFavorites", Other: "Empty favorites"})

		list.children = append(list.children, entry{
			label: tEmptyFavorites,
			icon:  "subsetting",
		})
	}

	if list.ptr >= len(list.children) {
		list.ptr = len(list.children) - 1
	}

	buildIndexes(&list.entry)
	genericAnimate(&list.entry)
}

// Generic stuff
func (s *sceneFavorites) Entry() *entry {
	return &s.entry
}

func (s *sceneFavorites) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneFavorites) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneFavorites) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneFavorites) update(dt float32) {
	genericInput(&s.entry, dt)
}

// Override rendering
func (s *sceneFavorites) render() {
	list := &s.entry

	_, h := menu.GetFramebufferSize()

	thumbnailDrawCursor(list)

	menu.ScissorStart(int32(510*menu.ratio), 0, int32(1310*menu.ratio), int32(h))

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			freeThumbnail(list, i)
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		if e.labelAlpha > 0 {
			drawThumbnail(
				list, i,
				e.system, e.gameName,
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, white.Alpha(e.iconAlpha),
			)
			menu.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				textColor.Alpha(e.iconAlpha))
			if e.path == state.GamePath && e.path != "" {
				menu.DrawCircle(
					680*menu.ratio,
					float32(h)*e.yp-14*menu.ratio+fontOffset,
					90*menu.ratio*e.scale,
					black.Alpha(e.iconAlpha))
				menu.DrawImage(menu.icons["resume"],
					680*menu.ratio-25*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-25*e.scale*menu.ratio+fontOffset,
					50*menu.ratio, 50*menu.ratio,
					e.scale, white.Alpha(e.iconAlpha))
			}

			// Offset on Y to vertically center label + sublabel if there is a sublabel
			slOffset := float32(0)
			if e.subLabel != "" {
				slOffset = 30 * menu.ratio * e.subLabelAlpha
			}

			menu.Font.SetColor(textColor.Alpha(e.labelAlpha))
			stack := 840 * menu.ratio
			menu.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset-slOffset,
				0.5*menu.ratio, e.label)
			stack += float32(int(menu.Font.Width(0.5*menu.ratio, e.label)))
			stack += 10

			for _, tag := range e.tags {
				if _, ok := menu.icons[tag]; ok {
					stack += 20
					menu.DrawImage(
						menu.icons[tag],
						stack, float32(h)*e.yp-22*menu.ratio-slOffset,
						60*menu.ratio, 44*menu.ratio, 1.0, white.Alpha(e.tagAlpha))
					menu.DrawBorder(stack, float32(h)*e.yp-22*menu.ratio-slOffset,
						60*menu.ratio, 44*menu.ratio, 0.05/menu.ratio, black.Alpha(e.tagAlpha/4))
					stack += 60 * menu.ratio
				}
			}

			menu.Font.SetColor(mediumGrey.Alpha(e.subLabelAlpha))
			menu.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset+60*menu.ratio-slOffset,
				0.5*menu.ratio, e.subLabel)
		}
	}

	menu.ScissorEnd()
}

func (s *sceneFavorites) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, _, a, b, x, _, _, _, guide := hintIcons()

	tHBarResume := l10n.T9(&i18n.Message{ID: "HBarResume", Other: "RESUME"})
	tHBarNavigate := l10n.T9(&i18n.Message{ID: "HBarNavigate", Other: "NAVIGATE"})
	tHBarBack := l10n.T9(&i18n.Message{ID: "HBarBack", Other: "BACK"})
	tHBarRun := l10n.T9(&i18n.Message{ID: "HBarRun", Other: "RUN"})
	tHBarDelete := l10n.T9(&i18n.Message{ID: "HBarDelete", Other: "DELETE"})

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, tHBarResume, h)
	}
	stackHint(&stack, upDown, tHBarNavigate, h)
	stackHint(&stack, b, tHBarBack, h)
	stackHint(&stack, a, tHBarRun, h)

	list := menu.stack[len(menu.stack)-1].Entry()
	if list.children[list.ptr].callbackX != nil {
		stackHint(&stack, x, tHBarDelete, h)
	}
}
