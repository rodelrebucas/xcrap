package scheduler

import (
	"github.com/carlescere/scheduler"
	"fmt"
	)

func RunAtMidnight() {
	job := func() {
		fmt.Println("Dummy function")
	   }
	   scheduler.Every(5).Seconds().Run(job)
	   scheduler.Every().Day().Run(job)
	   scheduler.Every().Sunday().At("08:30").Run(job)
}