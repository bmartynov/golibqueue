package golibqueue

import "fmt"

func FormatTopicQueueName(topic, queue string) string {
	return fmt.Sprintf("%s_%s", topic, queue)
}
