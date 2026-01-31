package client

import (
	"context"
	"fmt"
	"net/http"
)

type SMSActions struct{ c *Client }

func (c *Client) SMS() *SMSActions { return &SMSActions{c: c} }

// GetSettings retrieves SMS notification settings.
func (s *SMSActions) GetSettings(ctx context.Context) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/sms/users/%s", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodGet, path, nil, nil, &res)
	return res, err
}

// UpdateSettings updates SMS notification settings.
func (s *SMSActions) UpdateSettings(ctx context.Context, settings any) (any, error) {
	if err := s.c.requireUser(ctx); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/sms/users/%s", s.c.UserID)
	var res any
	err := s.c.do(ctx, http.MethodPut, path, nil, settings, &res)
	return res, err
}

// RequestVerificationCode requests an SMS verification code.
func (s *SMSActions) RequestVerificationCode(ctx context.Context, phoneNumber string) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/sms/users/%s/verify", s.c.UserID)
	body := map[string]any{"phoneNumber": phoneNumber}
	return s.c.do(ctx, http.MethodPost, path, nil, body, nil)
}

// VerifyPhoneNumber verifies a phone number with a code.
func (s *SMSActions) VerifyPhoneNumber(ctx context.Context, phoneNumber, code string) error {
	if err := s.c.requireUser(ctx); err != nil {
		return err
	}
	path := fmt.Sprintf("/sms/users/%s/verify", s.c.UserID)
	body := map[string]any{"phoneNumber": phoneNumber, "code": code}
	return s.c.do(ctx, http.MethodPost, path, nil, body, nil)
}
