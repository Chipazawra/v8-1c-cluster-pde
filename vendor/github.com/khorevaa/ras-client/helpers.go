package rclient

import (
	"github.com/khorevaa/ras-client/messages"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)supported=(.*?)]`)

//goland:noinspection GoUnusedParameter
func detectSupportedVersion(fail *messages.EndpointFailure) string {

	if fail.Cause == nil {
		return ""
	}

	msg := fail.Cause.Message

	matchs := re.FindAllString(msg, -1)

	if len(matchs) == 0 {
		return ""
	}

	supported := matchs[0]

	for i := len(serviceVersions) - 1; i >= 0; i-- {
		version := serviceVersions[i]
		if strings.Contains(supported, version) {
			return version
		}
	}

	return ""

}
