package webRTC

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pion/opus"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

const (
	// AudioSampleRate is the standard WebRTC sample rate (48kHz)
	AudioSampleRate = 48000
	// AudioChannels is the number of audio channels (mono = 1, stereo = 2)
	AudioChannels = 1
	// BitsPerSample is the number of bits per sample
	BitsPerSample = 16
)

// AudioTrackRecorder handles the recording of an audio track to WAV file
type AudioTrackRecorder struct {
	opusDecoder opus.Decoder
	track       *webrtc.TrackRemote
	samples     []int16
	sampleLock  sync.Mutex
	isRecording bool
	stopChan    chan struct{}
}

// NewAudioTrackRecorder creates a new AudioTrackRecorder instance
func NewAudioTrackRecorder(track *webrtc.TrackRemote) *AudioTrackRecorder {
	return &AudioTrackRecorder{
		track:    track,
		samples:  make([]int16, 0),
		stopChan: make(chan struct{}),
	}
}

// StartRecording begins recording the audio track
func (a *AudioTrackRecorder) StartRecording(peerConnection *webrtc.PeerConnection) error {
	if a.isRecording {
		return fmt.Errorf("recording is already in progress")
	}

	a.isRecording = true
	go a.readPackets()

	// Send PLI on sender
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			select {
			case <-a.stopChan:
				ticker.Stop()
				return
			default:
				rtcpSenders := peerConnection.GetSenders()
				if len(rtcpSenders) > 0 {
					rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{
						&rtcp.PictureLossIndication{MediaSSRC: uint32(a.track.SSRC())},
					})
					if rtcpSendErr != nil {
						fmt.Printf("failed to send RTCP packet: %v\n", rtcpSendErr)
					}
				}
			}
		}
	}()

	return nil
}

// StopRecording stops the recording and saves the WAV file
func (a *AudioTrackRecorder) StopRecording(filename string) error {
	if !a.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	close(a.stopChan)
	a.isRecording = false

	return a.saveToWAV(filename)
}

// readPackets continuously reads RTP packets from the track
func (a *AudioTrackRecorder) readPackets() {
	for {
		select {
		case <-a.stopChan:
			return
		default:
			packet, _, readErr := a.track.ReadRTP()
			fmt.Print("packet: ", packet)
			if readErr != nil {
				fmt.Printf("error reading RTP packet: %v\n", readErr)
				continue
			}

			if err := a.processRTPPacket(packet); err != nil {
				fmt.Printf("error processing RTP packet: %v\n", err)
			}
		}
	}
}

// processRTPPacket processes an individual RTP packet and extracts audio samples
func (a *AudioTrackRecorder) processRTPPacket(packet *rtp.Packet) error {
	// Initialize the Opus decoder if not already done
	fmt.Print("packet.Payload: ", packet.Payload, "\n")
	fmt.Print("track encoding ", a.track.Codec().MimeType, "\n")
	fmt.Print("track da", a.track.Codec().PayloadType, "\n")
	a.opusDecoder = opus.NewDecoder() // Assuming 48kHz and 2 channels

	a.sampleLock.Lock()
	defer a.sampleLock.Unlock()

	// Decode the Opus payload into PCM samples
	const maxDecodedSamples = 960 * 2 * 2 // 960 samples/frame, 2 channels, 2 bytes/sample
	decoded := make([]byte, maxDecodedSamples)

	_, isStereo, err := a.opusDecoder.Decode(packet.Payload, decoded)
	if err != nil {
		return err
	}

	// Determine sample size based on stereo or mono
	sampleSize := 2
	if isStereo {
		sampleSize *= 2
	}

	// Convert the decoded PCM samples to int16 and append to the sample buffer
	for i := 0; i < len(decoded); i += sampleSize {
		if i+1 >= len(decoded) {
			break
		}
		sample := int16(binary.LittleEndian.Uint16(decoded[i:]))
		fmt.Print("sample: ", sample)
		a.samples = append(a.samples, sample)
	}

	return nil
}

// saveToWAV saves the recorded audio samples to a WAV file
func (a *AudioTrackRecorder) saveToWAV(filename string) error {
	a.sampleLock.Lock()
	defer a.sampleLock.Unlock()

	// Create WAV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %v", err)
	}
	defer file.Close()

	// Calculate sizes
	dataSize := len(a.samples) * 2 // 2 bytes per sample
	fileSize := 36 + dataSize

	// Write WAV header
	header := &bytes.Buffer{}

	// RIFF header
	header.WriteString("RIFF")
	binary.Write(header, binary.LittleEndian, uint32(fileSize))
	header.WriteString("WAVE")

	// Format chunk
	header.WriteString("fmt ")
	binary.Write(header, binary.LittleEndian, uint32(16)) // Chunk size
	binary.Write(header, binary.LittleEndian, uint16(1))  // Audio format (PCM)
	binary.Write(header, binary.LittleEndian, uint16(AudioChannels))
	binary.Write(header, binary.LittleEndian, uint32(AudioSampleRate))
	binary.Write(header, binary.LittleEndian, uint32(AudioSampleRate*AudioChannels*BitsPerSample/8)) // Byte rate
	binary.Write(header, binary.LittleEndian, uint16(AudioChannels*BitsPerSample/8))                 // Block align
	binary.Write(header, binary.LittleEndian, uint16(BitsPerSample))

	// Data chunk
	header.WriteString("data")
	binary.Write(header, binary.LittleEndian, uint32(dataSize))

	// Write header to file
	if _, err := file.Write(header.Bytes()); err != nil {
		return fmt.Errorf("failed to write WAV header: %v", err)
	}

	// Write samples to file
	for _, sample := range a.samples {
		err := binary.Write(file, binary.LittleEndian, sample)
		if err != nil {
			return fmt.Errorf("failed to write audio sample: %v", err)
		}
	}

	return nil
}

func handleTrack(track *webrtc.TrackRemote, peerConnection *webrtc.PeerConnection) {
	// Check if it's an audio track
	if track.Kind() == webrtc.RTPCodecTypeAudio {
		// Create new recorder instance
		recorder := NewAudioTrackRecorder(track)

		// Start recording
		if err := recorder.StartRecording(peerConnection); err != nil {
			fmt.Printf("Failed to start recording: %v\n", err)
			return
		}

		// Record for a specific duration or implement your own stopping logic
		time.Sleep(time.Second * 10) // Example: Record for 30 seconds

		// Stop recording and save to file
		if err := recorder.StopRecording("output.wav"); err != nil {
			fmt.Printf("Failed to stop recording: %v\n", err)
			return
		}
	}
}
