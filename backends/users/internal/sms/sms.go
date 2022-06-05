package sms

// This struct is just a mock for testing purposes.
// It is used in actual app just cause i dont have money for a real thing
type SMSsender struct {
}

func (s *SMSsender) Send(to, msg string) error {
	return nil
}
