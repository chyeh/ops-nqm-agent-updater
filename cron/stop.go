package cron

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-nqm-agent-updater/g"
)

func StopNQMAgent(da *model.DesiredAgent) {
	nqmBinPath := filepath.Join(da.AgentVersionDir, da.Name)

	moduleStatus := g.CheckModuleStatus(nqmBinPath)
	if moduleStatus == g.NotRunning {
		// Skip stopping if the module is stopped
		return
	}

	fmt.Print("Stopping [", nqmBinPath, "] ")

	pidStr, _ := g.CheckModulePid(nqmBinPath)

	cmd := exec.Command("kill", "-9", pidStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	fmt.Println("with PID [", pidStr, "]...successfully!!")
	time.Sleep(1 * time.Second)

	moduleStatus = g.CheckModuleStatus(nqmBinPath)
	if moduleStatus == g.Running {
		fmt.Println("** stop failed **")
		return
	}
}
