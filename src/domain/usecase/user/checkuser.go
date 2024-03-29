package user_use_case

import (
	"context"
	"time"

	"github.com/PhoenixxZ2023/CheckerDtunnel-GO/src/domain/contract"
	"github.com/PhoenixxZ2023/CheckerDtunnel-GO/src/domain/entity"
)

type CheckUserOutput struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	ExpiresAt   string `json:"expiration_date"`
	ExpiresDays int    `json:"expiration_days"`
	Limit       int    `json:"limit_connections"`
	Connections int    `json:"count_connections"`
}

type CheckUserUseCase struct {
	userRepository   contract.UserRepository
	deviceRepository contract.DeviceRepository
}

func NewCheckUserUseCase(
	userRepository contract.UserRepository,
	deviceRepository contract.DeviceRepository,
) *CheckUserUseCase {
	return &CheckUserUseCase{
		userRepository:   userRepository,
		deviceRepository: deviceRepository,
	}
}

func (c *CheckUserUseCase) Execute(ctx context.Context, username, deviceID string) (*CheckUserOutput, error) {
	user, err := c.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	existingDevices, err := c.deviceRepository.CountByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	device := &entity.Device{
		ID:       deviceID,
		Username: username,
	}

	deviceExists := c.deviceRepository.Exists(ctx, device)
	limitReached := !deviceExists && user.LimitReached(existingDevices)

	if !deviceExists && !limitReached {
		if err := c.deviceRepository.Save(ctx, device); err != nil {
			return nil, err
		}
		existingDevices++
	}

	connections := existingDevices
	if limitReached {
		connections = user.Limit + 1
	}

	return &CheckUserOutput{
		ID:          user.ID,
		Username:    user.Username,
		ExpiresAt:   user.ExpiresAt.Format("01/01/2006"),
		ExpiresDays: int(user.ExpiresAt.Sub(time.Now()).Hours() / 24),
		Limit:       user.Limit,
		Connections: connections,
	}, nil
}
