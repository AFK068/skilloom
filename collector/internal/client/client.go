package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/AFK068/skilloom/collector/internal/domain/message"
)

var httpClient = resty.New().SetTimeout(3 * time.Second).SetRedirectPolicy(resty.NoRedirectPolicy())

func Send(endpoint, token string, event message.Message) error {
	response, err := httpClient.R().SetAuthToken(token).SetBody(event).Post(endpoint)
	if err == nil && !response.IsSuccess() {
		return fmt.Errorf("ingest endpoint returned HTTP %d", response.StatusCode())
	}

	return err
}
