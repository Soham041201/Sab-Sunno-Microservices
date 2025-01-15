package webRTC

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Soham041201/Sab-Sunno-Microservices/audio-service/internal/gemini"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"gopkg.in/hraban/opus.v2"
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
	track          *webrtc.TrackRemote
	peerConnection *webrtc.PeerConnection
	audioData      [][]int16
	isRecording    bool
	stopChan       chan struct{}
	lastTimestamp  uint32
	sampleRate     int
}

// NewAudioTrackRecorder creates a new AudioTrackRecorder instance
func NewAudioTrackRecorder(track *webrtc.TrackRemote) *AudioTrackRecorder {
	return &AudioTrackRecorder{
		track:     track,
		audioData: make([][]int16, 0),
		stopChan:  make(chan struct{}),
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

	if err := SaveAudioToWAV(a.audioData, a.sampleRate, "recorded_audio.wav"); err != nil {
		log.Fatal("Error saving audio to WAV file:", err)
	}

	return nil
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
	fmt.Printf("Sample Rate: %d Hz\n", a.track.Codec().ClockRate)
	fmt.Printf("Channels: %d\n", a.track.Codec().Channels)
	channels := int(a.track.Codec().Channels)
	sampleRate := int(a.track.Codec().ClockRate)
	a.sampleRate = sampleRate
	frameSizeMs := float32(60) // if you don't know, go with 60 ms. // Default to 60ms if frame duration is not available
	if packet.Header.Timestamp != 0 {
		// Calculate frame size based on timestamp difference
		if a.lastTimestamp == 0 {
			a.lastTimestamp = packet.Header.Timestamp
		}
		frameSizeMs = float32((packet.Header.Timestamp-a.lastTimestamp)*1000) / float32(a.track.Codec().ClockRate)
		a.lastTimestamp = packet.Header.Timestamp
	}
	frameSize := channels * int(frameSizeMs) * sampleRate / 1000
	pcm := make([]int16, frameSize)

	dec, err := opus.NewDecoder(sampleRate, channels)

	if err != nil {
		fmt.Printf("Error decoding opus packet: %v\n", err)
	}

	n, err := dec.Decode(packet.Payload, pcm)
	if err != nil {
		fmt.Printf("Error decoding opus packet: %v\n", err)
	}
	if n == 0 {
		// Handle cases where no audio data was decoded
		fmt.Println("Warning: No audio data decoded in this packet.")
		return nil // Or return an appropriate error
	}

	fmt.Printf("decoded pcm packet: %d \n ", n)
	samplesPerChannel := n / channels 
	allChannels := make([][]int16, channels)
	for i := 0; i < channels; i++ {
		allChannels[i] = make([]int16, samplesPerChannel) 
	}
	
	for i := 0; i < samplesPerChannel; i++ { 
		for j := 0; j < channels; j++ {
			allChannels[j][i] = pcm[i*channels+j] 
		}
	}

	for _, channel := range allChannels {
		// Convert int16 to bytes (assuming little-endian)
		bytes := make([]byte, len(channel)*2)
		for i, sample := range channel {
			binary.LittleEndian.PutUint16(bytes[i*2:], uint16(sample))
		}
		// Send audio data to Gemini
		err = gemini.HandleGeminiResponse(bytes, sampleRate, a.peerConnection)
		if err != nil {
			return fmt.Errorf("error sending audio to Gemini: %w", err)
		}
	}

	// _, _, err := decode.Decode(packet.Payload, pcm)

	// gemini.HandleGeminiResponse(int16ToBytes(pcm), sampleRate, a.peerConnection)
	// if err != nil {
	// 	fmt.Printf("Error decoding opus packet: %v\n", err)
	// }
	// gemini.SendTextMessage(packet.Payload)
	return nil

}

func writeWAVHeader(w io.Writer, numSamples int, sampleRate int, numChannels int) error {
	// RIFF chunk
	_, err := w.Write([]byte("RIFF"))
	if err != nil {
		return err
	}
	// Chunk size (placeholder)
	_, err = w.Write([]byte{0, 0, 0, 0})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("WAVE"))
	if err != nil {
		return err
	}

	// fmt sub-chunk
	_, err = w.Write([]byte("fmt "))
	if err != nil {
		return err
	}
	// Sub-chunk size (16 for PCM)
	_, err = w.Write([]byte{0x10, 0, 0, 0})
	if err != nil {
		return err
	}
	// Audio format (1 for PCM)
	_, err = w.Write([]byte{0x01, 0})
	if err != nil {
		return err
	}
	// Number of channels
	err = binary.Write(w, binary.LittleEndian, uint16(numChannels))
	if err != nil {
		return err
	}
	// Sample rate
	err = binary.Write(w, binary.LittleEndian, uint32(sampleRate))
	if err != nil {
		return err
	}
	// Byte rate (sample rate * numChannels * bytes per sample)
	byteRate := sampleRate * numChannels * 2
	err = binary.Write(w, binary.LittleEndian, uint32(byteRate))
	if err != nil {
		return err
	}
	// Block align (numChannels * bytes per sample)
	blockAlign := numChannels * 2
	err = binary.Write(w, binary.LittleEndian, uint16(blockAlign))
	if err != nil {
		return err
	}
	// Bits per sample
	_, err = w.Write([]byte{0x10, 0})
	if err != nil {
		return err
	}

	// data sub-chunk
	_, err = w.Write([]byte("data"))
	if err != nil {
		return err
	}
	// Sub-chunk size (numSamples * numChannels * bytes per sample)
	dataSize := numSamples * numChannels * 2
	err = binary.Write(w, binary.LittleEndian, uint32(dataSize))
	if err != nil {
		return err
	}

	return nil
}

// saveToWAV saves the recorded audio samples to a WAV file
func SaveAudioToWAV(audioData [][]int16, sampleRate int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating WAV file: %w", err)
	}
	defer file.Close()

	// Write WAV header
	err = writeWAVHeader(file, len(audioData[0]), sampleRate, len(audioData))
	if err != nil {
		return fmt.Errorf("error writing WAV header: %w", err)
	}

	// Write audio data
	for _, channel := range audioData {
		for _, sample := range channel {
			err = binary.Write(file, binary.LittleEndian, sample)
			if err != nil {
				return fmt.Errorf("error writing audio data: %w", err)
			}
		}
	}

	return nil
}

func HandleTrack(track *webrtc.TrackRemote, peerConnection *webrtc.PeerConnection) {
	// Check if it's an audio track
	if track.Kind() == webrtc.RTPCodecTypeAudio {
		// Create new recorder instance
		recorder := NewAudioTrackRecorder(track)

		recorder.peerConnection = peerConnection

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
