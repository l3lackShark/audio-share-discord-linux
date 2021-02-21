package pactl

import (
	"fmt"
	"strings"

	"github.com/l3lackShark/audio-share-discord-linux/shell"
)

//Checks for pulseaudio
func CheckForPulseaudio() bool {
	raw := string(shell.Cmd("pgrep pulseaudio", true))
	if raw != "" {
		return true
	}
	return false
}

//GetSinks returns list of available alsa sinks
func GetSinks() []string {
	raw := string(shell.Cmd("pactl list sinks | grep -A 1 Name", true))
	blocks := strings.Split(raw, "--")
	return blocks
}

//GetSources returns list of available alsa sources
func GetSources() []string {
	raw := string(shell.Cmd("pactl list sources | grep -A 1 Name", true))
	blocks := strings.Split(raw, "--")

	return blocks
}

//GetMicVolume is used for a workaround to stop 100% mic reinitialisation
func GetMicVolume(input string) string {
	raw := string(shell.Cmd(fmt.Sprintf("pactl list sources | grep -A 15 %s", input), true))
	spl := strings.Split(raw, "Volume:")
	spl = strings.Split(raw, "/")
	out := strings.TrimSpace(spl[2])
	out = strings.TrimSuffix(out, "%")
	shell.Cmd(fmt.Sprintf("pactl set-source-mute %s 1", input), true)

	return out
}

//RestoreMicVolume is used for a workaround to stop 100% mic reinitialisation
func RestoreMicVolume(deviceName string, value string) {
	shell.Cmd(fmt.Sprintf("pactl set-source-volume %s %s%%", deviceName, value), true)
}

//CreateVirualCables creates virtual cables that we will pass our sinks and sources to
func CreateVirualCables() []string {
	var out []string
	out = append(out, string(shell.Cmd(`pactl load-module module-null-sink sink_name=VirtSoundFirst sink_properties=device.description="FunnelSink"`, true)))
	out = append(out, string(shell.Cmd(`pactl load-module module-null-sink sink_name=VirtSoundSecond sink_properties=device.description="Do_Not_Touch_Sink"`, true)))
	out = append(out, string(shell.Cmd(`pactl load-module module-null-sink sink_name=VirtSoundThird sink_properties=device.description="OutputMixer"`, true)))
	return out
}

//UnloadCables removes the virtual cables from pa
func UnloadCables(input *[]string) error {
	for _, id := range *input {
		raw := string(shell.Cmd(fmt.Sprintf("pactl unload-module %s", id), true))
		if raw != "" {
			return fmt.Errorf("ERROR: %s", raw)
		}
	}
	return nil
}

//LoadListenCalbe loads the loopback cable for output
func LoadListenCalbe(input string) string {
	return string(shell.Cmd(fmt.Sprintf(`pactl load-module module-loopback source=VirtSoundFirst.monitor sink=%s latency_msec=1`, input), true))
}

//LoadVirtualMic mixes everything to get the input device conaining everything
func LoadVirtualMic(input string) []string {
	var out []string
	out = append(out, string(shell.Cmd(`pactl load-module module-loopback source=VirtSoundFirst.monitor sink=VirtSoundThird latency_msec=20`, true)))
	out = append(out, string(shell.Cmd(`pactl load-module module-loopback source=VirtSoundThird.monitor sink=VirtSoundSecond latency_msec=20`, true)))
	out = append(out, string(shell.Cmd(fmt.Sprintf(`pactl load-module module-loopback source=%s sink=VirtSoundSecond latency_msec=20`, input), true)))
	out = append(out, string(shell.Cmd(`pactl load-module module-virtual-source source_name=VirtMic master=VirtSoundSecond.monitor`, true)))
	return out
}
