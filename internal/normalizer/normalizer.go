// Package normalizer convert audio files with volume normalization to vaw
package normalizer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/go-mp3"
)

const (
	sourceExt      = ".mp3"
	destinationExt = ".wav"
	stereo         = 2

	outputWavFolder = "./output/wav"
	outputMp3Folder = "./output/mp3"
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
	Normalize(folder string, factor, tolerance float64, convertBackToMp3 bool) error
}

type norma struct {
	factor          float64
	tolerance       float64
	messageCallback interface{}
}

// Normalize starts normalizing files in the folder
func (t *norma) Normalize(folder string, factor, tolerance float64, convertBackToMp3 bool) error {
	t.factor = factor
	t.tolerance = tolerance
	t.message(
		fmt.Sprintf("Start processing folder %s, Normalization factor: %.2f, over amplification tolerance %.2f", folder, factor, tolerance),
	)

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	fileCount := t.numberOfFiles(files, sourceExt)
	t.message(fmt.Sprintf("Number of mp3 files: %d", fileCount))
	processed := 0

	os.MkdirAll(outputWavFolder, os.ModePerm)
	os.MkdirAll(outputMp3Folder, os.ModePerm)

	for _, file := range files {
		if filepath.Ext(file.Name()) == sourceExt {
			processed++
			filePath := filepath.Join(folder, file.Name())
			t.message(fmt.Sprintf("%d/%d Processing %s", fileCount, processed, filePath))
			t.normalizeFile(filePath)
			
			t.normalizeWAV(
				t.replaceFileExtension(filePath, destinationExt),
			)
		}
	}

	t.message("Done converting to Wav and normalizing\n")

	if convertBackToMp3 {
		t.convertAllVawToMP3(outputWavFolder, outputMp3Folder);
	}

	return nil
}

func (t *norma) numberOfFiles(files []fs.DirEntry, ext string) int {
	cnt := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
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

	wavFile, err := os.Create("./output/wav/" +  t.replaceFileExtension(fileName, destinationExt))
	if err != nil {
		return fmt.Errorf("failed to create %s file: %v", destinationExt, err)
	}
	defer wavFile.Close()

	sampleRate := int32(decoder.SampleRate())
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

func (t *norma) message(s string) {
	if t.messageCallback != nil {
		if callback, ok := t.messageCallback.(func(string)); ok {
			callback(s)
		}
	}
}

func (t *norma)replaceFileExtension(filePath, newExt string) string {
	oldExt := filepath.Ext(filePath)
	return strings.TrimSuffix(filePath, oldExt) + newExt
}

func (t *norma)deleteFileIfExists(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
