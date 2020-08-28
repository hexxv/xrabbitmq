package consumer

type Model uint8

const (
	ModelSimple = iota
	ModelWork
	ModelPublish
	ModelRouting
	ModelTopic
)

func (m Model) String() string {
	switch m {
	case ModelSimple:
		return "simple model consumer"
	case ModelWork:
		return "work model consumer"
	case ModelPublish:
		return "publish model consumer"
	case ModelRouting:
		return "routing model consumer"
	case ModelTopic:
		return "topic model consumer"
	default:
		return "unknown model consumer"
	}
}
