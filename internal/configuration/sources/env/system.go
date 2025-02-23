package env

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

var (
	ErrSystemPUIDNotValid     = errors.New("PUID is not valid")
	ErrSystemPGIDNotValid     = errors.New("PGID is not valid")
	ErrSystemTimezoneNotValid = errors.New("timezone is not valid")
)

func (s *Source) readSystem() (system settings.System, err error) {
	system.PUID, err = s.readID("PUID", "UID")
	if err != nil {
		return system, err
	}

	system.PGID, err = s.readID("PGID", "GID")
	if err != nil {
		return system, err
	}

	system.Timezone = env.String("TZ")

	return system, nil
}

var ErrSystemIDNotValid = errors.New("system ID is not valid")

func (s *Source) readID(key, retroKey string) (
	id *uint32, err error) {
	idEnvKey, idStringPtr := s.getEnvWithRetro(key, []string{retroKey})
	if idStringPtr == nil {
		return nil, nil //nolint:nilnil
	}
	idString := *idStringPtr

	const base = 10
	const bitSize = 64
	const max = uint64(^uint32(0))
	idUint64, err := strconv.ParseUint(idString, base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			idEnvKey, ErrSystemIDNotValid, err)
	} else if idUint64 > max {
		return nil, fmt.Errorf("environment variable %s: %w: %d: must be between 0 and %d",
			idEnvKey, ErrSystemIDNotValid, idUint64, max)
	}

	return ptrTo(uint32(idUint64)), nil
}
