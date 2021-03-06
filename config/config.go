package config

import (
	"os"
	"regexp"
	"strconv"
)

var DbDsn string = os.Getenv("PRESENT_DB_DSN")

var SlackIncomingWebhookUrl string = os.Getenv("PRESENT_SLACK_INCOMMING_URL")

var Names []string = []string{"present"}
var Wait int = 15 * 60
var NoopLimit int = 0
var RankingsHour int = -1
var AccesslogUrlBase string = ""

func init() {

	names := os.Getenv("PRESENT_NAME")
	if names != "" {
		Names = regexp.MustCompile(",").Split(names, -1)
	}

	wait := os.Getenv("PRESENT_WAIT")
	if w, err := strconv.Atoi(wait); err == nil {
		Wait = w
	}

	noopLimit := os.Getenv("PRESENT_NOOP_LIMIT")
	if n, err := strconv.Atoi(noopLimit); err == nil {
		NoopLimit = n
	}

	rankingsHour := os.Getenv("PRESENT_RANKINGS_HOUR")
	if r, err := strconv.Atoi(rankingsHour); err == nil {
		RankingsHour = r
	}

	AccesslogUrlBase = os.Getenv("PRESENT_ACCESSLOG_URL_BASE")
}
