package normalizer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (t *norma) convertAllVawToMP3(sourceFolder, targetFolder string) error {
	files, err := os.ReadDir(sourceFolder)
	if err != nil {
		return err
	}

	fileCount := t.numberOfFiles(files, destinationExt)
	t.message(fmt.Sprintf("Number of vaw files convert back to MP3 files: %d", fileCount))
	processed := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) == destinationExt {
			processed++
			targetFile := filepath.Join(targetFolder, t.replaceFileExtension(file.Name(), sourceExt))
			t.message(fmt.Sprintf("%d/%d converting back to MP3 %s", fileCount, processed, targetFile))
			err := t.deleteFileIfExists(targetFile);
			if err != nil {
				t.message(fmt.Sprintf("file already exists, could not delete `%s` error: %s", file.Name(), err.Error()))
				continue;
			}
			
			err = t.vawToMP3(sourceFolder + "/" + file.Name(), targetFile);
			if err != nil {
			 t.message(fmt.Sprintf("could not convert `%s` error: %s", file.Name(), err.Error()))
			}
		}
	}

	t.message("Done converting back to MP3 files\n")

	return nil
}

func (t *norma) vawToMP3(inputWav, outputMp3 string) error {

	cmd := exec.Command("ffmpeg", "-i", inputWav, outputMp3)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

	// // Alternatiely Use the 'lame' command to convert WAV to MP3, did not work on ubuntu for some non stereo wav files
	// cmd := exec.Command("lame", inputWav, outputMp3)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// return cmd.Run()
}