package logic

type Task struct {
	ID           string            `json:"id"`
	SubscribesTo []string          `json:"subscribesTo"`
	Image        string            `json:"image"`
	Env          map[string]string `json:"env"`
	CPU          float64           `json:"cpu"`
	Mem          float64           `json:"mem"`
}

type Event struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type IOPair struct {
	Task  *Task
	Event *Event
}
