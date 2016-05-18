package cron

import (
	"log"
	"path"
	"time"

	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-nqm-agent-updater/g"
	f "github.com/toolkits/file"
)

func BuildHeartbeatRequest(hostname string, agentDirs []string) model.HeartbeatRequest {
	req := model.HeartbeatRequest{Hostname: hostname}

	realAgents := []*model.RealAgent{}
	now := time.Now().Unix()

	for _, agentDir := range agentDirs {
		if path.Base(agentDir) != "nqm-agent" {
			continue
		}
		if g.NQMRunningVersion == "<UNDEFINED>" {
			continue
		}

		status := ""
		nqmVersionPath := path.Join(g.SelfDir, agentDir, g.NQMRunningVersion)
		switch g.CheckModuleStatus(nqmVersionPath) {
		case g.Running:
			status = "Running"
		case g.NotRunning:
			status = "NotRunning"
		}

		realAgent := &model.RealAgent{
			Name:      agentDir,
			Version:   g.NQMRunningVersion,
			Status:    status,
			Timestamp: now,
		}

		realAgents = append(realAgents, realAgent)
	}

	req.RealAgents = realAgents
	return req
}

func ListAgentDirs() ([]string, error) {
	agentDirs, err := f.DirsUnder(g.SelfDir)
	if err != nil {
		log.Println("list dirs under", g.SelfDir, "fail", err)
	}
	return agentDirs, err
}
