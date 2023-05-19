package models

// Status holds the job status
type Status int

func (s Status) String() string {
	switch s {
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case Queued:
		return "Queued"
	case Running:
		return "Running"
	}
	return "Error please check logs"

}
