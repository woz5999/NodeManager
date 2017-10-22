package constants

var (
	// Version is the NodeMan version and is set at build time.
	Version string
)

const (
	// UserAgentBase is the base string for the User-Agent HTTP header.
	UserAgentBase = "LogicMonitor NodeMan/"
	// InstanceTerminating is the string indicating an AWS ASG terminate action
	InstanceTerminating = "autoscaling:EC2_INSTANCE_TERMINATING"
	// AsgContinue ASG continue action
	AsgActionContinue = "CONTINUE"
)
