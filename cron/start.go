package cron

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Cepave/ops-common/model"
	"github.com/Cepave/ops-common/utils"
	"github.com/Cepave/ops-nqm-agent-updater/g"
	"github.com/toolkits/file"
)

func prequisite(da *model.DesiredAgent) error {
	var err error
	if err = InsureDesiredAgentDirExists(da); err != nil {
		return err
	}
	if err = InsureNewVersionFiles(da); err != nil {
		return err
	}
	if err = Untar(da); err != nil {
		return err
	}
	return nil
}

func StartNQMAgent(da *model.DesiredAgent) {
	if err := prequisite(da); err != nil {
		return
	}
	nqmBinPath := filepath.Join(da.AgentVersionDir, da.Name)

	moduleStatus := g.CheckModuleStatus(nqmBinPath)
	if moduleStatus == g.NotRunning {
		fmt.Print("Starting [", da.Name, "]...")

		logPath := filepath.Join(da.AgentVersionDir, g.LogFile)
		LogOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Println("Error in opening file:", err)
			return
		}
		defer LogOutput.Close()

		cmd := exec.Command(nqmBinPath)
		cmd.Stdout = LogOutput
		cmd.Stderr = LogOutput
		dir, _ := os.Getwd()
		cmd.Dir = dir
		cmd.Start()
		fmt.Println("successfully!!")
		time.Sleep(1 * time.Second)
		moduleStatus = g.CheckModuleStatus(nqmBinPath)
		if moduleStatus == g.NotRunning {
			log.Fatalln("** start failed **")
		}
	}
}

func Untar(da *model.DesiredAgent) error {
	cmd := exec.Command("tar", "zxf", da.TarballFilename)
	cmd.Dir = da.AgentVersionDir
	err := cmd.Run()
	if err != nil {
		log.Println("tar zxf", da.TarballFilename, "fail", err)
		return err
	}

	return nil
}

func InsureNewVersionFiles(da *model.DesiredAgent) error {
	if FilesReady(da) {
		return nil
	}
	content, err := ioutil.ReadFile("./password")
	password := strings.Trim(string(content), "\n")
	if err != nil {
		panic(err)
	}
	downloadTarballCmd := exec.Command("wget", "--no-check-certificate", "--auth-no-challenge", "--user=owl", "--password="+password, da.TarballUrl, "-O", da.TarballFilename)
	downloadTarballCmd.Dir = da.AgentVersionDir
	err = downloadTarballCmd.Run()
	if err != nil {
		log.Println("wget -q --no-check-certificate --auth-no-challenge --user=owl --password="+password, da.TarballUrl, "-O", da.TarballFilename, "fail", err)
		return err
	}

	downloadMd5Cmd := exec.Command("wget", "--no-check-certificate", "--auth-no-challenge", "--user=owl", "--password="+password, da.Md5Url, "-O", da.Md5Filename)
	downloadMd5Cmd.Dir = da.AgentVersionDir
	err = downloadMd5Cmd.Run()
	if err != nil {
		log.Println("wget -q --no-check-certificate --auth-no-challenge --user=owl --password="+password, da.Md5Url, "-O", da.Md5Filename, "fail", err)
		return err
	}

	if utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename) {
		return nil
	} else {
		return fmt.Errorf("md5sum -c fail")
	}
}

func FilesReady(da *model.DesiredAgent) bool {
	if !file.IsExist(da.Md5Filepath) {
		return false
	}

	if !file.IsExist(da.TarballFilepath) {
		return false
	}

	return utils.Md5sumCheck(da.AgentVersionDir, da.Md5Filename)
}

func InsureDesiredAgentDirExists(da *model.DesiredAgent) error {
	err := file.InsureDir(da.AgentDir)
	if err != nil {
		log.Println("insure dir", da.AgentDir, "fail", err)
		return err
	}

	err = file.InsureDir(da.AgentVersionDir)
	if err != nil {
		log.Println("insure dir", da.AgentVersionDir, "fail", err)
	}
	return err
}
