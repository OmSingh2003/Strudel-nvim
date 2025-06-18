package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Strudel pattern structure
type StrudelPattern struct {
	Name     string      `json:"name"`
	Pattern  string      `json:"pattern"`
	BPM      int         `json:"bpm"`
	Cycle    int         `json:"cycle"`
	Elements []string    `json:"elements"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Audio event structure
type AudioEvent struct {
	Instrument string  `json:"instrument"`
	Note       string  `json:"note"`
	Velocity   float64 `json:"velocity"`
	Duration   float64 `json:"duration"`
	Time       float64 `json:"time"`
}

// Pattern evaluation result
type EvaluationResult struct {
	Success    bool         `json:"success"`
	Pattern    StrudelPattern `json:"pattern"`
	Events     []AudioEvent `json:"events"`
	Message    string       `json:"message"`
	Error      string       `json:"error,omitempty"`
	Timestamp  time.Time    `json:"timestamp"`
}

// evaluateStrudelPattern processes a Strudel pattern string and returns a JSON result
func evaluateStrudelPattern(patternStr string) string {
	log.Printf("Evaluating Strudel pattern: %s", patternStr)
	
	result := EvaluationResult{
		Timestamp: time.Now(),
	}
	
	// Parse the pattern
	pattern, err := parseStrudelPattern(patternStr)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Parse error: %v", err)
		result.Message = "Failed to parse Strudel pattern"
	} else {
		result.Success = true
		result.Pattern = pattern
		result.Events = generateAudioEvents(pattern)
		result.Message = fmt.Sprintf("Successfully evaluated pattern '%s' with %d events", pattern.Name, len(result.Events))
		
		// Send to audio engine (placeholder)
		go sendToAudioEngine(result.Events)
	}
	
	// Convert result to JSON
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"success": false, "error": "JSON encoding error: %v"}`, err)
	}
	
	return string(jsonResult)
}

// parseStrudelPattern parses a Strudel pattern string into a structured format
func parseStrudelPattern(patternStr string) (StrudelPattern, error) {
	pattern := StrudelPattern{
		BPM:      120, // Default BPM
		Cycle:    1,   // Default cycle
		Metadata: make(map[string]interface{}),
	}
	
	// Clean the pattern string
	patternStr = strings.TrimSpace(patternStr)
	
	if patternStr == "" {
		return pattern, fmt.Errorf("empty pattern")
	}
	
	// Extract pattern name if present (e.g., $: d1 $ sound "bd sn")
	nameRegex := regexp.MustCompile(`^(\w+)\s*\$\s*(.+)$`)
	if matches := nameRegex.FindStringSubmatch(patternStr); len(matches) == 3 {
		pattern.Name = matches[1]
		pattern.Pattern = matches[2]
	} else {
		pattern.Name = "unnamed"
		pattern.Pattern = patternStr
	}
	
	// Parse BPM if present
	bpmRegex := regexp.MustCompile(`bpm\s+(\d+)`)
	if matches := bpmRegex.FindStringSubmatch(pattern.Pattern); len(matches) == 2 {
		if bpm, err := strconv.Atoi(matches[1]); err == nil {
			pattern.BPM = bpm
		}
	}
	
	// Extract elements from the pattern
	pattern.Elements = extractPatternElements(pattern.Pattern)
	
	return pattern, nil
}

// extractPatternElements extracts individual elements from a Strudel pattern
func extractPatternElements(patternStr string) []string {
	// Simple pattern extraction - looks for quoted strings and individual tokens
	var elements []string
	
	// Extract quoted strings (e.g., "bd sn hh")
	quotedRegex := regexp.MustCompile(`"([^"]+)"`)
	quotedMatches := quotedRegex.FindAllStringSubmatch(patternStr, -1)
	for _, match := range quotedMatches {
		// Split quoted content by spaces
		quotedElements := strings.Fields(match[1])
		elements = append(elements, quotedElements...)
	}
	
	// If no quoted strings found, split by spaces
	if len(elements) == 0 {
		elements = strings.Fields(patternStr)
	}
	
	// Filter out Strudel keywords
	keywords := map[string]bool{
		"sound": true, "note": true, "gain": true, "pan": true,
		"delay": true, "reverb": true, "$": true, "|": true,
		"fast": true, "slow": true, "rev": true, "bpm": true,
	}
	
	var filteredElements []string
	for _, element := range elements {
		cleanElement := strings.Trim(element, "()[]{}")
		if !keywords[cleanElement] && cleanElement != "" {
			filteredElements = append(filteredElements, cleanElement)
		}
	}
	
	return filteredElements
}

// generateAudioEvents converts pattern elements into audio events
func generateAudioEvents(pattern StrudelPattern) []AudioEvent {
	var events []AudioEvent
	
	if len(pattern.Elements) == 0 {
		return events
	}
	
	// Calculate timing based on BPM
	beatDuration := 60.0 / float64(pattern.BPM) // Duration of one beat in seconds
	stepDuration := beatDuration / float64(len(pattern.Elements))
	
	// Generate events for each element
	for i, element := range pattern.Elements {
		if element == "~" || element == "" {
			continue // Skip rests
		}
		
		event := AudioEvent{
			Instrument: detectInstrument(element),
			Note:       element,
			Velocity:   0.8, // Default velocity
			Duration:   stepDuration * 0.9, // 90% of step duration
			Time:       float64(i) * stepDuration,
		}
		
		events = append(events, event)
	}
	
	return events
}

// detectInstrument tries to determine the instrument type from the element
func detectInstrument(element string) string {
	drumMap := map[string]string{
		"bd": "kick", "kick": "kick",
		"sn": "snare", "snare": "snare",
		"hh": "hihat", "hihat": "hihat",
		"oh": "openhat", "openhat": "openhat",
		"cp": "clap", "clap": "clap",
		"cy": "cymbal", "cymbal": "cymbal",
	}
	
	if instrument, exists := drumMap[strings.ToLower(element)]; exists {
		return instrument
	}
	
	// Check if it's a musical note
	noteRegex := regexp.MustCompile(`^[a-g][#b]?[0-9]?$`)
	if noteRegex.MatchString(strings.ToLower(element)) {
		return "synth"
	}
	
	return "sample" // Default to sample
}

// sendToAudioEngine sends audio events to the audio engine
func sendToAudioEngine(events []AudioEvent) {
	if audioEngine == nil {
		log.Printf("Audio engine not initialized")
		return
	}
	
	log.Printf("🎵 Playing %d audio events:", len(events))
	
	// Schedule and play each event
	for i, event := range events {
		log.Printf("  🎶 Event %d: %s (%s) at %.3fs", i+1, event.Note, event.Instrument, event.Time)
		
		// Schedule the event to play at the correct time
		go func(evt AudioEvent, delay time.Duration) {
			if delay > 0 {
				time.Sleep(delay)
			}
			
			// Play the sound
			if err := audioEngine.playSound(evt.Instrument, evt.Note); err != nil {
				log.Printf("Error playing sound: %v", err)
			} else {
				log.Printf("🔊 Played: %s (%s)", evt.Note, evt.Instrument)
			}
		}(event, time.Duration(event.Time*float64(time.Second)))
	}
}

