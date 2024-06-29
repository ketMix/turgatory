package game

import (
	"fmt"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/kettek/ebijam24/assets"
)

type Track struct {
	player *audio.Player
	volume float64
	pan    float64

	// panning goes from -1 to 1
	// -1: 100% left channel, 0% right channel
	// 0: 100% both channels
	// 1: 0% left channel, 100% right channel
	panstream *StereoPanStream
}

func (t *Track) Play() {
	t.player.Rewind()
	t.player.Play()
}

func (t *Track) Pause() {
	t.player.Pause()
}

func (t *Track) SetVolume(volume float64) {
	t.volume = volume
	t.player.SetVolume(volume)
}

func (t *Track) SetPan(pan float64) {
	t.pan = pan
	t.panstream.SetPan(pan)
}

type AudioController struct {
	audioContext     *audio.Context
	tracks           map[RoomKind]*Track
	backgroundTracks []*Track
	sfx              map[string]*Track
	tracksPaused     bool
	sfxPaused        bool
}
type PanVol struct {
	Pan float64
	Vol float64
}

const VOL_MULT = 2.0

func NewAudioController() *AudioController {
	audioContext := audio.NewContext(44100)
	tracks := make(map[RoomKind]*Track)

	// Get list of room kinds, then create a mapping to the bytes
	for i := 0; i < int(RoomKindEnd); i++ {
		roomKind := RoomKind(i)
		name := roomKind.String()

		stream, err := assets.LoadSound("room", name)
		if err != nil {
			// it's fine, just means the track doesn't exist
			fmt.Println("Error loading track ", name, err)
		}

		if stream != nil {
			panstream := NewStereoPanStream(audio.NewInfiniteLoop(stream, stream.Length()))
			panstream.SetPan(0.0)

			player, err := audioContext.NewPlayer(panstream)

			player.SetVolume(0)
			if err != nil {
				fmt.Println("Error creating player for track ", name, err)
			}
			tracks[roomKind] = &Track{
				player:    player,
				panstream: panstream,
				volume:    0,
				pan:       0,
			}
		}
	}

	// Add background track
	trackNames := []string{"bass", "kick"}
	backgroundTracks := make([]*Track, 0)
	for _, name := range trackNames {
		stream, err := assets.LoadSound("room", name)
		if err != nil {
			fmt.Println("Error loading background tracks ", err)
		}

		panstream := NewStereoPanStream(audio.NewInfiniteLoop(stream, stream.Length()))
		panstream.SetPan(0.0)
		player, err := audioContext.NewPlayer(panstream)
		player.SetVolume(VOL_MULT)
		if err != nil {
			fmt.Println("Error creating player for background track ", err)
		}

		backgroundTrack := &Track{
			player:    player,
			panstream: panstream,
			volume:    0,
			pan:       0,
		}
		backgroundTracks = append(backgroundTracks, backgroundTrack)
	}

	return &AudioController{
		audioContext:     audioContext,
		tracks:           tracks,
		backgroundTracks: backgroundTracks,
		sfx:              make(map[string]*Track),
		tracksPaused:     true,
		sfxPaused:        false,
	}
}

func (a *AudioController) PlayRoomTracks() {
	a.tracksPaused = false
	for _, track := range a.tracks {
		track.Play()
	}
	for _, track := range a.backgroundTracks {
		track.Play()
	}
}

func (a *AudioController) PauseRoomTracks() {
	for _, track := range a.tracks {
		track.Pause()
	}
	for _, track := range a.backgroundTracks {
		track.Pause()
	}
	a.tracksPaused = true
}

func (a *AudioController) SetVol(roomKind RoomKind, volume float64) {
	if track, ok := a.tracks[roomKind]; ok {
		track.SetVolume(volume * VOL_MULT)
	}
}

func (a *AudioController) SetPan(roomKind RoomKind, pan float64) {
	if track, ok := a.tracks[roomKind]; ok {
		track.SetPan(pan)
	}
}

// For the tracks that we have,
// set the volume and pan
// If the track doesn't exist, set it to 0
// Allows muting rooms that are removed
func (a *AudioController) SetStoryPanVol(roomPanVol map[RoomKind]PanVol) {
	for track := range a.tracks {
		if panvol, ok := roomPanVol[track]; ok {
			a.SetVol(track, panvol.Vol)
			a.SetPan(track, panvol.Pan)
		} else {
			a.SetVol(track, 0)
		}
	}
}

func (a *AudioController) PlaySfx(name string, vol float64, pan float64) {
	if a.sfxPaused {
		return
	}
	if sfx, ok := a.sfx[name]; ok {
		sfx.SetVolume(vol)
		sfx.SetPan(pan)

		sfx.Play()
	} else {
		stream, err := assets.LoadSound("sfx", name)
		if err != nil {
			fmt.Println("Error loading sfx ", name, err)
			return
		}

		panstream := NewStereoPanStream(stream)
		panstream.SetPan(0)

		player, err := a.audioContext.NewPlayer(panstream)
		if err != nil {
			fmt.Println("Error creating player for sfx ", name, err)
			return
		}

		sfx := &Track{
			player:    player,
			panstream: panstream,
			volume:    1,
			pan:       0,
		}

		a.sfx[name] = sfx
		a.PlaySfx(name, vol, pan)
	}
}

/**
 *	This section copied from https://github.com/hajimehoshi/ebiten/blob/main/examples/audiopanning/main.go
 */

// StereoPanStream is an audio buffer that changes the stereo channel's signal
// based on the Panning.
type StereoPanStream struct {
	io.ReadSeeker
	pan float64 // -1: left; 0: center; 1: right
	buf []byte
}

func (s *StereoPanStream) Read(p []byte) (int, error) {
	// If the stream has a buffer that was read in the previous time, use this first.
	var bufN int
	if len(s.buf) > 0 {
		bufN = copy(p, s.buf)
		s.buf = s.buf[bufN:]
	}

	readN, err := s.ReadSeeker.Read(p[bufN:])
	if err != nil && err != io.EOF {
		return 0, err
	}

	// Align the buffer size in multiples of 4. The extra part is pushed to the buffer for the
	// next time.
	totalN := bufN + readN
	extra := totalN - totalN/4*4
	s.buf = append(s.buf, p[totalN-extra:totalN]...)
	alignedN := totalN - extra

	// This implementation uses a linear scale, ranging from -1 to 1, for stereo or mono sounds.
	// If pan = 0.0, the balance for the sound in each speaker is at 100% left and 100% right.
	// When pan is -1.0, only the left channel of the stereo sound is audible, when pan is 1.0,
	// only the right channel of the stereo sound is audible.
	// https://docs.unity3d.com/ScriptReference/AudioSource-panStereo.html
	ls := math.Min(s.pan*-1+1, 1)
	rs := math.Min(s.pan+1, 1)
	for i := 0; i < alignedN; i += 4 {
		lc := int16(float64(int16(p[i])|int16(p[i+1])<<8) * ls)
		rc := int16(float64(int16(p[i+2])|int16(p[i+3])<<8) * rs)

		p[i] = byte(lc)
		p[i+1] = byte(lc >> 8)
		p[i+2] = byte(rc)
		p[i+3] = byte(rc >> 8)
	}
	return alignedN, err
}

func (s *StereoPanStream) SetPan(pan float64) {
	s.pan = math.Min(math.Max(-1, pan), 1)
}

func (s *StereoPanStream) Pan() float64 {
	return s.pan
}

// NewStereoPanStream returns a new StereoPanStream with a buffered src.
//
// The src's format must be linear PCM (16bits little endian, 2 channel stereo)
// without a header (e.g. RIFF header). The sample rate must be same as that
// of the audio context.
func NewStereoPanStream(src io.ReadSeeker) *StereoPanStream {
	return &StereoPanStream{
		ReadSeeker: src,
	}
}
