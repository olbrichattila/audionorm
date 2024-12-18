package normalizer

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

func (t *norma) normalizeWAV(wavPath string) error {
	fmt.Println(wavPath);
	wavFile, err := os.OpenFile(wavPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file: %v", destinationExt, err)
	}
	defer wavFile.Close()

	// Skip WAV headers (44 bytes)
	wavFile.Seek(44, 0)

	buf := make([]byte, 4096)
	var maxSample int16
	var sampleCount int64
	peekBuffer := make([]int64, 32768)

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
			abs := t.abs(sample)
			if abs > maxSample {
				maxSample = abs
			}
			peekBuffer[abs]++
			sampleCount++
		}
	}

	// Calculate peek with some threshold
	if t.tolerance > 0 {
		for i := 32767; i > 0; i-- {
			percentage := float64(peekBuffer[i]) / float64(sampleCount) * 10000000
			if percentage > t.tolerance {
				t.message(fmt.Sprintf("  - tolerance at %d position", 32767-i))
				maxSample = int16(i)
				break
			}
		}
	}

	// Calculate normalization factor
	factor := float64(math.MaxInt16) / float64(maxSample) * float64(t.factor)
	t.message(fmt.Sprintf("  - calculated factor: %0.4f, max Sample: %d", factor, maxSample))

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
			normalizedSample := float64(sample) * factor

			// In case of over amplification, remove some distortion
			if normalizedSample > 32767 {
				normalizedSample = 32767
			} else if normalizedSample < -32767 {
				normalizedSample = -32767
			}

			binary.LittleEndian.PutUint16(buf[i:i+2], uint16(int16(normalizedSample)))
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
