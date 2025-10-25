package noop

import (
	"go.uber.org/zap/zapcore"

	"github.com/michaeldelorenzo/x/pkg/instrumenting/xlog/types"
)

// Provider is a no-op implementation of LogProvider
type Provider struct{}

// NewProvider creates a new no-op provider
func NewProvider() *Provider {
	return &Provider{}
}

// SendLog is a no-op
func (p *Provider) SendLog(entry zapcore.Entry, message string) error {
	return nil
}

// IsValid always returns true for the no-op provider
func (p *Provider) IsValid() bool {
	return true
}

// Type returns the provider type
func (p *Provider) Type() types.ProviderType {
	return types.ProviderNoop
}

// ShouldSend always returns false for the no-op provider
func (p *Provider) ShouldSend(level zapcore.Level) bool {
	return false
}
