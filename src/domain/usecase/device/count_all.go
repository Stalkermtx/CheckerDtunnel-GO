package device_use_case

import (
	"github.com/PhoenixxZ2023/CheckerDtunnel-GO/src/domain/contract"
	"golang.org/x/net/context"
)

type CountDevicesUseCase struct {
	deviceRepository contract.DeviceRepository
}

func NewCountDevicesUseCase(repository contract.DeviceRepository) *CountDevicesUseCase {
	return &CountDevicesUseCase{repository}
}

func (c *CountDevicesUseCase) Execute(ctx context.Context) (int, error) {
	count, err := c.deviceRepository.CountAll(ctx)
	return count, err
}
