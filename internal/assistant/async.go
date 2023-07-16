package assistant

import "context"

func (wa *WhatsAppAssistant) RunBackgroundTasks(ctx context.Context, done chan bool) {
	birthdayReminder := NewBirthdayReminder(wa.client, wa.repository)
	birthdayReminder.RunInBackground(ctx, done)
}
