package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/utils"
	"os"
	"path"
	"strconv"
	"time"
)

var service utils.External

const xrayGid = 3003

type ServiceCommand struct{}

func (this *ServiceCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify operation, available operation [start|stop|restart|status]").WithPrefix("service").WithPathObj(*this)
	}
	if len(args) > 1 {
		return errors.New("service: too many arguments")
	}
	switch args[0] {
	case "start":
		log.HandleInfo("service: starting xray")
		if err := startService(); err != nil {
			return err
		}
		log.HandleInfo("service: xray is running, pid is " + getServicePid())
	case "stop":
		log.HandleInfo("service: stopping xray")
		stopService()
		log.HandleInfo("service: xray is stopped")
	case "restart":
		log.HandleInfo("service: restarting xray")
		stopService()
		if err := startService(); err != nil {
			return err
		}
		log.HandleInfo("service: xray is running, pid is " + getServicePid())
	case "status":
		pidStr := getServicePid()
		if len(pidStr) > 0 {
			log.HandleInfo("service: xray is running, pid is " + pidStr)
		} else {
			log.HandleInfo("service: xray is stopped")
		}
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [start|stop|restart|status]").WithPrefix("service").WithPathObj(*this)
	}
	return nil
}

// startService start xray service
func startService() error {
	listenFlag := false
	serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "error.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return errors.New("open xray log file failed, ", err).WithPrefix("service")
	}
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return errors.New("open xray config file failed, ", err).WithPrefix("service")
	} else {
		if confInfo.IsDir() {
			service = utils.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.Core, "run", "-confdir", builds.Config.XrayHelper.CoreConfig)
		} else {
			service = utils.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.Core, "run", "-c", builds.Config.XrayHelper.CoreConfig)
		}
	}
	service.AppendEnv("XRAY_LOCATION_ASSET=" + builds.Config.XrayHelper.BaseDir)
	if err := service.SetUidGid(0, xrayGid); err != nil {
		return err
	}
	service.Start()
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if utils.CheckPort("tcp", "127.0.0.1", builds.Config.Proxy.TproxyPort) {
			listenFlag = true
			break
		}
	}
	if listenFlag {
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
			_ = service.Kill()
			return errors.New("write xray pid failed, ", err).WithPrefix("service")
		}
	} else {
		_ = service.Kill()
		log.HandleDebug(service.Err())
		return errors.New("start xray service failed").WithPrefix("service")
	}
	return nil
}

// stopService stop xray service
func stopService() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		pid, _ := strconv.Atoi(string(pidFile))
		if serviceProcess, err := os.FindProcess(pid); err == nil {
			_ = serviceProcess.Kill()
			_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid"))
		} else {
			log.HandleDebug(err)
		}
	} else {
		log.HandleDebug(err)
	}
}

// getServicePid get xray pid from pid file
func getServicePid() string {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "xray.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		return string(pidFile)
	} else {
		log.HandleDebug(err)
	}
	return ""
}
