package shell

import "os/exec"

//Cmd executes shell command and gives back the o
func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			//	println("some error found", err)
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		//	println("some error found2", err)
	}
	return out

}
