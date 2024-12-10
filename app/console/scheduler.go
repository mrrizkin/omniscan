package console

import (
	"github.com/mrrizkin/omniscan/app/providers/scheduler"
	"github.com/mrrizkin/omniscan/app/repositories"
)

func Schedule(
	schedule *scheduler.Scheduler,
	eStatementRepository *repositories.EStatementRepository,
) {
	// Example usage:
	// schedule.Add("@every 1m", func() {
	// 	log.Info("Scheduled task ran", "at", time.Now().Format(time.RFC3339))
	// })
	//
	// Add your scheduled tasks here. Use cron syntax for defining intervals.
	// Refer to the scheduler documentation for more advanced usage.

	// deleting expired e-statements
	schedule.Add("* * * * *", func() {
		eStatementRepository.Bomb()
	})
}
