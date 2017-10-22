package constants

var (
	// Version is the NodeMan version and is set at build time.
	Version string
)

const (
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase       = "LogicMonitor NodeMan/"
	InstanceTerminating = "autoscaling:EC2_INSTANCE_TERMINATING"
)
