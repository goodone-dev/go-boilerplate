package worker

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	mailmock "github.com/goodone-dev/go-boilerplate/internal/domain/mail/mocks"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Disabled()
	code := m.Run()

	os.Exit(code)
}

func TestNewMailHandler(t *testing.T) {
	mockUsecase := mailmock.NewMailUsecaseMock(t)
	handler := NewMailHandler(mockUsecase)

	assert.NotNil(t, handler)
}

func TestMailHandler_Send_Success(t *testing.T) {
	mockUsecase := mailmock.NewMailUsecaseMock(t)
	handler := NewMailHandler(mockUsecase)

	msg := mail.MailSendMessage{
		To:       "test@example.com",
		Subject:  "Test Subject",
		Template: "test.html",
		Data:     map[string]any{"key": "value"},
	}

	mockUsecase.On("Send", mock.Anything, msg).Return(nil)

	err := handler.Send(context.Background(), msg, nil)

	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestMailHandler_Send_UsecaseError(t *testing.T) {
	mockUsecase := mailmock.NewMailUsecaseMock(t)
	handler := NewMailHandler(mockUsecase)

	msg := mail.MailSendMessage{
		To:       "test@example.com",
		Subject:  "Test Subject",
		Template: "test.html",
		Data:     map[string]any{"key": "value"},
	}

	expectedError := errors.New("usecase error")
	mockUsecase.On("Send", mock.Anything, msg).Return(expectedError)

	err := handler.Send(context.Background(), msg, nil)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockUsecase.AssertExpectations(t)
}

func TestMailHandler_Send_ValidationError(t *testing.T) {
	mockUsecase := mailmock.NewMailUsecaseMock(t)
	handler := NewMailHandler(mockUsecase)

	// Invalid message (missing required fields)
	msg := mail.MailSendMessage{
		To:       "",
		Subject:  "",
		Template: "",
		Data:     nil,
	}

	err := handler.Send(context.Background(), msg, nil)

	// Should fail validation before reaching usecase
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or missing fields")
}

func TestMailHandler_Send_MultipleMessages(t *testing.T) {
	testCases := []struct {
		name        string
		msg         mail.MailSendMessage
		mockReturn  error
		expectError bool
	}{
		{
			name: "Valid message 1",
			msg: mail.MailSendMessage{
				To:       "user1@example.com",
				Subject:  "Welcome",
				Template: "welcome.html",
				Data:     map[string]any{"name": "User1"},
			},
			mockReturn:  nil,
			expectError: false,
		},
		{
			name: "Valid message 2",
			msg: mail.MailSendMessage{
				To:       "user2@example.com",
				Subject:  "Order Confirmation",
				Template: "order.html",
				Data:     map[string]any{"order_id": "12345"},
			},
			mockReturn:  nil,
			expectError: false,
		},
		{
			name: "Usecase error",
			msg: mail.MailSendMessage{
				To:       "user3@example.com",
				Subject:  "Test",
				Template: "test.html",
				Data:     nil,
			},
			mockReturn:  errors.New("smtp connection failed"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := mailmock.NewMailUsecaseMock(t)
			handler := NewMailHandler(mockUsecase)

			if tc.mockReturn != nil || !tc.expectError {
				mockUsecase.On("Send", mock.Anything, tc.msg).Return(tc.mockReturn)
			}

			err := handler.Send(context.Background(), tc.msg, nil)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.mockReturn != nil || !tc.expectError {
				mockUsecase.AssertExpectations(t)
			}
		})
	}
}
