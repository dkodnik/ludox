package l10n

import (
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// The list of supported languages.
var supportLng = map[string]string{
	"en": "English",
	"ru": "Русский",
}

type Lng struct {
	l    *i18n.Localizer
	t    string // the current language is "en", "ru", ...
	path string // the path to the *.toml files
}

var lng *Lng

func Get0() *Lng {
	return lng
}

func checkLang(in string) string {
	name, ok := supportLng[in]
	if ok {
		return name
	}
	return in
}

func GetAllNameLang() []string {
	ret := []string{}
	for _, f := range supportLng {
		ret = append(ret, f)
	}

	return ret
}

// NameLang to Lang
func GetNameLang2Lang(in string) string {
	ret := ""
	for i, f := range supportLng {
		if f == in {
			ret = i
		}
	}
	return ret
}

// Lang to NameLang
func GetLang2NameLang(in string) string {
	return checkLang(in)
}

// List of supported languages
func GetAllLang(path string) []string {
	// Returns languages, depending on the availability of files in the 'path'
	// and possibly supported languages of the 'supportLng' system. In any
	// case, there is always an 'en'.
	retdat := []string{}

	m, err := filepath.Glob(path + "/active.*.toml")
	if err != nil {
		log.Println(err)
		return []string{} //en - def
	}

	for _, val := range m {
		val = strings.Replace(val, ".toml", "", -1)
		setLng := val[len(val)-2:]

		_, ok := supportLng[setLng]
		if ok {
			retdat = append(retdat, setLng)
		}
	}

	if len(retdat) == 0 {
		return []string{} //en - def
	}

	return retdat
}

func Init(setLng string, path string) {
	filters := GetAllLang(path)
	def := "en"

	if !slices.Contains(filters, setLng) {
		setLng = def
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, f := range filters {
		// No need to load active.en.toml since we are providing default translations.
		bundle.MustLoadMessageFile(path + "/active." + f + ".toml")
	}

	lng = &Lng{
		path: path,
		t:    setLng,
		l:    i18n.NewLocalizer(bundle, setLng),
	}
}

func ReInit(setLng string) {
	filters := GetAllLang(lng.path)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, f := range filters {
		bundle.MustLoadMessageFile(lng.path + "/active." + f + ".toml")
	}

	lng.l = i18n.NewLocalizer(bundle, setLng)
}

// Translate(t9)
func T9(id2 *i18n.Message) string {
	return lng.l.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: id2,
	})
}
