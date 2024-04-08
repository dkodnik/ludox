package menu

import (
	"github.com/libretro/ludo/ludos"
	ntf "github.com/libretro/ludo/notifications"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type sceneWiFi struct {
	entry
}

func buildWiFi() Scene {
	var list sceneWiFi

	tWiFiMenu := l10n.T9(&i18n.Message{ID: "WiFiMenu", Other: "WiFi Menu"})

	list.label = tWiFiMenu //"WiFi Menu"

	tLooking4Networks := l10n.T9(&i18n.Message{ID: "Looking4Networks", Other: "Looking for networks"})

	list.children = append(list.children, entry{
		label: tLooking4Networks, //"Looking for networks",
		icon:  "reload",
	})

	list.segueMount()

	go func() {
		networks, err := ludos.ScanNetworks()
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		}

		if len(networks) > 0 {
			list.children = []entry{}
			for _, network := range networks {
				network := network
				list.children = append(list.children, entry{
					label:       network.SSID,
					icon:        "menu_network",
					stringValue: func() string { return ludos.NetworkStatus(network) },
					callbackOK: func() {
						list.segueNext()
						menu.Push(buildKeyboard(
							"Passphrase for "+network.SSID,
							func(pass string) {
								go func() {
									if err := ludos.ConnectNetwork(network, pass); err != nil {
										ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
									}
								}()
							},
						))
					},
				})
				list.segueMount()
				menu.tweens.FastForward()
			}
		} else {
			tNoNetworkFound := l10n.T9(&i18n.Message{ID: "NoNetworkFound", Other: "No network found"})

			list.children[0].label = tNoNetworkFound //"No network found"
			list.children[0].icon = "close"
		}
	}()

	return &list
}

func (s *sceneWiFi) Entry() *entry {
	return &s.entry
}

func (s *sceneWiFi) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneWiFi) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneWiFi) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneWiFi) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneWiFi) render() {
	genericRender(&s.entry)
}

func (s *sceneWiFi) drawHintBar() {
	w, h := menu.GetFramebufferSize()
	menu.DrawRect(0, float32(h)-70*menu.ratio, float32(w), 70*menu.ratio, 0, lightGrey)

	_, upDown, _, a, b, _, _, _, _, _ := hintIcons()

	tHBarNavigate := l10n.T9(&i18n.Message{ID: "HBarNavigate", Other: "NAVIGATE"})
	tHBarBack := l10n.T9(&i18n.Message{ID: "HBarBack", Other: "BACK"})
	tHBarConnect := l10n.T9(&i18n.Message{ID: "HBarConnect", Other: "CONNECT"})

	var stack float32
	stackHint(&stack, upDown, tHBarNavigate, h)
	stackHint(&stack, b, tHBarBack, h)
	if s.children[0].callbackOK != nil {
		stackHint(&stack, a, tHBarConnect, h)
	}
}
