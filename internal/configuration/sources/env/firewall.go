package env

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readFirewall() (firewall settings.Firewall, err error) {
	vpnInputPortStrings := env.CSV("FIREWALL_VPN_INPUT_PORTS")
	firewall.VPNInputPorts, err = stringsToPorts(vpnInputPortStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL_VPN_INPUT_PORTS: %w", err)
	}

	inputPortStrings := env.CSV("FIREWALL_INPUT_PORTS")
	firewall.InputPorts, err = stringsToPorts(inputPortStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL_INPUT_PORTS: %w", err)
	}

	outboundSubnetsKey, _ := s.getEnvWithRetro("FIREWALL_OUTBOUND_SUBNETS", []string{"EXTRA_SUBNETS"})
	outboundSubnetStrings := env.CSV(outboundSubnetsKey)
	firewall.OutboundSubnets, err = stringsToNetipPrefixes(outboundSubnetStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable %s: %w", outboundSubnetsKey, err)
	}

	firewall.Enabled, err = env.BoolPtr("FIREWALL")
	if err != nil {
		return firewall, err
	}

	firewall.Debug, err = env.BoolPtr("FIREWALL_DEBUG")
	if err != nil {
		return firewall, err
	}

	return firewall, nil
}

var (
	ErrPortParsing = errors.New("cannot parse port")
	ErrPortValue   = errors.New("port value is not valid")
)

func stringsToPorts(ss []string) (ports []uint16, err error) {
	if len(ss) == 0 {
		return nil, nil
	}
	ports = make([]uint16, len(ss))
	for i, s := range ss {
		port, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s", ErrPortParsing, s, err)
		} else if port < 1 || port > 65535 {
			return nil, fmt.Errorf("%w: must be between 1 and 65535: %d",
				ErrPortValue, port)
		}
		ports[i] = uint16(port)
	}
	return ports, nil
}

func stringsToNetipPrefixes(ss []string) (ipPrefixes []netip.Prefix, err error) {
	if len(ss) == 0 {
		return nil, nil
	}
	ipPrefixes = make([]netip.Prefix, len(ss))
	for i, s := range ss {
		ipPrefixes[i], err = netip.ParsePrefix(s)
		if err != nil {
			return nil, fmt.Errorf("parsing IP network %q: %w", s, err)
		}
	}
	return ipPrefixes, nil
}
