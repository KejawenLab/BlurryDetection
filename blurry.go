// Author: Muhamad Surya Iksanudin<surya.kejawen@gmail.com>
package identify

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type BlurryDetection struct {
	IsBlur bool
	Score  float64
}

func (b BlurryDetection) Detect(ImagePath string) (BlurryDetection, error) {
	if !b.commandAvailable() {
		return b, errors.New("ImageMagick is not available.")
	}

	if !b.validate(ImagePath) {
		return b, errors.New("Please remove space in your image.")
	}

	output, err := b.run(ImagePath)
	if err != nil {
		return b, err
	}

	result := strings.Split(output, "\n")
	b.Score, _ = strconv.ParseFloat(strings.Trim(strings.Split(result[len(result)-2], "(")[1], ")"), 64)
	if 0.1 > b.Score {
		b.IsBlur = true
	} else {
		b.IsBlur = false
	}

	return b, nil
}

func (b BlurryDetection) commandAvailable() bool {
	cmd := exec.Command("identify", "-version")
	_, err := cmd.Output()

	return err == nil
}

func (b BlurryDetection) validate(ImagePath string) bool {
	return len(strings.Split(ImagePath, " ")) == 1
}

func (b BlurryDetection) run(ImagePath string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("identify -verbose %s | grep deviation", ImagePath))
	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(stdout), err
}
