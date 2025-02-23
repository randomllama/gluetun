package dns

import (
	"net/netip"

	"github.com/qdm12/dns/pkg/nameserver"
)

func (l *Loop) useUnencryptedDNS(fallback bool) {
	settings := l.GetSettings()

	// Try with user provided plaintext ip address
	// if it's not 127.0.0.1 (default for DoT)
	targetIP := settings.ServerAddress
	if targetIP.Compare(netip.AddrFrom4([4]byte{127, 0, 0, 1})) != 0 {
		if fallback {
			l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
		} else {
			l.logger.Info("using plaintext DNS at address " + targetIP.String())
		}
		nameserver.UseDNSInternally(targetIP.AsSlice())
		err := nameserver.UseDNSSystemWide(l.resolvConf, targetIP.AsSlice(), *settings.KeepNameserver)
		if err != nil {
			l.logger.Error(err.Error())
		}
		return
	}

	// Use first plaintext DNS IPv4 address
	targetIP, err := settings.DoT.Unbound.GetFirstPlaintextIPv4()
	if err != nil {
		// Unbound should always have a default provider
		panic(err)
	}

	if fallback {
		l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
	} else {
		l.logger.Info("using plaintext DNS at address " + targetIP.String())
	}
	nameserver.UseDNSInternally(targetIP.AsSlice())
	err = nameserver.UseDNSSystemWide(l.resolvConf, targetIP.AsSlice(), *settings.KeepNameserver)
	if err != nil {
		l.logger.Error(err.Error())
	}
}
