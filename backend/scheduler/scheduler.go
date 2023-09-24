package scheduler

import "github.com/robfig/cron"

// StartScheduler is a function to run service automatically within request time without hit service
func StartScheduler() {
	c := cron.New()

	c.Start()
}
