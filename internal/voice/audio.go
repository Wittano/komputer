package voice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes      = (frameSize * 2) * 2 // max size of opus data
)

// TODO Added context to canceling action when users doesn't listening audio (e.g. force stop or every users from channel
func PlayAudio(vc *discordgo.VoiceConnection, path string, stop <-chan bool) (err error) {
	cmd := exec.Command("ffmpeg", "-i", path, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	buf := bufio.NewReaderSize(output, 16384)

	// Starts the ffmpeg command
	err = cmd.Start()
	if err != nil {
		return
	}

	// wait and release child resources in undefiled time
	go cmd.Wait()

	//when stop is sent, kill ffmpeg
	go func() {
		<-stop
		cmd.Process.Kill()
	}()

	// Send "speaking" packet over the voice websocket
	err = vc.Speaking(true)
	if err != nil {
		return
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			return
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	stopPlaying := make(chan bool)
	go func() {
		sendPCM(vc, send)
		stopPlaying <- true
	}()

	for {
		// read data from ffmpeg stdout
		audioBuf := make([]int16, frameSize*channels)
		err = binary.Read(buf, binary.LittleEndian, &audioBuf)
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			return
		}
		if err != nil {
			return
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audioBuf:
		case <-stopPlaying:
			return
		}
	}
}

func sendPCM(v *discordgo.VoiceConnection, pcm chan []int16) (err error) {
	if pcm == nil {
		return nil
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		return
	}

	for {

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
