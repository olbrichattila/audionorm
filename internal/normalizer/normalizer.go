package normalizer

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/go-mp3"
)

const (
	sourceExt      = ".mp3"
	destinationExt = ".wav"
	stereo         = 2
)

// New creates a normalized, if message callback is provided then it will be called with status messages, Signature is func(string)
func New(messageCallback interface{}) AudioNormalizer {
	if messageCallback != nil {
		_, ok := messageCallback.(func(string))
		if !ok {
			panic("Message callback is not a func(string)")
		}
	}

	return &norma{messageCallback: messageCallback}
}

// AudioNormalizer abstracts the normalizer logic
type AudioNormalizer interface {
	Normalize(folder string, factor float64) error
}

type norma struct {
	factor          float64
	messageCallback interface{}
}

// Normalize starts normalizing files in the folder
func (t *norma) Normalize(folder string, factor float64) error {
	t.factor = factor
	t.message(
		fmt.Sprintf("Start processing folder %s, Normalization factor: %.2f", folder, factor),
	)

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	fileCount := t.numberOfFiles(files)
	t.message(fmt.Sprintf("Number of mp3 files: %d", fileCount))
	processed := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) == sourceExt {
			processed++
			filePath := filepath.Join(folder, file.Name())
			t.message(fmt.Sprintf("%d/%d Processing %s", fileCount, processed, filePath))
			t.normalizeFile(filePath)
			t.normalizeWAV(filePath + destinationExt)
		}
	}

	t.message("Done\n")

	return nil
}

func (t *norma) numberOfFiles(files []fs.DirEntry) int {
	cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == sourceExt {
			cnt++
		}
	}
	return cnt
}

func (t *norma) normalizeFile(fileName string) error {
	mp3File, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open MP3 file: %v", err)
	}
	defer mp3File.Close()

	decoder, err := mp3.NewDecoder(mp3File)
	if err != nil {
		return fmt.Errorf("failed to decode MP3: %v", err)
	}

	wavFile, err := os.Create(fileName + destinationExt)
	if err != nil {
		return fmt.Errorf("failed to create %s file: %v", destinationExt, err)
	}
	defer wavFile.Close()

	// sampleRate := int32(decoder.SampleRate())
	sampleRate := int32(44100)
	numChannels := int16(stereo)
	bitsPerSample := int16(16)
	data := []byte{}
	t.message(fmt.Sprintf("  - sample rate: %d", sampleRate))

	wavHeader := t.createWAVHeader(int32(len(data)), sampleRate, numChannels, bitsPerSample)
	if _, err := wavFile.Write(wavHeader); err != nil {
		return fmt.Errorf("failed to write %s header: %v", destinationExt, err)
	}

	buf := make([]byte, 4096)
	for {
		n, err := decoder.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read MP3 data: %v", err)
		}
		if n == 0 {
			break
		}

		wavFile.Write(buf[:n])
	}
	return nil

}

func (t *norma) normalizeWAV(wavPath string) error {
	wavFile, err := os.OpenFile(wavPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file: %v", destinationExt, err)
	}
	defer wavFile.Close()

	// Skip WAV headers (44 bytes)
	wavFile.Seek(44, 0)

	buf := make([]byte, 4096)
	var maxSample int16 = 0

	// First pass: Find the maximum sample value
	for {
		n, err := wavFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("calculate max volume: failed to read %s data: %v", destinationExt, err)
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i += 2 {
			sample := int16(binary.LittleEndian.Uint16(buf[i : i+2]))
			if t.abs(sample) > maxSample {
				maxSample = t.abs(sample)
			}
		}
	}

	// Calculate normalization factor
	factor := float64(math.MaxInt16) / float64(maxSample) * float64(t.factor)
	t.message(fmt.Sprintf("  - calculated factor: %2.f, max Sample: %d", factor, maxSample))

	// Second pass: Apply normalization
	wavFile.Seek(44, 0)
	for {
		n, err := wavFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("apply: failed to read %s data: %v", destinationExt, err)
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i += 2 {
			sample := int16(binary.LittleEndian.Uint16(buf[i : i+2]))
			normalizedSample := int16(float64(sample) * factor)
			binary.LittleEndian.PutUint16(buf[i:i+2], uint16(normalizedSample))
		}

		// Write normalized data back
		wavFile.Seek(-int64(n), io.SeekCurrent)
		wavFile.Write(buf[:n])
	}

	return nil
}

func (*norma) createWAVHeader(dataSize, sampleRate int32, numChannels, bitsPerSample int16) []byte {
	blockAlign := numChannels * bitsPerSample / 8
	byteRate := sampleRate * int32(blockAlign)

	header := make([]byte, 44)
	copy(header[0:], []byte("RIFF"))
	binary.LittleEndian.PutUint32(header[4:], uint32(36+dataSize))
	copy(header[8:], []byte("WAVE"))
	copy(header[12:], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:], 16)
	binary.LittleEndian.PutUint16(header[20:], 1)
	binary.LittleEndian.PutUint16(header[22:], uint16(numChannels))
	binary.LittleEndian.PutUint32(header[24:], uint32(sampleRate))
	binary.LittleEndian.PutUint32(header[28:], uint32(byteRate))
	binary.LittleEndian.PutUint16(header[32:], uint16(blockAlign))
	binary.LittleEndian.PutUint16(header[34:], uint16(bitsPerSample))
	copy(header[36:], []byte("data"))
	binary.LittleEndian.PutUint32(header[40:], uint32(dataSize))

	return header
}

func (*norma) abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func (t *norma) message(s string) {
	if t.messageCallback != nil {
		if callback, ok := t.messageCallback.(func(string)); ok {
			callback(s)
		}
	}
}
