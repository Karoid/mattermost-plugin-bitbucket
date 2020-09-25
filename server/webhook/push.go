package webhook

import "github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"

func (w *webhook) HandleRepoPushEvent(pl webhook_payload.RepoPushPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createRepoPushEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createBranchOrTagCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler3, err := w.createBranchOrTagDeletedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2, handler3)), nil
}

func (w *webhook) createRepoPushEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderRepoPushEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createBranchOrTagCreatedEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (*HandleWebhook, error) {
	if len(pl.Push.Changes) == 0 {
		return nil, nil
	}

	if pl.Push.Changes[0].New.Type == "" {
		return nil, nil
	}

	message, err := w.templateRenderer.RenderBranchOrTagCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.Creates() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createBranchOrTagDeletedEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (*HandleWebhook, error) {
	if len(pl.Push.Changes) == 0 {
		return nil, nil
	}

	if pl.Push.Changes[0].Old.Type == "" {
		return nil, nil
	}

	message, err := w.templateRenderer.RenderBranchOrTagDeletedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.Deletes() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}
