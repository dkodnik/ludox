package settings

import (
	"os/user"
	"path/filepath"

	"github.com/adrg/xdg"
)

func defaultSettings() Settings {
	usr, _ := user.Current()
	return Settings{
		VideoFullscreen:   false,
		VideoMonitorIndex: 0,
		VideoFilter:       "Pixel Perfect",
		MapAxisToDPad:     false,
		AudioVolume:       0.5,
		MenuAudioVolume:   0.25,
		ShowHiddenFiles:   false,
		CoreForPlaylist: map[string]string{
			"Atari - 2600":                                   "stella2014_libretro",
			"Atari - 5200":                                   "atari800_libretro",
			"Atari - 7800":                                   "prosystem_libretro",
			"Atari - Jaguar":                                 "virtualjaguar_libretro",
			"Atari - Lynx":                                   "handy_libretro",
			"Atari - ST":                                     "hatari_libretro",
			"Bandai - WonderSwan Color":                      "mednafen_wswan_libretro",
			"Bandai - WonderSwan":                            "mednafen_wswan_libretro",
			"Cave Story":                                     "nxengine_libretro",
			"ChaiLove":                                       "chailove_libretro",
			"Coleco - ColecoVision":                          "bluemsx_libretro",
			"FBNeo - Arcade Games":                           "fbneo_libretro",
			"GCE - Vectrex":                                  "vecx_libretro",
			"Magnavox - Odyssey2":                            "o2em_libretro",
			"Microsoft - MSX":                                "bluemsx_libretro",
			"Microsoft - MSX2":                               "bluemsx_libretro",
			"NEC - PC Engine SuperGrafx":                     "mednafen_pce_libretro",
			"NEC - PC Engine - TurboGrafx 16":                "mednafen_pce_libretro",
			"NEC - PC Engine CD - TurboGrafx-CD":             "mednafen_pce_libretro",
			"NEC - PC-FX":                                    "mednafen_pcfx_libretro",
			"Nintendo - Family Computer Disk System":         "fceumm_libretro",
			"Nintendo - Game Boy Advance":                    "mgba_libretro",
			"Nintendo - Game Boy Color":                      "gambatte_libretro",
			"Nintendo - Game Boy":                            "gambatte_libretro",
			"Nintendo - Nintendo 64":                         "mupen64plus_next_libretro",
			"Nintendo - Nintendo Entertainment System":       "fceumm_libretro",
			"Nintendo - Nintendo DS":                         "melonds_libretro",
			"Nintendo - Pokemon Mini":                        "pokemini_libretro",
			"Nintendo - Super Nintendo Entertainment System": "snes9x_libretro",
			"Nintendo - Virtual Boy":                         "mednafen_vb_libretro",
			"Sega - 32X":                                     "picodrive_libretro",
			"Sega - Game Gear":                               "genesis_plus_gx_libretro",
			"Sega - Master System - Mark III":                "genesis_plus_gx_libretro",
			"Sega - Mega Drive - Genesis":                    "genesis_plus_gx_libretro",
			"Sega - Mega-CD - Sega CD":                       "genesis_plus_gx_libretro",
			"Sega - PICO":                                    "picodrive_libretro",
			"Sega - Saturn":                                  "mednafen_saturn_libretro",
			"Sega - SG-1000":                                 "genesis_plus_gx_libretro",
			"SNK - Neo Geo Pocket Color":                     "mednafen_ngp_libretro",
			"SNK - Neo Geo Pocket":                           "mednafen_ngp_libretro",
			"Sony - PlayStation":                             playstationCore,
		},
		Language:             "en",
		FileDirectory:        usr.HomeDir,
		CoresDirectory:       "./cores",
		AssetsDirectory:      "./assets",
		DatabaseDirectory:    "./database",
		SavestatesDirectory:  filepath.Join(xdg.DataHome, "ludo", "savestates"),
		SavefilesDirectory:   filepath.Join(xdg.DataHome, "ludo", "savefiles"),
		ScreenshotsDirectory: filepath.Join(xdg.DataHome, "ludo", "screenshots"),
		SystemDirectory:      filepath.Join(xdg.DataHome, "ludo", "system"),
		PlaylistsDirectory:   filepath.Join(xdg.DataHome, "ludo", "playlists"),
		ThumbnailsDirectory:  filepath.Join(xdg.DataHome, "ludo", "thumbnails"),
		LanguagesDirectory:   "./i18n",
	}
}
