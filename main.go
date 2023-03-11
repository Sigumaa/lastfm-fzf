package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sigumaa/lfu"
	"github.com/joho/godotenv"
	"github.com/ktr0731/go-fuzzyfinder"
	r "github.com/mattn/go-runewidth"
	"log"
	"os"
	"sync"
)

var cache sync.Map

func main() {
	key, user, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	c := lfu.New(user, key)
	f, err := c.Friends(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	names = append(names, user)
	for i := 0; i < f.Len(); i++ {
		fn, err := f.Name(i)
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, fn)
	}

	_, err = fuzzyfinder.FindMulti(
		names,
		func(i int) string {
			return names[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return userInfo(i, w, h, key, names)
		}))
	if err != nil {
		log.Fatal(err)
	}
}

func userInfo(i, w, h int, key string, names []string) string {
	if v, ok := cache.Load(names[i]); ok {
		if info, ok := v.(string); ok {
			return info
		}
	}

	cu := lfu.New(names[i], key)
	ui, err := cu.Info(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	u1, _ := cu.RecentTracks(context.TODO())
	u2, _ := cu.TopArtists(context.TODO())
	u3, _ := cu.TopTracks(context.TODO())
	u4, _ := cu.TopAlbums(context.TODO())
	rt := u1.RecentTracks.Track[0].Name
	rta := u1.RecentTracks.Track[0].Artist.Text
	rts := fmt.Sprintf("Recent Tracks: %s - %s", rt, rta)
	tar := fmt.Sprintf("Top Artists: %s", u2.TopArtists.Artist[0].Name)
	ttr := fmt.Sprintf("Top Tracks: %s", u3.TopTracks.Track[0].Name)
	tal := fmt.Sprintf("Top Albums: %s", u4.TopAlbums.Album[0].Name)

	lfi := fmt.Sprintf("Last.fm %s Info", names[i])
	pc := fmt.Sprintf("Play Count: %s", ui.PlayCount())
	info := fmt.Sprintf("%s\n%s\n%s\n\n\n%s\n%s\n%s\n", r.Wrap(lfi, w/2-5), r.Wrap(pc, w/2-5), r.Wrap(rts, w/2-5), r.Wrap(tar, w/2-5), r.Wrap(ttr, w/2-5), r.Wrap(tal, w/2-5))

	cache.Store(names[i], info)
	return info
}

func LoadConfig() (string, string, error) {
	if err := godotenv.Load(); err != nil {
		return "", "", errors.New("error loading .env file")
	}
	if os.Getenv("API_KEY") == "" {
		return "", "", errors.New("API_KEY is not set")
	}
	if os.Getenv("USER_NAME") == "" {
		return "", "", errors.New("USER_NAME is not set")
	}
	return os.Getenv("API_KEY"), os.Getenv("USER_NAME"), nil
}
