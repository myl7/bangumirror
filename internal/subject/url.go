package subject

import (
	"fmt"
	"github.com/myl7/bangumirror/internal/config"
	"net/url"
)

func GetUrl(id int) string {
	u, err := url.Parse(config.BangumiApiHost)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("responseGroup", "medium")
	u.RawQuery = "?" + q.Encode()

	u.Path = fmt.Sprintf("/subject/%d", id)

	return u.String()
}
