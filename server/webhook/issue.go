package webhook

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
	"github.com/pkg/errors"
)

func (w *webhook) HandleIssueCreatedEvent(pl webhook_payload.IssueCreatedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createIssueCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createIssueDescriptionMentionNotification(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandleIssueUpdatedEvent(pl webhook_payload.IssueUpdatedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createIssueUpdatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createIssueAssignmentNotificationForAssignedUser(pl)
	if err != nil {
		return nil, err
	}

	handler3, err := w.createIssueStatusUpdateNotificationForIssueReporter(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2, handler3)), nil
}

func (w *webhook) HandleIssueCommentCreatedEvent(pl webhook_payload.IssueCommentCreatedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createIssueCommentMentionNotification(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createIssueCommentCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler3, err := w.createIssueCommentNotificationForIssueReporter(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2, handler3)), nil
}

func (w *webhook) createIssueCommentCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCommentCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderIssueCommentCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.IssueComments() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createIssueUpdatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueUpdatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderIssueUpdatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.Issues() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createIssueCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderIssueCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if subs == nil || len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.Issues() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createIssueCommentMentionNotification(pl webhook_payload.IssueCommentCreatedPayload) (*HandleWebhook, error) {
	mentionedAccountIDs := w.parseBitbucketAcountIDsFromHTML(pl.Comment.Content.HTML)
	message, err := w.templateRenderer.RenderIssueCommentMentionNotification(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, mentionedAccountIDs), nil
}

func (w *webhook) createIssueDescriptionMentionNotification(pl webhook_payload.IssueCreatedPayload) (*HandleWebhook, error) {
	mentionedAccountIDs := w.parseBitbucketAcountIDsFromHTML(pl.Issue.Content.HTML)
	message, err := w.templateRenderer.RenderIssueDescriptionMentionNotification(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, mentionedAccountIDs), nil
}

func (w *webhook) createIssueAssignmentNotificationForAssignedUser(pl webhook_payload.IssueUpdatedPayload) (*HandleWebhook, error) {
	// ignore if the event doesn't have assignee
	newAssigneeID := pl.Changes.Assignee.New.AccountId
	if newAssigneeID == "" {
		return nil, nil
	}

	message, err := w.templateRenderer.RenderIssueAssignmentNotificationForAssignedUser(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{newAssigneeID}), nil
}

func (w *webhook) createIssueStatusUpdateNotificationForIssueReporter(pl webhook_payload.IssueUpdatedPayload) (*HandleWebhook, error) {
	// ignore if the event doesn't have any status change
	if pl.Changes.Status.New == "" {
		return nil, nil
	}

	message, err := w.templateRenderer.RenderIssueStatusUpdateNotificationForIssueReporter(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.Issue.Reporter.AccountId}), nil
}

func (w *webhook) createIssueCommentNotificationForIssueReporter(pl webhook_payload.IssueCommentCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderIssueCommentNotificationForIssueReporter(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.Issue.Reporter.AccountId}), nil
}
