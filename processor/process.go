package processor

import (
	"net/url"
	"os"
	"os/exec"
	"path"
)

var name string

func getName() string {
	if name != "" {
		return name
	}
	name = os.Getenv("FFMPEG_PATH")
	if name == "" {
		name = "ffmpeg"
	}
	return name
}

var output string

func getOutput() (string, error) {
	if output != "" {
		return output, nil
	}
	key := os.Getenv("RTMP_KEY")
	url, err := url.Parse(os.Getenv("RTMP_URL"))
	if err != nil {
		return "", err
	}
	url.Path = path.Join(url.Path, key)
	output = url.String()
	return output, nil
}

var (
	args = []string{
		"-re",
		"-stream_loop", "-1",
		"-loglevel", "fatal",
		"-i", "",
		"-preset", "ultrafast",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "flv",
	}
	inputIndex = 6
)

func getArgs(input string) ([]string, error) {
	toReturn := args
	toReturn[inputIndex] = input
	output, err := getOutput()
	if err != nil {
		return toReturn, err
	}
	toReturn = append(toReturn, output)
	return toReturn, nil
}

var process *exec.Cmd

func Stop() (bool, error) {
	if Processing() {
		err := process.Process.Kill()
		process = nil
		return true, err
	}
	return false, nil
}

func Process(input string) error {
	if process != nil {
		if _, err := Stop(); err != nil {
			return err
		}
	}
	args, err := getArgs(input)
	if err != nil {
		return err
	}
	process = exec.Command(getName(), args...)
	process.Stdout = os.Stdout
	process.Stdin = os.Stdin
	err = process.Start()
	if err == nil {
		go func() {
			process.Wait()
			Stop()
		}()
	}
	return err
}

func Processing() bool {
	return process != nil
}
