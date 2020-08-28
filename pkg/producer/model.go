package producer

type Model uint8

const (
	ModelSimple = iota
	ModelWork
	ModelPublish
	ModelRouting
	ModelRoutingDynamic
	ModelTopic
	ModelTopicDynamic
)

func (m Model) String() string {
	switch m {
	case ModelSimple:
		return "simple model producer"
	case ModelWork:
		return "work model producer"
	case ModelPublish:
		return "publish model producer"
	case ModelRouting:
		return "routing model producer"
	case ModelRoutingDynamic:
		return "routingDynamic model producer"
	case ModelTopic:
		return "topic model producer"
	case ModelTopicDynamic:
		return "topicDynamic model producer"
	default:
		return "unknown model producer"
	}
}
