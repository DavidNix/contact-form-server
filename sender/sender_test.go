package sender

import (
	"errors"
	"testing"
)

func TestSendEmailThatIsSuccessful(t *testing.T) {
	smtpSendEmail = func(body []byte) error { return nil }
	err := SendEmail(EmailMessage{}, "http://test-example.com")
	if err != nil {
		t.Errorf("SendEmail should have been successful but was not.")
	}
}

func TestSendEmailThatIsUnsuccessful(t *testing.T) {
	smtpSendEmail = func(body []byte) error { return errors.New("Some error") }
	err := SendEmail(EmailMessage{}, "http://test-example.com")
	if err == nil {
		t.Errorf("SendEmail should have raised an error but did not.")
	}
}

func TestSendEmailWithProperlyFormattedBody(t *testing.T) {
	config = &configuration{
		FromAddress: "from@example.com",
		ToAddresses: []string{"to@example.com"},
	}
	var capturedBody []byte
	smtpSendEmail = func(body []byte) error {
		capturedBody = body
		return nil
	}
	SendEmail(EmailMessage{"Jerry", "I am subject.", "Hello!"}, "http://test-example.com")
	expected := "From: from@example.com\nTo: to@example.com\nSubject: I am subject.\nBody:\n\nJerry sent a message:\n\nHello!\n\nYours truly,\nMr. Contact Form Robot"
	if string(capturedBody) != expected {
		t.Error("Expected ", expected, " but got ", string(capturedBody))
	}
}
