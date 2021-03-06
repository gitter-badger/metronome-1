package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"

	"github.com/runabove/metronome/src/metronome/constants"
	"github.com/runabove/metronome/src/metronome/core"
)

// Task holds task attributes.
type Task struct {
	GUID      string    `json:"guid",sql:"guid,pk"`
	ID        string    `json:"id",sql:"-"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Schedule  string    `json:"schedule"`
	URN       string    `json:"URN"`
	CreatedAt time.Time `json:"created_at"`
}

// Tasks is a Task list
type Tasks []Task

// ToKafka serialize a Task to Kafka.
func (t *Task) ToKafka() *sarama.ProducerMessage {
	if len(t.GUID) == 0 {
		t.GUID = core.Sha256(t.UserID + t.ID)
	}

	return &sarama.ProducerMessage{
		Topic: constants.KafkaTopicTasks,
		Key:   sarama.StringEncoder(t.GUID),
		Value: sarama.StringEncoder(fmt.Sprintf("%v %v %v %v %v %v", t.UserID, t.ID, t.Schedule, t.URN, url.QueryEscape(t.Name), t.CreatedAt.Unix())),
	}
}

// FromKafka unserialize a Task from Kafka.
func (t *Task) FromKafka(msg *sarama.ConsumerMessage) error {
	key := string(msg.Key)
	segs := strings.Split(string(msg.Value), " ")
	if len(segs) != 6 {
		return fmt.Errorf("unprocessable task(%v) - bad segments", key)
	}

	name, err := url.QueryUnescape(segs[4])
	if err != nil {
		return fmt.Errorf("unprocessable task(%v) - bad name", key)
	}

	timestamp, err := strconv.Atoi(segs[5])
	if err != nil {
		return fmt.Errorf("unprocessable task(%v) - bad timestamp", key)
	}

	t.GUID = key
	t.UserID = segs[0]
	t.ID = segs[1]
	t.Schedule = segs[2]
	t.URN = segs[3]
	t.Name = name
	t.CreatedAt = time.Unix(int64(timestamp), 0)

	return nil
}

// ToJSON serialize a Task as JSON.
func (t *Task) ToJSON() string {
	out, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(out)
}
