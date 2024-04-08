package menu

import (
	"fmt"
	"os"

	//"os/user"
	"sort"

	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/libretro"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneTabs struct {
	entry
}

func buildTabs() Scene {
	var list sceneTabs
	list.label = "Ludo"

	tMainMenuTab := l10n.T9(&i18n.Message{ID: "MainMenuTab", Other: "Main Menu"})
	tMainMenuSub := l10n.T9(&i18n.Message{ID: "MainMenuSub", Other: "Load cores and games manually"})

	list.children = append(list.children, entry{
		label:    tMainMenuTab, //"Main Menu",
		subLabel: tMainMenuSub, //"Load cores and games manually",
		icon:     "main",
		callbackOK: func() {
			menu.Push(buildMainMenu())
		},
	})

	tSettingsTab := l10n.T9(&i18n.Message{ID: "SettingsTab", Other: "Settings"})
	tSettingsSub := l10n.T9(&i18n.Message{ID: "SettingsSub", Other: "Configure Ludo"})

	list.children = append(list.children, entry{
		label:    tSettingsTab, //"Settings",
		subLabel: tSettingsSub, //"Configure Ludo",
		icon:     "setting",
		callbackOK: func() {
			menu.Push(buildSettings())
		},
	})

	tHistoryTab := l10n.T9(&i18n.Message{ID: "HistoryTab", Other: "History"})
	tHistorySub := l10n.T9(&i18n.Message{ID: "HistorySub", Other: "Play again"})

	list.children = append(list.children, entry{
		label:    tHistoryTab, //"History",
		subLabel: tHistorySub, //"Play again",
		icon:     "history",
		callbackOK: func() {
			menu.Push(buildHistory())
		},
	})

	list.children = append(list.children, getPlaylists()...)

	tAddGamesTab := l10n.T9(&i18n.Message{ID: "AddGamesTab", Other: "Add games"})
	tAddGamesSub := l10n.T9(&i18n.Message{ID: "AddGamesSub", Other: "Scan your collection"})

	tScanDir := l10n.T9(&i18n.Message{ID: "ScanDir", Other: "<Scan this directory>"})

	list.children = append(list.children, entry{
		label:    tAddGamesTab, //"Add games",
		subLabel: tAddGamesSub, //"Scan your collection",
		icon:     "add",
		callbackOK: func() {
			//usr, _ := user.Current()
			menu.Push(buildExplorer(settings.Current.FileDirectory, nil,
				func(path string) {
					scanner.ScanDir(path, refreshTabs)
				},
				&entry{
					label: tScanDir, //"<Scan this directory>",
					icon:  "scan",
				},
				nil,
			))
		},
	})

	list.segueMount()

	return &list
}

// refreshTabs is called after playlist scanning is complete. It inserts the new
// playlists in the tabs, and makes sure that all the icons are positioned and
// sized properly.
func refreshTabs() {
	e := menu.stack[0].Entry()
	l := len(e.children)
	pls := getPlaylists()

	// This assumes that the 3 first tabs are not playlists, and that the last
	// tab is the scanner.
	e.children = append(e.children[:3], append(pls, e.children[l-1:]...)...)

	// Update which tab is the active tab after the refresh
	if e.ptr >= 3 {
		e.ptr += len(pls) - (l - 4)
	}

	// Ensure new icons are styled properly
	for i := range e.children {
		if i == e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.75
			e.children[i].width = 500
		} else if i < e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.25
			e.children[i].width = 128
		} else if i > e.ptr {
			e.children[i].iconAlpha = 1
			e.children[i].scale = 0.25
			e.children[i].width = 128
		}
	}

	// Adapt the tabs scroll value
	if len(menu.stack) == 1 {
		menu.scroll = float32(e.ptr * 128)
	} else {
		e.children[e.ptr].margin = 1360
		menu.scroll = float32(e.ptr*128 + 680)
	}
}

// getPlaylists browse the filesystem for CSV files, parse them and returns
// a list of menu entries. It is used in the tabs, but could be used somewhere
// else too.
func getPlaylists() []entry {
	playlists.Load()

	// To store the keys in slice in sorted order
	var keys []string
	for k := range playlists.Playlists {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	txtI18nGames := l10n.T9(&i18n.Message{ID: "nGames", Other: "%d Games"})

	var pls []entry
	for _, path := range keys {
		path := path
		filename := utils.FileName(path)
		count := playlists.Count(path)
		label := playlists.ShortName(filename)
		pls = append(pls, entry{
			label:    label,
			subLabel: fmt.Sprintf(txtI18nGames, count),
			icon:     filename,
			callbackOK: func() {
				menu.Push(buildPlaylist(path))
			},
			callbackX: func() { askDeletePlaylistConfirmation(func() { deletePlaylist(path) }) },
		})
	}
	return pls
}

func deletePlaylist(path string) {
	err := os.Remove(path)
	if err != nil {
		txtI18n := l10n.T9(&i18n.Message{ID: "CouldNotDelPlaylist", Other: "Could not delete playlist: %s"})
		ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n, err.Error())
		return
	}
	menu.stack[0].Entry().ptr++
	delete(playlists.Playlists, path)
	refreshTabs()
}

func (tabs *sceneTabs) Entry() *entry {
	return &tabs.entry
}

func (tabs *sceneTabs) segueMount() {
	for i := range tabs.children {
		e := &tabs.children[i]

		if i == tabs.ptr {
			e.labelAlpha = 1
			e.iconAlpha = 1
			e.scale = 0.75
			e.width = 500
		} else if i < tabs.ptr {
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		} else if i > tabs.ptr {
			e.labelAlpha = 0
			e.iconAlpha = 1
			e.scale = 0.25
			e.width = 128
		}
	}

	tabs.animate()
}

func (tabs *sceneTabs) segueBack() {
	tabs.animate()
}

func (tabs *sceneTabs) animate() {
	for i := range tabs.children {
		e := &tabs.children[i]

		var labelAlpha, scale, width float32
		if i == tabs.ptr {
			labelAlpha = 1
			scale = 0.75
			width = 500
		} else if i < tabs.ptr {
			labelAlpha = 0
			scale = 0.25
			width = 128
		} else if i > tabs.ptr {
			labelAlpha = 0
			scale = 0.25
			width = 128
		}

		menu.tweens[&e.labelAlpha] = gween.New(e.labelAlpha, labelAlpha, 0.15, ease.OutSine)
		menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 1, 0.15, ease.OutSine)
		menu.tweens[&e.scale] = gween.New(e.scale, scale, 0.15, ease.OutSine)
		menu.tweens[&e.width] = gween.New(e.width, width, 0.15, ease.OutSine)
		menu.tweens[&e.margin] = gween.New(e.margin, 0, 0.15, ease.OutSine)
	}
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, float32(tabs.ptr*128), 0.15, ease.OutSine)
}

func (tabs *sceneTabs) segueNext() {
	cur := &tabs.children[tabs.ptr]
	menu.tweens[&cur.margin] = gween.New(cur.margin, 1360, 0.15, ease.OutSine)
	menu.tweens[&menu.scroll] = gween.New(menu.scroll, menu.scroll+680, 0.15, ease.OutSine)
	for i := range tabs.children {
		e := &tabs.children[i]
		if i != tabs.ptr {
			menu.tweens[&e.iconAlpha] = gween.New(e.iconAlpha, 0, 0.15, ease.OutSine)
		}
	}
}

func (tabs *sceneTabs) update(dt float32) {
	// Right
	repeatRight(dt, input.NewState[0][libretro.DeviceIDJoypadRight] == 1, func() {
		tabs.ptr++
		if tabs.ptr >= len(tabs.children) {
			tabs.ptr = 0
		}
		audio.PlayEffect(audio.Effects["down"])
		tabs.animate()
	})

	// Left
	repeatLeft(dt, input.NewState[0][libretro.DeviceIDJoypadLeft] == 1, func() {
		tabs.ptr--
		if tabs.ptr < 0 {
			tabs.ptr = len(tabs.children) - 1
		}
		audio.PlayEffect(audio.Effects["up"])
		tabs.animate()
	})

	// OK
	if input.Released[0][libretro.DeviceIDJoypadA] == 1 {
		if tabs.children[tabs.ptr].callbackOK != nil {
			audio.PlayEffect(audio.Effects["ok"])
			tabs.segueNext()
			tabs.children[tabs.ptr].callbackOK()
		}
	}

	// X
	if input.Released[0][libretro.DeviceIDJoypadX] == 1 {
		if tabs.children[tabs.ptr].callbackX != nil {
			tabs.children[tabs.ptr].callbackX()
		}
	}
}

func (tabs sceneTabs) render() {
	_, h := menu.GetFramebufferSize()

	stackWidth := 710 * menu.ratio
	for i, e := range tabs.children {

		cf := colorful.Hcl(float64(i)*20, 0.5, 0.5)
		c := video.Color{R: float32(cf.R), G: float32(cf.B), B: float32(cf.G), A: e.iconAlpha}

		x := -menu.scroll*menu.ratio + stackWidth + e.width/2*menu.ratio

		stackWidth += e.width*menu.ratio + e.margin*menu.ratio

		if e.labelAlpha > 0 {
			menu.Font.SetColor(c.Alpha(e.labelAlpha))
			lw := menu.Font.Width(0.5*menu.ratio, e.label)
			menu.Font.Printf(x-lw/2, float32(int(float32(h)/2+250*menu.ratio)), 0.5*menu.ratio, e.label)
			lw = menu.Font.Width(0.4*menu.ratio, e.subLabel)
			menu.Font.Printf(x-lw/2, float32(int(float32(h)/2+330*menu.ratio)), 0.4*menu.ratio, e.subLabel)
		}

		menu.DrawImage(menu.icons["hexagon"],
			x-220*e.scale*menu.ratio, float32(h)/2-220*e.scale*menu.ratio,
			440*menu.ratio, 440*menu.ratio, e.scale, c)

		menu.DrawImage(menu.icons[e.icon],
			x-128*e.scale*menu.ratio, float32(h)/2-128*e.scale*menu.ratio,
			256*menu.ratio, 256*menu.ratio, e.scale, white.Alpha(e.iconAlpha))
	}
}

func (tabs sceneTabs) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, _, leftRight, a, _, x, _, _, _, guide := hintIcons()

	tHBarResume := l10n.T9(&i18n.Message{ID: "HBarResume", Other: "RESUME"})
	tHBarNavigate := l10n.T9(&i18n.Message{ID: "HBarNavigate", Other: "NAVIGATE"})
	tHBarOpen := l10n.T9(&i18n.Message{ID: "HBarOpen", Other: "OPEN"})
	tHBarDelete := l10n.T9(&i18n.Message{ID: "HBarDelete", Other: "DELETE"})

	var stack float32
	if state.CoreRunning {
		stackHint(&stack, guide, tHBarResume, h)
	}
	stackHint(&stack, leftRight, tHBarNavigate, h)
	stackHint(&stack, a, tHBarOpen, h)

	list := menu.stack[0].Entry()
	if list.children[list.ptr].callbackX != nil {
		stackHint(&stack, x, tHBarDelete, h)
	}
}
