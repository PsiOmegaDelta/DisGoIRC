package bot

import (
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Config struct {
	IRC     IRCConfig         `json:"irc"`
	Discord DiscordConfig     `json:"discord"`
	Mapping map[string]string `json:"mapping"`
}

var (
	conf           Config
	inverseMapping map[string]string
	modifiedMapping map[string]string
)

func Init(c Config) {
	conf = c
	inverseMapping = map[string]string{}
	modifiedMapping = map[string]string{}
	for k, v := range conf.Mapping {
		ircChannelPassword := strings.Split(k, " ")
		ircChannel := ircChannelPassword[0]
		inverseMapping[v] = ircChannel
		modifiedMapping[ircChannel] = v
	}
	dInit()
	iInit()
}

func incomingIRC(nick, channel, message string) {
	log.Infof("IRC %s <%s> %s", channel, nick, message)

	discordChan, ok := modifiedMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping IRC:%s to DIS:%s", channel, discordChan)

	message = fmtIrcToDiscord(message)

	dOutgoing(nick, discordChan, message)
}

func incomingDiscord(nick, channel, message string) {
	log.Infof("DIS %s <%s> %s", channel, nick, message)

	ircChan, ok := inverseMapping[channel]
	if !ok {
		return
	}

	log.Debugf("Mapping DIS:%s to IRC:%s", channel, ircChan)

	iOutgoing(nick, ircChan, message)
}

var specialIrc = regexp.MustCompile("|[0-9]{0,2}(,[0-9]{1,2})?")

func fmtReplaceInPairs(msg, find, replace string) string {
	r := regexp.MustCompile(find)
	active := false
	msg = r.ReplaceAllStringFunc(msg, func(a string) string {
		active = !active
		return replace
	})
	if active {
		msg = msg + replace
	}
	return msg
}

func fmtIrcToDiscord(msg string) string {
	msg = specialIrc.ReplaceAllString(msg, "")
	msg = fmtReplaceInPairs(msg, "", "**")
	msg = fmtReplaceInPairs(msg, "", "__")
	msg = fmtReplaceInPairs(msg, "", "*")
	return msg
}
