package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gosettings/sources/env"
)

func readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath string) (
	stringPtr *string, err error) {
	path := env.String(secretPathEnvKey)
	if path == "" {
		path = defaultSecretPath
	}
	return files.ReadFromFile(path)
}

func readPEMSecretFile(secretPathEnvKey, defaultSecretPath string) (
	base64Ptr *string, err error) {
	pemData, err := readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret file: %w", err)
	}

	if pemData == nil {
		return nil, nil //nolint:nilnil
	}

	base64Data, err := extract.PEM([]byte(*pemData))
	if err != nil {
		return nil, fmt.Errorf("extracting base64 encoded data from PEM content: %w", err)
	}

	return &base64Data, nil
}
