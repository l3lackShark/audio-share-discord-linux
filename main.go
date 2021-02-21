package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/l3lackShark/audio-share-discord-linux/pactl"
)

var ids []string

func main() {
	if runtime.GOOS != "linux" {
		panic("OS IS NOT SUPPORTED")
	}

	if !pactl.CheckForPulseaudio() {
		panic("Pulseaudio is not installed")
	}

	setupCloseHandler()
	defer pactl.UnloadCables(&ids)
	rawSinks := pactl.GetSinks()
	i := 1
	for _, sink := range rawSinks {
		fmt.Println(fmt.Sprintf("%d: %s", i, sink))
		i++
	}
	fmt.Println("Select your output device!")
	var answer int
	_, err := fmt.Scanln(&answer)
	for err != nil || answer > len(rawSinks) || answer < 1 {
		fmt.Println("Select your output device!")
		_, err = fmt.Scanln(&answer)
	}
	spl := strings.Split(rawSinks[answer-1], "Name: ")
	spl = strings.Split(spl[1], "\n")

	parsedAlsaSink := spl[0]
	fmt.Println("Parsed alsa sink:", parsedAlsaSink)
	ids = pactl.CreateVirualCables()
	ids = append(ids, pactl.LoadListenCalbe(parsedAlsaSink))

	rawSources := pactl.GetSources()

	i = 1
	for _, source := range rawSources {
		fmt.Println(fmt.Sprintf("%d: %s", i, source))
		i++
	}
	fmt.Println("Select your input device!")
	var answerI int
	_, err = fmt.Scanln(&answerI)
	for err != nil || answerI > len(rawSources) || answerI < 1 {
		fmt.Println("Select your input device!")
		_, err = fmt.Scanln(&answerI)
	}
	spl = strings.Split(rawSources[answerI-1], "Name: ")
	spl = strings.Split(spl[1], "\n")

	parsedAlsaSource := spl[0]
	fmt.Println("Parsed alsa source:", parsedAlsaSource)
	pactl.GetMicVolume(parsedAlsaSource)
	ids = append(ids, pactl.LoadVirtualMic(parsedAlsaSource)...)

	fmt.Println()
	fmt.Println(`Your main mic is now muted and is set to 100% volume! (EARRAPE WARNING), change input device in Discord to "<...>VirtMic<...>", unmute your main mic and set it's appropriate sound level (pavucontrol/pulsemixer). Then move any programs that you want audio streaming ON to "FunnelSink" and change the volume of "OutputMixer" according to your friend's liking. You should also disable automatic input sensitivity and set it's value to the lowest possible in the Discord settings.`)
	fmt.Println()
	fmt.Println("Press ENTER or Ctrl + C to quit")
	var x string
	_, _ = fmt.Scanln(&x)

}

//Ctrl + C handler
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Exiting gracefully...")
		if len(ids) > 0 {
			pactl.UnloadCables(&ids)
		}
		os.Exit(0)
	}()
}
