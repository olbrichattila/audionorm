# audionorm: Audio Normalization Tool for MP3 to WAV Conversion

`audionorm` is a powerful and easy-to-use command-line tool written in Go. It allows you to normalize the volume of audio files in bulk. The tool processes **MP3 files** from a specified folder and outputs normalized **WAV files** with consistent volume levels, making it ideal for podcast editing, music libraries, and audio projects.

## Key Features

- **Batch Processing**: Normalize multiple MP3 files at once.
- **Custom Normalization Factor**: Set your desired volume level (0 to 1).
- **MP3 to WAV Conversion**: Outputs normalized files in WAV format.
- **Ease of Use**: Simple command-line interface.

## Installation

To install `audionorm`, ensure you have [Go](https://golang.org/) installed, then run:

```
go install github.com/olbrichattila/audionorm/cmd/audionorm@latest
```

## How to Use audionorm
Basic Command Syntax
```
audionorm <path> -factor=<value> -help
```
### Parameters
- path: (Optional) Specifies the folder containing MP3 files. Defaults to the current working directory if not provided.
- -factor: (Optional) A normalization factor between 0 and 1 (e.g., 0.8 for 80% of max volume). Defaults to 1 (no reduction in volume).
- -help: (Optional) Displays usage instructions and exits.

### Examples of Usage
1. Normalize audio in the current directory with default settings:

```
audionorm
```
2. Normalize audio in the current directory with a factor of 0.8:

```
audionorm -factor=0.8
```
3. Normalize audio in a specific folder (./myfolder) with default settings:

```
audionorm ./myfolder
```
4. Normalize audio in ./myfolder with a factor of 0.8:

```
audionorm ./myfolder -factor=0.8
```
5. Display help information and usage instructions:
```
audionorm -help
```

6. Set a low over-amplification tolerance (e.g., 2):
```
audionorm ./folder -tolerance=2
```

7. Set a hight over-amplification tolerance (e.g., 2):
```
audionorm ./folder -tolerance=20
```

## Tolerance Feature
The tolerance feature allows the program to handle naturally occurring high volumes without clipping or distortion. By setting a tolerance level, the program will permit a certain amount of over-amplification, enabling loud segments to remain audible and impactful.

- Key Parameter: -tolerance=<value>
- Range: 0 to 20
- 0 (default): Disables over-amplification.
- 1-20: Gradually increases the permissible over-amplification level, where higher values allow greater tolerance.

### How It Works
When -tolerance is set, the utility adjusts its normalization process to allow some high-volume segments to exceed the standard threshold. This is particularly useful for dynamic audio content, such as music or dialogue, where peaks contribute to the overall experience.

For instance, a tolerance value of 2 permits minor over-amplification, while a value of 20 allows significant amplification for dramatic peaks.

## What Does the Normalization Factor Do?
The normalization factor adjusts the output volume:

- 1: Retains the original volume.
- 0.8: Reduces the volume to 80%.
- Values closer to 0: Significantly lower the volume.

This flexibility ensures that your audio output meets your specific needs, whether you are fine-tuning a podcast or preparing a uniform music library.

Benefits of Using audionorm
- Save time by batch processing audio files.
- Ensure consistent audio quality across all files.
- Easy integration into automated workflows and scripts.

### License
This project is licensed under the MIT License. Feel free to use, modify, and distribute this tool.

Coming soon.
- Bitrate setup,
- wav as input file
