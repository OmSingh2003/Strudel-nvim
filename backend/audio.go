package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
	bufferSize = 4096
)

// AudioEngine handles sound synthesis and playback
type AudioEngine struct {
	context *audio.Context
	mu      sync.Mutex
	playing map[string]*audio.Player
}

// Initialize the audio engine
var audioEngine *AudioEngine

func initAudioEngine() error {
	ctx := audio.NewContext(sampleRate)
	audioEngine = &AudioEngine{
		context: ctx,
		playing: make(map[string]*audio.Player),
	}
	log.Println("🔊 Audio engine initialized")
	return nil
}

// generateKick creates a kick drum sound
func generateKick() []byte {
	duration := 0.5 // 500ms
	samples := int(float64(sampleRate) * duration)
	buffer := make([]byte, samples*4) // 16-bit stereo

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		
		// Kick drum: low frequency sine wave with exponential decay
		freq := 60.0 * math.Exp(-t*8) // Frequency starts at 60Hz and decays
		envelope := math.Exp(-t * 5)   // Exponential decay
		sample := envelope * math.Sin(2*math.Pi*freq*t) * 0.5
		
		// Convert to 16-bit signed integer
		value := int16(sample * 32767)
		
		// Write stereo samples (little endian)
		buffer[i*4] = byte(value)
		buffer[i*4+1] = byte(value >> 8)
		buffer[i*4+2] = byte(value)
		buffer[i*4+3] = byte(value >> 8)
	}
	
	return buffer
}

// generateSnare creates a snare drum sound
func generateSnare() []byte {
	duration := 0.2 // 200ms
	samples := int(float64(sampleRate) * duration)
	buffer := make([]byte, samples*4) // 16-bit stereo

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		
		// Snare: noise + tone with quick decay
		noise := (math.Sin(t*12345) + math.Sin(t*23456) + math.Sin(t*34567)) / 3
		tone := math.Sin(2 * math.Pi * 200 * t) // 200Hz tone
		envelope := math.Exp(-t * 15) // Quick decay
		
		sample := envelope * (0.7*noise + 0.3*tone) * 0.3
		
		// Convert to 16-bit signed integer
		value := int16(sample * 32767)
		
		// Write stereo samples (little endian)
		buffer[i*4] = byte(value)
		buffer[i*4+1] = byte(value >> 8)
		buffer[i*4+2] = byte(value)
		buffer[i*4+3] = byte(value >> 8)
	}
	
	return buffer
}

// generateHihat creates a hihat sound
func generateHihat() []byte {
	duration := 0.1 // 100ms
	samples := int(float64(sampleRate) * duration)
	buffer := make([]byte, samples*4) // 16-bit stereo

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		
		// Hihat: high frequency noise with very quick decay
		noise := (math.Sin(t*54321) + math.Sin(t*65432) + math.Sin(t*76543)) / 3
		envelope := math.Exp(-t * 25) // Very quick decay
		
		sample := envelope * noise * 0.2
		
		// Convert to 16-bit signed integer
		value := int16(sample * 32767)
		
		// Write stereo samples (little endian)
		buffer[i*4] = byte(value)
		buffer[i*4+1] = byte(value >> 8)
		buffer[i*4+2] = byte(value)
		buffer[i*4+3] = byte(value >> 8)
	}
	
	return buffer
}

// generateSynth creates a simple synth tone
func generateSynth(frequency float64) []byte {
	duration := 0.3 // 300ms
	samples := int(float64(sampleRate) * duration)
	buffer := make([]byte, samples*4) // 16-bit stereo

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		
		// Simple saw wave with ADSR envelope
		attack := 0.05
		decay := 0.1
		sustain := 0.5
		release := 0.15
		
		var envelope float64
		if t < attack {
			envelope = t / attack
		} else if t < attack+decay {
			envelope = 1.0 - (1.0-sustain)*(t-attack)/decay
		} else if t < duration-release {
			envelope = sustain
		} else {
			envelope = sustain * (duration-t) / release
		}
		
		// Saw wave
		sample := envelope * (2.0*(frequency*t-math.Floor(frequency*t)) - 1.0) * 0.2
		
		// Convert to 16-bit signed integer
		value := int16(sample * 32767)
		
		// Write stereo samples (little endian)
		buffer[i*4] = byte(value)
		buffer[i*4+1] = byte(value >> 8)
		buffer[i*4+2] = byte(value)
		buffer[i*4+3] = byte(value >> 8)
	}
	
	return buffer
}

// playSound plays a synthesized sound
func (ae *AudioEngine) playSound(instrument, note string) error {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	
	var buffer []byte
	
	switch instrument {
	case "kick":
		buffer = generateKick()
	case "snare":
		buffer = generateSnare()
	case "hihat":
		buffer = generateHihat()
	case "synth":
		// Convert note to frequency (simplified)
		frequency := noteToFrequency(note)
		buffer = generateSynth(frequency)
	default:
		// Default to kick for unknown instruments
		buffer = generateKick()
	}
	
	// Create audio player
	player := ae.context.NewPlayerFromBytes(buffer)
	if player == nil {
		return fmt.Errorf("failed to create audio player")
	}
	
	// Play the sound
	player.Play()
	
	// Clean up after playback
	go func() {
		time.Sleep(time.Second) // Wait for sound to finish
		player.Close()
	}()
	
	return nil
}

// noteToFrequency converts a note name to frequency
func noteToFrequency(note string) float64 {
	// Simplified note to frequency conversion
	noteMap := map[string]float64{
		"c4": 261.63, "d4": 293.66, "e4": 329.63, "f4": 349.23,
		"g4": 392.00, "a4": 440.00, "b4": 493.88,
		"c5": 523.25, "d5": 587.33, "e5": 659.25,
	}
	
	if freq, exists := noteMap[note]; exists {
		return freq
	}
	
	// Default to A4 if note not found
	return 440.0
}

