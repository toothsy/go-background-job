package constants

const (
	// constants for job status
	Queued = iota
	Running
	Completed
	Failed
	DeleteJob
	// constants for job type
	// mainly authenticate or upload
	Authenticate
	Upload
)
