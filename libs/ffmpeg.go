package libs

import (
	"log"
	"os/exec"

	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

var FfmpegCmd *fluentffmpeg.Command

func init() {
	// check
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute ffmpeg: %v", err)
	}

	FfmpegCmd = fluentffmpeg.NewCommand("")
}