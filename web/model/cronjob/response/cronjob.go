package response

// @Author: morris
// NAME
// SCHEDULE
// SUSPEND
// ACTIVE
// LAST SCHEDULE
// AGE
type CronJob struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	//CronJob的执行时间表
	Schedule string `json:"schedule"`
	//于暂停CronJob的调度。当设置为true时，CronJob将停止生成新的Job，但保留现有的Job
	Suspend *bool `json:"suspend"`
	//CronJob当前处于活跃状态的Job数量
	Active int `json:"active"`
	//这个字段记录了CronJob上一次成功调度的时间
	LastSchedule int64 `json:"lastSchedule"`
	//创建时间
	Age int64 `json:"age"`
}
