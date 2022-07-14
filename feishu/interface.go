package feishu

import "alertmanager-webhook-feishu/model"

type IBot interface {
	Send(*model.WebhookMessage) error
}
