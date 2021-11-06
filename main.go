package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/gookit/color"
	"github.com/l3lackShark/audio-share-discord-linux/pactl"
)

var ids []string

func main() {
	color.Println("\n <op=reverse> audio-share-discord-linux </> \n")

	if runtime.GOOS != "linux" {
		panic("OS IS NOT SUPPORTED")
	}

	if !pactl.CheckForPulseaudio() {
		panic("pulseaudio is not running!")
	}

	setupCloseHandler()
	defer pactl.UnloadCables(&ids)
	rawSinks := pactl.GetSinks()
	i := 1
	for _, sink := range rawSinks {
		fmt.Println(fmt.Sprintf("%d: %s", i, sink))
		i++
	}
	fmt.Printf("Select your output device [%d-%d]: ", 1, len(rawSinks))
	var answer int
	color.Set(color.White)
	_, err := fmt.Scanln(&answer)
	color.Reset()
	for err != nil || answer > len(rawSinks) || answer < 1 {
		fmt.Printf("Select your output device [%d-%d]: ", 1, len(rawSinks))
		color.Set(color.White)
		_, err = fmt.Scanln(&answer)
		color.Reset()
	}
	spl := strings.Split(rawSinks[answer-1], "Name: ")
	spl = strings.Split(spl[1], "\n")

	parsedAlsaSink := spl[0]
	ids = pactl.CreateVirualCables()
	ids = append(ids, pactl.LoadListenCalbe(parsedAlsaSink))

	rawSources := pactl.GetSources()

	i = 1
	for _, source := range rawSources {
		fmt.Println(fmt.Sprintf("%d: %s", i, source))
		i++
	}
	fmt.Printf("Select your input device [%d-%d]: ", 1, len(rawSources))
	var answerI int
	color.Set(color.White)
	_, err = fmt.Scanln(&answerI)
	color.Reset()
	for err != nil || answerI > len(rawSources) || answerI < 1 {
		fmt.Printf("Select your input device [%d-%d]: ", 1, len(rawSources))
		color.Set(color.White)
		_, err = fmt.Scanln(&answerI)
		color.Reset()
	}
	spl = strings.Split(rawSources[answerI-1], "Name: ")
	spl = strings.Split(spl[1], "\n")

	parsedAlsaSource := spl[0]
	pactl.GetMicVolume(parsedAlsaSource)
	ids = append(ids, pactl.LoadVirtualMic(parsedAlsaSource)...)

	color.Println(`
Your main mic is now muted and is set to 100% volume! (EARRAPE WARNING),

  <fg=cyan>1.</> Change input device in Discord to <fg=white;op=bold>Virtual Source VirtMic on Monitor of Do_Not_Touch_Sink</>
  <fg=cyan>2.</> Unmute your main mic and set it's appropriate sound level (pavucontrol/pulsemixer)
  <fg=cyan>3.</> Move any programs that you want audio streaming ON to <fg=white;op=bold>FunnelSink</>
     and change the volume of <fg=lightWhite;op=bold>OutputMixer</> according to your friend's liking. 
  <fg=cyan>4.</> You should also disable automatic input sensitivity and set it's value to the
     lowest possible in the Discord settings.

Press ENTER or Ctrl + C to quit`)
	var x string
	_, _ = fmt.Scanln(&x)

}

//Ctrl + C handler
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rExiting gracefully...")
		if len(ids) > 0 {
			pactl.UnloadCables(&ids)
		}
		color.Reset()
		os.Exit(0)
	}()
}
