// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Dep struct {
	Name     string `json:"name"`
	Detached bool   `json:"detached"`
}

type DetachedTask struct {
	ID       string        `json:"id"`
	Tasks    []*Task       `json:"tasks"`
	Outputs  []*Output     `json:"outputs"`
	Status   string        `json:"status"`
	StartAt  time.Time     `json:"startAt"`
	Duration time.Duration `json:"duration"`
	EndAt    *time.Time    `json:"endAt"`
}

type Elk struct {
	Version string                 `json:"version"`
	Env     map[string]interface{} `json:"env"`
	EnvFile string                 `json:"envFile"`
	Vars    map[string]interface{} `json:"vars"`
	Tasks   []*Task                `json:"tasks"`
}

type Log struct {
	Out    string `json:"out"`
	Format string `json:"format"`
	Error  string `json:"error"`
}

type Output struct {
	Task  string   `json:"task"`
	Out   []string `json:"out"`
	Error []string `json:"error"`
}

type RunConfig struct {
	Start    *time.Time     `json:"start"`
	Deadline *time.Time     `json:"deadline"`
	Timeout  *time.Duration `json:"timeout"`
	Delay    *time.Duration `json:"delay"`
}

type Task struct {
	Title       string                 `json:"title"`
	Tags        []string               `json:"tags"`
	Name        string                 `json:"name"`
	Cmds        []*string              `json:"cmds"`
	Env         map[string]interface{} `json:"env"`
	Vars        map[string]interface{} `json:"vars"`
	EnvFile     string                 `json:"envFile"`
	Description string                 `json:"description"`
	Dir         string                 `json:"dir"`
	Log         *Log                   `json:"log"`
	Sources     *string                `json:"sources"`
	Deps        []*Dep                 `json:"deps"`
	IgnoreError bool                   `json:"ignoreError"`
}

type TaskDep struct {
	Name        string `json:"name"`
	Detached    bool   `json:"detached"`
	IgnoreError bool   `json:"ignoreError"`
}

type TaskInput struct {
	Name        string                 `json:"name"`
	Title       *string                `json:"title"`
	Tags        []string               `json:"tags"`
	Cmds        []string               `json:"cmds"`
	Env         map[string]interface{} `json:"env"`
	Vars        map[string]interface{} `json:"vars"`
	EnvFile     *string                `json:"envFile"`
	Description *string                `json:"description"`
	Dir         *string                `json:"dir"`
	Log         *TaskLog               `json:"log"`
	Sources     *string                `json:"sources"`
	Deps        []*TaskDep             `json:"deps"`
	IgnoreError *bool                  `json:"ignoreError"`
}

type TaskLog struct {
	Out    string         `json:"out"`
	Error  string         `json:"error"`
	Format *TaskLogFormat `json:"format"`
}

type TaskProperties struct {
	Vars        map[string]interface{} `json:"vars"`
	Env         map[string]interface{} `json:"env"`
	EnvFile     *string                `json:"envFile"`
	IgnoreError *bool                  `json:"ignoreError"`
}

type DetachedTaskStatus string

const (
	DetachedTaskStatusWaiting DetachedTaskStatus = "waiting"
	DetachedTaskStatusRunning DetachedTaskStatus = "running"
	DetachedTaskStatusSuccess DetachedTaskStatus = "success"
	DetachedTaskStatusError   DetachedTaskStatus = "error"
)

var AllDetachedTaskStatus = []DetachedTaskStatus{
	DetachedTaskStatusWaiting,
	DetachedTaskStatusRunning,
	DetachedTaskStatusSuccess,
	DetachedTaskStatusError,
}

func (e DetachedTaskStatus) IsValid() bool {
	switch e {
	case DetachedTaskStatusWaiting, DetachedTaskStatusRunning, DetachedTaskStatusSuccess, DetachedTaskStatusError:
		return true
	}
	return false
}

func (e DetachedTaskStatus) String() string {
	return string(e)
}

func (e *DetachedTaskStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DetachedTaskStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DetachedTaskStatus", str)
	}
	return nil
}

func (e DetachedTaskStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TaskLogFormat string

const (
	TaskLogFormatAnsic       TaskLogFormat = "ANSIC"
	TaskLogFormatUnixDate    TaskLogFormat = "UnixDate"
	TaskLogFormatRubyDate    TaskLogFormat = "RubyDate"
	TaskLogFormatRfc822      TaskLogFormat = "RFC822"
	TaskLogFormatRfc822z     TaskLogFormat = "RFC822Z"
	TaskLogFormatRfc850      TaskLogFormat = "RFC850"
	TaskLogFormatRfc1123     TaskLogFormat = "RFC1123"
	TaskLogFormatRfc1123z    TaskLogFormat = "RFC1123Z"
	TaskLogFormatRfc3339     TaskLogFormat = "RFC3339"
	TaskLogFormatRFC3339Nano TaskLogFormat = "RFC3339Nano"
	TaskLogFormatKitchen     TaskLogFormat = "Kitchen"
)

var AllTaskLogFormat = []TaskLogFormat{
	TaskLogFormatAnsic,
	TaskLogFormatUnixDate,
	TaskLogFormatRubyDate,
	TaskLogFormatRfc822,
	TaskLogFormatRfc822z,
	TaskLogFormatRfc850,
	TaskLogFormatRfc1123,
	TaskLogFormatRfc1123z,
	TaskLogFormatRfc3339,
	TaskLogFormatRFC3339Nano,
	TaskLogFormatKitchen,
}

func (e TaskLogFormat) IsValid() bool {
	switch e {
	case TaskLogFormatAnsic, TaskLogFormatUnixDate, TaskLogFormatRubyDate, TaskLogFormatRfc822, TaskLogFormatRfc822z, TaskLogFormatRfc850, TaskLogFormatRfc1123, TaskLogFormatRfc1123z, TaskLogFormatRfc3339, TaskLogFormatRFC3339Nano, TaskLogFormatKitchen:
		return true
	}
	return false
}

func (e TaskLogFormat) String() string {
	return string(e)
}

func (e *TaskLogFormat) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TaskLogFormat(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TaskLogFormat", str)
	}
	return nil
}

func (e TaskLogFormat) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
