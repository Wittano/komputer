package voice

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

const (
	channels     int = 2                   // 1 for mono, 2 for stereo
	frameRate    int = 48000               // audio sampling rate
	frameSize    int = 960                 // uint16 size of each audio frame
	maxBytes         = (frameSize * 2) * 2 // max size of opus data
	audioBufSize     = 16384
)

func PlayAudio(ctx context.Context, vc *discordgo.VoiceConnection, path string, stop <-chan struct{}) (err error) {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	cmd := exec.Command("ffmpeg", "-i", path, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer output.Close()

	buf := bufio.NewReaderSize(output, audioBufSize)

	// Starts the ffmpeg command
	if err = cmd.Start(); err != nil {
		return
	}

	// wait and release child resources in undefiled time
	go cmd.Wait()

	// Send "speaking" packet over the voice websocket
	if err = vc.Speaking(true); err != nil {
		return
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			return
		}
	}()

	pcm := make(chan []int16, 2)
	defer close(pcm)

	stopPlaying := make(chan struct{})
	defer close(stopPlaying)
	go func(ctx context.Context, vs *discordgo.VoiceConnection, pcm <-chan []int16) {
		sendPCM(ctx, vc, pcm)
		stopPlaying <- struct{}{}
	}(ctx, vc, pcm)

	//when stop is sent, kill ffmpeg
	go func() {
		for {
			select {
			case <-stop:
			case <-ctx.Done():
				cmd.Process.Kill()
				stopPlaying <- struct{}{}
			}
		}
	}()

	for {
		// read data from ffmpeg stdout
		audioBuf := make([]int16, frameSize*channels)
		if err = binary.Read(buf, binary.LittleEndian, &audioBuf); err != nil {
			return
		}

		// Send received PCM to the sendPCM channel
		select {
		case <-ctx.Done():
		case pcm <- audioBuf:
		case <-stopPlaying:
			return
		}
	}
}

func sendPCM(ctx context.Context, v *discordgo.VoiceConnection, pcm <-chan []int16) error {
	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			// read pcm from chan, exit if channel is closed.
			recv, ok := <-pcm
			if !ok {
				return errors.New("PCM Channel closed")
			}

			// try encoding pcm frame with Opus
			opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
			if err != nil {
				return err
			}

			if v.Ready == false || v.OpusSend == nil {
				// Sending errors here might not be suited
				return nil
			}
			// send encoded opus data to the sendOpus channel
			v.OpusSend <- opus
		}
	}
}

func DuractionAudio(path string) (duration time.Duration, err error) {
	cmd := exec.Command("ffprobe", "-i", path, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv='p=0'")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err = cmd.Start(); err != nil {
		return
	}

	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer output.Close()

	rawTime, err := io.ReadAll(output)
	if err != nil {
		return
	}

	return time.ParseDuration(string(rawTime) + "s")
}