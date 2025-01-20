package sound

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var player *oto.Player

// loadMP3 loads an MP3 file and returns an audio player.
func loadMP3(filePath string) (*oto.Player, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	c, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return nil, err
	}

	player := c.NewPlayer()

	return player, nil
}

// playSound plays the sound.
func playSound() {
	if player != nil {
		player.Write()
	}
}

func TestPlayer() {
	var er error
	s, er = mp3.DecodeF32(bytes.NewReader(raudio.Ragtime_mp3))
	if er != nil {
		return nil, er
	}

	var err error
	player, err = loadMP3("assets/sound/Coin.mp3")
	if err != nil {
		log.Fatalf("Failed to load MP3: %v", err)
	}
	defer player.Close()

}
