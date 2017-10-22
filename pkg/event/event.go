package event

// Event an AWS ASG scaling event message body
type Event struct {
	EC2InstanceID        string `json:"EC2InstanceId"`
	AccountID            string `json:"AccountId"`
	AutoScalingGroupName string `json:"AutoScalingGroupName"`
	LifecycleActionToken string `json:"LifecycleActionToken"`
	LifecycleHookName    string `json:"LifecycleHookName"`
	LifecycleTransition  string `json:"LifecycleTransition"`
	NotificationMetadata string `json:"NotificationMetadata"`
	RequestID            string `json:"RequestID"`
	Service              string `json:"Service"`
	Time                 string `json:"Time"`
}
