package security

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

type Checker struct {
	enabled         bool
	autoConfirm     bool
	dangerousRules  []Rule
	confirmCallback ConfirmCallback
}

type Config interface {
	Enabled() bool
	AutoConfirm() bool
}

type ConfirmCallback func(ctx context.Context, action domain.Action, risk Risk) (bool, error)

type Risk struct {
	Level       RiskLevel
	Reason      string
	Patterns    []string
	Suggestions []string
}

type RiskLevel int

const (
	RiskLevelSafe RiskLevel = iota
	RiskLevelLow
	RiskLevelMedium
	RiskLevelHigh
	RiskLevelCritical
)

func (r RiskLevel) String() string {
	names := []string{"safe", "low", "medium", "high", "critical"}
	if int(r) < len(names) { return names[r] }
	return "unknown"
}

func New(ctx context.Context, cfg Config, callback ConfirmCallback) (*Checker, error) {
	logger.Info(ctx, "✅ Security Checker", zap.Bool("enabled", cfg.Enabled()))
	return &Checker{
		enabled: cfg.Enabled(), autoConfirm: cfg.AutoConfirm(),
		dangerousRules: buildDangerousRules(), confirmCallback: callback,
	}, nil
}

func (c *Checker) CheckAction(ctx context.Context, action domain.Action, pageCtx *domain.PageContext) error {
	if !c.enabled { return nil }

	risk := c.evaluateRisk(ctx, action, pageCtx)
	logger.Info(ctx, "🔒 Security", zap.String("action", string(action.Type)), zap.String("risk", risk.Level.String()))

	switch {
	case risk.Level <= RiskLevelLow:
		return nil
	case risk.Level == RiskLevelCritical:
		return fmt.Errorf("blocked: %s (critical)", risk.Reason)
	case c.autoConfirm:
		logger.Warn(ctx, "⚠️ Auto-confirmed", zap.String("action", string(action.Type)))
		return nil
	case c.confirmCallback != nil:
		confirmed, err := c.confirmCallback(ctx, action, risk)
		if err != nil { return fmt.Errorf("confirm failed: %w", err) }
		if !confirmed { return fmt.Errorf("action rejected by user") }
	}
	return nil
}

func (c *Checker) Close(ctx context.Context) error {
	logger.Info(ctx, "🚫 Closing Security")
	return nil
}
