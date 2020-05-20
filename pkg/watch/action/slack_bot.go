package action

type (
	SlackActionRepository interface {
		PutSlackPostTimestamp(jobID string, itemID string, timestamp int64) error
		GetSlackPostTimestamp(jobID, itemID string) (int64, error)
	}
)
