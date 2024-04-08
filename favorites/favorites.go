package favorites

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// Game represents a game in the favorites file
type Game struct {
	Path      string // Absolute path of the game on the filesystem
	Name      string // Human readable name of the game, comes from the RDB
	System    string // Name of the game console
	CorePath  string // Absolute path to the libretro core
	Savestate string // Absolute path of the last savestate on this game
}

// Favorites is a list of games
type Favorites []Game

// List is the list of recently played games
var List Favorites

// Push pushes a game onto the stack
func Push(g Game) {
	List = append([]Game{g}, List...)

	// Deduplicate
	l := Favorites{}
	exist := map[string]bool{}
	for _, g := range List {
		if !exist[g.Path] {
			l = append(l, g)
			exist[g.Path] = true
		}
	}
	List = l

	err := Save()
	if err != nil {
		log.Println(err)
	}
}

// Load loads favorites.csv in memory
func Load() error {
	file, err := os.Open(filepath.Join(xdg.DataHome, "ludo", "favorites.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewReader(bufio.NewReader(file))

	List = Favorites{}
	for {
		record, err := wr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		List = append(List, Game{
			Path:     record[0],
			Name:     record[1],
			System:   record[2],
			CorePath: record[3],
		})
	}

	return nil
}

// Save persists the favorites as a csv file
func Save() error {
	file, err := os.Create(filepath.Join(xdg.DataHome, "ludo", "favorites.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(bufio.NewWriter(file))
	defer wr.Flush()

	for _, game := range List {
		wr.Write([]string{
			game.Path,
			game.Name,
			game.System,
			game.CorePath,
		})
	}

	return nil
}
