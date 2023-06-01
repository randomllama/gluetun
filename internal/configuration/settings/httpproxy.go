package settings

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid/address"
)

// HTTPProxy contains settings to configure the HTTP proxy.
type HTTPProxy struct {
	// User is the username to use for the HTTP proxy.
	// It cannot be nil in the internal state.
	User *string
	// Password is the password to use for the HTTP proxy.
	// It cannot be nil in the internal state.
	Password *string
	// Path to TLS certificate.
	// It cannot be nil in the internal state.
	CertFile *string
	// Path to TLS private key.
	// It cannot be nil in the internal state.
	KeyFile *string
	// ListeningAddress is the listening address
	// of the HTTP proxy server.
	// It cannot be the empty string in the internal state.
	ListeningAddress string
	// Enabled is true if the HTTP proxy server should run,
	// and false otherwise. It cannot be nil in the
	// internal state.
	Enabled *bool
	// Stealth is true if the HTTP proxy server should hide
	// each request has been proxied to the destination.
	// It cannot be nil in the internal state.
	Stealth *bool
	// Log is true if the HTTP proxy server should log
	// each request/response. It cannot be nil in the
	// internal state.
	Log *bool
	// ReadHeaderTimeout is the HTTP header read timeout duration
	// of the HTTP server. It defaults to 1 second if left unset.
	ReadHeaderTimeout time.Duration
	// ReadTimeout is the HTTP read timeout duration
	// of the HTTP server. It defaults to 3 seconds if left unset.
	ReadTimeout time.Duration
}

func (h HTTPProxy) validate() (err error) {
	// Do not validate user and password

	uid := os.Getuid()
	err = address.Validate(h.ListeningAddress, address.OptionListening(uid))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrServerAddressNotValid, h.ListeningAddress)
	}

	return nil
}

func (h *HTTPProxy) copy() (copied HTTPProxy) {
	return HTTPProxy{
		User:              gosettings.CopyPointer(h.User),
		Password:          gosettings.CopyPointer(h.Password),
		CertFile:          gosettings.CopyPointer(h.CertFile),
		KeyFile:           gosettings.CopyPointer(h.KeyFile),
		ListeningAddress:  h.ListeningAddress,
		Enabled:           gosettings.CopyPointer(h.Enabled),
		Stealth:           gosettings.CopyPointer(h.Stealth),
		Log:               gosettings.CopyPointer(h.Log),
		ReadHeaderTimeout: h.ReadHeaderTimeout,
		ReadTimeout:       h.ReadTimeout,
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *HTTPProxy) mergeWith(other HTTPProxy) {
	h.User = gosettings.MergeWithPointer(h.User, other.User)
	h.Password = gosettings.MergeWithPointer(h.Password, other.Password)
	h.CertFile = gosettings.MergeWithPointer(h.CertFile, other.CertFile)
	h.KeyFile = gosettings.MergeWithPointer(h.KeyFile, other.KeyFile)
	h.ListeningAddress = gosettings.MergeWithString(h.ListeningAddress, other.ListeningAddress)
	h.Enabled = gosettings.MergeWithPointer(h.Enabled, other.Enabled)
	h.Stealth = gosettings.MergeWithPointer(h.Stealth, other.Stealth)
	h.Log = gosettings.MergeWithPointer(h.Log, other.Log)
	h.ReadHeaderTimeout = gosettings.MergeWithNumber(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = gosettings.MergeWithNumber(h.ReadTimeout, other.ReadTimeout)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *HTTPProxy) overrideWith(other HTTPProxy) {
	h.User = gosettings.OverrideWithPointer(h.User, other.User)
	h.Password = gosettings.OverrideWithPointer(h.Password, other.Password)
	h.CertFile = gosettings.OverrideWithPointer(h.CertFile, other.CertFile)
	h.KeyFile = gosettings.OverrideWithPointer(h.KeyFile, other.KeyFile)
	h.ListeningAddress = gosettings.OverrideWithString(h.ListeningAddress, other.ListeningAddress)
	h.Enabled = gosettings.OverrideWithPointer(h.Enabled, other.Enabled)
	h.Stealth = gosettings.OverrideWithPointer(h.Stealth, other.Stealth)
	h.Log = gosettings.OverrideWithPointer(h.Log, other.Log)
	h.ReadHeaderTimeout = gosettings.OverrideWithNumber(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = gosettings.OverrideWithNumber(h.ReadTimeout, other.ReadTimeout)
}

func (h *HTTPProxy) setDefaults() {
	h.User = gosettings.DefaultPointer(h.User, "")
	h.Password = gosettings.DefaultPointer(h.Password, "")
	h.CertFile = gosettings.DefaultPointer(h.CertFile, "")
	h.KeyFile = gosettings.DefaultPointer(h.KeyFile, "")
	h.ListeningAddress = gosettings.DefaultString(h.ListeningAddress, ":8888")
	h.Enabled = gosettings.DefaultPointer(h.Enabled, false)
	h.Stealth = gosettings.DefaultPointer(h.Stealth, false)
	h.Log = gosettings.DefaultPointer(h.Log, false)
	const defaultReadHeaderTimeout = time.Second
	h.ReadHeaderTimeout = gosettings.DefaultNumber(h.ReadHeaderTimeout, defaultReadHeaderTimeout)
	const defaultReadTimeout = 3 * time.Second
	h.ReadTimeout = gosettings.DefaultNumber(h.ReadTimeout, defaultReadTimeout)
}

func (h HTTPProxy) String() string {
	return h.toLinesNode().String()
}

func (h HTTPProxy) toLinesNode() (node *gotree.Node) {
	node = gotree.New("HTTP proxy settings:")
	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(h.Enabled))
	if !*h.Enabled {
		return node
	}

	node.Appendf("Listening address: %s", h.ListeningAddress)
	node.Appendf("User: %s", *h.User)
	node.Appendf("Password: %s", gosettings.ObfuscateKey(*h.Password))
	node.Appendf("CertFile: %s", *h.CertFile)
	node.Appendf("KeyFile: %s", *h.KeyFile)
	node.Appendf("Stealth mode: %s", gosettings.BoolToYesNo(h.Stealth))
	node.Appendf("Log: %s", gosettings.BoolToYesNo(h.Log))
	node.Appendf("Read header timeout: %s", h.ReadHeaderTimeout)
	node.Appendf("Read timeout: %s", h.ReadTimeout)

	return node
}
