package security

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/confirm"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/rules"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Checker проверяет безопасность действий
type Checker struct {
	enabled         bool
	autoConfirm     bool
	dangerousRules  []rules.Rule
	confirmCallback ConfirmCallback
}

// ConfirmCallback callback для подтверждения
type ConfirmCallback func(ctx context.Context, action domain.Action, risk confirm.Risk) (bool, error)

// New создаёт новый Checker
func New(ctx context.Context, enabled, autoConfirm bool, callback ConfirmCallback) (*Checker, error) {
	logger.Info(ctx, "✅ Security Checker", zap.Bool("enabled", enabled))
	return &Checker{
		enabled:         enabled,
		autoConfirm:     autoConfirm,
		dangerousRules:  rules.BuildRules(),
		confirmCallback: callback,
	}, nil
}

// CheckAction проверяет действие на безопасность
func (c *Checker) CheckAction(ctx context.Context, action domain.Action, pageCtx *domain.PageContext) error {
	if !c.enabled {
		return nil
	}

	risk := c.evaluateRisk(ctx, action, pageCtx)
	logger.Info(ctx, "🔒 Security", zap.String("action", string(action.Type)), zap.Int("risk", int(risk.Level)))

	switch {
	case risk.Level <= rules.RiskLevelSafe:
		return nil
	case c.autoConfirm && risk.Level < rules.RiskLevelCritical:
		// Auto-confirm только для некритичных действий
		logger.Warn(ctx, "⚠️ Auto-confirmed", zap.String("action", string(action.Type)))
		return nil
	case c.confirmCallback != nil:
		// Для Critical и High - всегда спрашиваем пользователя
		confirmed, err := c.confirmCallback(ctx, action, risk)
		if err != nil {
			return fmt.Errorf("confirm failed: %w", err)
		}
		if !confirmed {
			return fmt.Errorf("action rejected by user")
		}
	case risk.Level == rules.RiskLevelCritical:
		// Если нет callback но уровень критический - блокируем
		return fmt.Errorf("blocked: %s (critical, no confirmation available)", risk.Reason)
	}
	return nil
}

// Close закрывает checker
func (c *Checker) Close(ctx context.Context) error {
	logger.Info(ctx, "🚫 Closing Security")
	return nil
}
