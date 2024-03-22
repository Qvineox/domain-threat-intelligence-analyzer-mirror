package record

import (
	"log/slog"
	"os/exec"
	"strings"
)

func PingDomain(name string) bool {
	out, err := exec.Command("ping", name, "-c 4", "-i 3", "-w 5").Output()
	if err != nil {
		slog.Warn("ping failed: " + err.Error())
		return false
	}

	if strings.Contains(string(out), "Destination Host Unreachable") {
		return false
	}

	return true
}
