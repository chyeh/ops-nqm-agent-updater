package cron

import (
	"log"

	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-nqm-agent-updater/g"
)

func HandleHeartbeatResponse(respone *model.HeartbeatResponse) {
	if respone.ErrorMessage != "" {
		log.Println("receive error message:", respone.ErrorMessage)
		return
	}

	das := respone.DesiredAgents
	if das == nil || len(das) == 0 {
		return
	}

	for _, da := range das {
		da.FillAttrs(g.SelfDir)
		HandleDesiredAgent(da)
	}
}

func HandleDesiredAgent(da *model.DesiredAgent) {
	if da.Cmd == "start" {
		StartNQMAgent(da)
	} else if da.Cmd == "stop" {
		StopNQMAgent(da)
	} else {
		log.Println("unknown cmd", da)
	}
}
