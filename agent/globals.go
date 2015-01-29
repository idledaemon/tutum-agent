package agent

import (
	"log"
	"os"
)

var (
	FlagDebugMode   *bool
	FlagLogToStdout *bool
	FlagStandalone  *bool
	FlagDockerHost  *string
	FlagDockerOpts  *string
	FlagTutumHost   *string
	FlagTutumToken  *string
	FlagTutumUUID   *string

	Conf                      Configuration
	Logger                    *log.Logger
	DockerProcess             *os.Process
	ScheduleToTerminateDocker = false
	DockerBinaryURL           = "https://files.tutum.co/packages/docker/latest.json"
)

const (
	VERSION               = "0.11.3"
	defaultCertCommonName = ""
	defaultDockerHost     = "tcp://0.0.0.0:2375"
	defaultTutumHost      = "https://dashboard.tutum.co/"
)

const (
	TutumHome = "/etc/tutum/agent"
	DockerDir = "/usr/lib/tutum"
	LogDir    = "/var/log/tutum"

	DockerSymbolicLink     = "/usr/bin/docker"
	DockerLogFileName      = "docker.log"
	TutumLogFileName       = "agent.log"
	KeyFileName            = "key.pem"
	CertFileName           = "cert.pem"
	CAFileName             = "ca.pem"
	ConfigFileName         = "tutum-agent.conf"
	DockerBinaryName       = "docker"
	DockerNewBinaryName    = "docker.new"
	DockerNewBinarySigName = "docker.new.sig"

	RegEndpoint       = "api/agent/node/"
	DockerDefaultHost = "unix:///var/run/docker.sock"

	MaxWaitingTime    = 200 //seconds
	HeartBeatInterval = 5   //second

	RenicePriority = -10
)