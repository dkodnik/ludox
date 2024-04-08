// Package savefiles takes care of saving the game SRAM to the filesystem
package savefiles

import (
	"C"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/libretro/ludo/libretro"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"

	"github.com/libretro/ludo/l10n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var mutex sync.Mutex

// path returns the path of the SRAM file for the current core
func path() string {
	return filepath.Join(
		settings.Current.SavefilesDirectory,
		utils.FileName(state.GamePath)+".srm")
}

// SaveSRAM saves the game SRAM to the filesystem
func SaveSRAM() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !state.CoreRunning {
		txtI18n := l10n.T9(&i18n.Message{ID: "CoreNotRunning", Other: "core not running"})
		return errors.New(txtI18n)
	}

	len := state.Core.GetMemorySize(libretro.MemorySaveRAM)
	ptr := state.Core.GetMemoryData(libretro.MemorySaveRAM)
	if ptr == nil || len == 0 {
		txtI18n := l10n.T9(&i18n.Message{ID: "Unable2GetSRAMAddress", Other: "unable to get SRAM address"})
		return errors.New(txtI18n)
	}

	// convert the C array to a go slice
	bytes := C.GoBytes(ptr, C.int(len))
	err := os.MkdirAll(settings.Current.SavefilesDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	fd, err := os.Create(path())
	if err != nil {
		return err
	}

	_, err = fd.Write(bytes)
	if err != nil {
		fd.Close()
		return err
	}

	err = fd.Close()
	if err != nil {
		return err
	}

	return fd.Sync()
}

// LoadSRAM load the game SRAM from the filesystem
func LoadSRAM() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !state.CoreRunning {
		txtI18n := l10n.T9(&i18n.Message{ID: "CoreNotRunning", Other: "core not running"})
		return errors.New(txtI18n)
	}

	len := state.Core.GetMemorySize(libretro.MemorySaveRAM)
	ptr := state.Core.GetMemoryData(libretro.MemorySaveRAM)
	if ptr == nil || len == 0 {
		txtI18n := l10n.T9(&i18n.Message{ID: "Unable2GetSRAMAddress", Other: "unable to get SRAM address"})
		return errors.New(txtI18n)
	}

	// this *[1 << 30]byte points to the same memory as ptr, allowing to
	// overwrite this memory
	destination := (*[1 << 30]byte)(unsafe.Pointer(ptr))[:len:len]
	source, err := ioutil.ReadFile(path())
	if err != nil {
		return err
	}
	copy(destination, source)

	return nil
}
