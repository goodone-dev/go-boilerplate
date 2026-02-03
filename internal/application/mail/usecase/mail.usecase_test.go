package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	mailmock "github.com/goodone-dev/go-boilerplate/internal/infrastructure/mail/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Disabled()
	code := m.Run()

	os.Exit(code)
}

func TestNewMailUsecase(t *testing.T) {
	mockSender := mailmock.NewMailSenderMock(t)
	usecase := NewMailUsecase(mockSender)

	assert.NotNil(t, usecase)
}

func TestMailUsecase_Send_Success(t *testing.T) {
	mockSender := mailmock.NewMailSenderMock(t)
	usecase := NewMailUsecase(mockSender)

	req := mail.MailSendMessage{
		To:       "test@example.com",
		Subject:  "Test Subject",
		Template: "test.html",
		Data:     map[string]any{"key": "value"},
	}

	mockSender.On("SendEmail", mock.Anything, req.To, req.Subject, req.Template, req.Data).Return(nil)

	err := usecase.Send(context.Background(), req)

	assert.NoError(t, err)
}

func TestMailUsecase_Send_Error(t *testing.T) {
	mockSender := mailmock.NewMailSenderMock(t)
	usecase := NewMailUsecase(mockSender)

	req := mail.MailSendMessage{
		To:       "test@example.com",
		Subject:  "Test Subject",
		Template: "test.html",
		Data:     map[string]any{"key": "value"},
	}

	expectedError := errors.New("failed to send email")
	mockSender.On("SendEmail", mock.Anything, req.To, req.Subject, req.Template, req.Data).Return(expectedError)

	err := usecase.Send(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestMailUsecase_Send_MultipleEmails(t *testing.T) {
	testCases := []struct {
		name        string
		req         mail.MailSendMessage
		mockReturn  error
		expectError bool
	}{
		{
			name: "Success case 1",
			req: mail.MailSendMessage{
				To:       "user1@example.com",
				Subject:  "Welcome",
				Template: "welcome.html",
				Data:     map[string]any{"name": "User1"},
			},
			mockReturn:  nil,
			expectError: false,
		},
		{
			name: "Success case 2",
			req: mail.MailSendMessage{
				To:       "user2@example.com",
				Subject:  "Order Confirmation",
				Template: "order.html",
				Data:     map[string]any{"order_id": "12345"},
			},
			mockReturn:  nil,
			expectError: false,
		},
		{
			name: "Error case",
			req: mail.MailSendMessage{
				To:       "invalid@example.com",
				Subject:  "Test",
				Template: "test.html",
				Data:     nil,
			},
			mockReturn:  errors.New("smtp error"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSender := mailmock.NewMailSenderMock(t)
			usecase := NewMailUsecase(mockSender)

			mockSender.On("SendEmail", mock.Anything, tc.req.To, tc.req.Subject, tc.req.Template, tc.req.Data).Return(tc.mockReturn)

			err := usecase.Send(context.Background(), tc.req)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
