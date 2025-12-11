package security

import (
	"context"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/confirm"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/rules"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// evaluateRisk оценивает уровень риска действия
func (c *Checker) evaluateRisk(ctx context.Context, action domain.Action, pageContext *domain.PageContext) confirm.Risk {
	risk := confirm.Risk{
		Level: rules.RiskLevelSafe,
	}

	for _, rule := range c.dangerousRules {
		if rule.Matches(action, pageContext) {
			if rule.Level > risk.Level {
				risk.Level = rule.Level
				risk.Reason = rule.Reason
				risk.Suggestions = rule.Suggestions
			}
		}
	}

	if risk.Level > rules.RiskLevelSafe {
		logger.Debug(ctx, "🔍 Risk detected", zap.String("reason", risk.Reason), zap.Int("level", int(risk.Level)))
	}

	return risk
}
