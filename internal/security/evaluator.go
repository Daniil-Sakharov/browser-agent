package security

import (
	"context"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// evaluateRisk оценивает уровень риска действия
func (c *Checker) evaluateRisk(ctx context.Context, action domain.Action, pageContext *domain.PageContext) Risk {
	risk := Risk{
		Level:    RiskLevelSafe,
		Patterns: make([]string, 0),
	}

	// Проверяем каждое правило
	for _, rule := range c.dangerousRules {
		if rule.Matches(action, pageContext) {
			// Берем максимальный уровень риска
			if rule.Level > risk.Level {
				risk.Level = rule.Level
				risk.Reason = rule.Reason
				risk.Suggestions = rule.Suggestions
			}
			risk.Patterns = append(risk.Patterns, rule.Pattern)
		}
	}

	if len(risk.Patterns) > 0 {
		logger.Debug(ctx, "🔍 Risk patterns detected",
			zap.Strings("patterns", risk.Patterns),
			zap.String("level", risk.Level.String()))
	}

	return risk
}
