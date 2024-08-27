//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Topic -trimprefix Topic
package linode

type Topic int

const (
	TopicStatus Topic = iota
	TopicBoot
	TopicReboot
	TopicShutdown
)
