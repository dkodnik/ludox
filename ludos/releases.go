package ludos

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"

	ntf "github.com/libretro/ludo/notifications"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// UpdatesDir is where releases should be saved to
const UpdatesDir = "/storage/.update/"
const releasesURL = "https://api.github.com/repos/libretro/LudOS/releases"

var client = grab.NewClient()
var downloading bool
var progress float64
var done bool

// Arch is the cpu architecture of LudOS
var Arch = os.Getenv("LIBREELEC_ARCH")

// Version is the version tag of LudOS
var Version = os.Getenv("VERSION")

// GHAsset is an asset attached to a github release
type GHAsset struct {
	Name               string
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GHRelease is a LudOS release hosted on github
type GHRelease struct {
	Name   string
	Assets []GHAsset
}

// GetReleases will get and decode the json from github api, and return the
// list of LudOS releases
func GetReleases() (*[]GHRelease, error) {
	r, err := http.Get(releasesURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	releases := []GHRelease{}
	err = json.NewDecoder(r.Body).Decode(&releases)
	return &releases, err
}

// FilterAssets finds and return the asset matching the LIBREELEC_ARCH
func FilterAssets(assets []GHAsset) *GHAsset {
	for _, asset := range assets {
		if strings.Contains(asset.Name, Arch) {
			return &asset
		}
	}
	return nil
}

// DownloadRelease will download a LudOS release from github
func DownloadRelease(path, url string) {
	if downloading {
		txtI18n := l10n.T9(&i18n.Message{ID: "DownloadAlreadyProgress", Other: "A download is already in progress"})
		ntf.DisplayAndLog(ntf.Error, "Menu", txtI18n)
		return
	}

	txtI18n := l10n.T9(&i18n.Message{ID: "DownloadingUpdate0", Other: "Downloading update 0%%"})
	n := ntf.DisplayAndLog(ntf.Info, "Menu", txtI18n)
	downloading = true
	defer func() { downloading = false }()

	req, err := grab.NewRequest(path, url)
	if err != nil {
		n.Update(ntf.Error, err.Error())
		return
	}

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			txtI18n := l10n.T9(&i18n.Message{ID: "DownloadingUpdate0f", Other: "Downloading update %.0f%%%% "})
			n.Update(ntf.Info, txtI18n, 100*resp.Progress())
			progress = resp.Progress()

		case <-resp.Done:
			// download is complete
			downloading = false
			done = true
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		n.Update(ntf.Error, err.Error())
		downloading = false
		done = false
		return
	}

	txtI18n = l10n.T9(&i18n.Message{ID: "DoneDownloading", Other: "Done downloading. You can now reboot your system."})
	n.Update(ntf.Success, txtI18n)
}

// IsDownloading returns true if the updater is currently downloading a release
func IsDownloading() bool {
	return downloading
}

// IsDone returns true when the download is finished
func IsDone() bool {
	return done
}

// GetProgress returns the download progress
func GetProgress() float64 {
	return progress
}
