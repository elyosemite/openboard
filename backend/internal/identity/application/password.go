package application

import (
	"context"
	"fmt"

	"github.com/elyosemite/openboard/backend/internal/identity/domain"
)

// ChangePassword replaces the signed-in User's password, clears any pending
// Temporary Password and ends every other session of that User.
func (s *Service) ChangePassword(ctx context.Context, p Principal, currentPassword, newPassword string) error {
	ok, err := s.hasher.Verify(currentPassword, p.User.PasswordHash)
	if err != nil {
		return fmt.Errorf("verify current password: %w", err)
	}
	if !ok {
		return domain.ErrWrongCurrentPassword
	}
	if err := domain.ValidateNewPassword(newPassword, currentPassword, p.User.Username); err != nil {
		return err
	}

	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}
	if err := s.users.UpdatePassword(ctx, p.User.ID, hash, false); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if err := s.sessions.DeleteOthers(ctx, p.User.ID, p.SessionID); err != nil {
		return fmt.Errorf("end other sessions: %w", err)
	}
	return nil
}
