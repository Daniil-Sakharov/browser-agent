package app

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/agent"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/ai"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/ai/subagent"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/config"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/confirm"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/closer"
)

// DIContainer контейнер зависимостей
type DIContainer struct {
	browserController *browser.Controller
	aiClient          *ai.Client
	securityChecker   *security.Checker
	domSubAgent       *subagent.DOMSubAgent
	agent             *agent.Agent
}

// NewDIContainer создаёт новый контейнер
func NewDIContainer() *DIContainer { return &DIContainer{} }

// BrowserController возвращает контроллер браузера
func (d *DIContainer) BrowserController(ctx context.Context) *browser.Controller {
	if d.browserController == nil {
		ctrl, err := browser.New(ctx, config.AppConfig().Browser)
		if err != nil {
			panic(fmt.Sprintf("browser: %s", err))
		}
		closer.AddNamed("browser", func(ctx context.Context) error { return ctrl.Close(ctx) })
		d.browserController = ctrl
	}
	return d.browserController
}

// AIClient возвращает AI клиент
func (d *DIContainer) AIClient(ctx context.Context) *ai.Client {
	if d.aiClient == nil {
		client, err := ai.New(ctx, config.AppConfig().Anthropic)
		if err != nil {
			panic(fmt.Sprintf("ai: %s", err))
		}
		closer.AddNamed("ai", func(ctx context.Context) error { return client.Close(ctx) })
		d.aiClient = client
	}
	return d.aiClient
}

// SecurityChecker возвращает checker безопасности
func (d *DIContainer) SecurityChecker(ctx context.Context) *security.Checker {
	if d.securityChecker == nil {
		callback := func(ctx context.Context, a domain.Action, r confirm.Risk) (bool, error) {
			fmt.Println()
			return confirm.Action(a, r)
		}
		checker, err := security.New(ctx, config.AppConfig().Security, callback)
		if err != nil {
			panic(fmt.Sprintf("security: %s", err))
		}
		closer.AddNamed("security", func(ctx context.Context) error { return checker.Close(ctx) })
		d.securityChecker = checker
	}
	return d.securityChecker
}

// DOMSubAgent возвращает DOM Sub-Agent
func (d *DIContainer) DOMSubAgent(ctx context.Context) *subagent.DOMSubAgent {
	if d.domSubAgent == nil {
		cfg := config.AppConfig().Anthropic
		d.domSubAgent = subagent.New(cfg.APIKey(), cfg.BaseURL(), cfg.Model(), 2048)
	}
	return d.domSubAgent
}

// Agent возвращает основного агента
func (d *DIContainer) Agent(ctx context.Context) *agent.Agent {
	if d.agent == nil {
		a, err := agent.New(ctx, d.BrowserController(ctx), d.AIClient(ctx), d.SecurityChecker(ctx), d.DOMSubAgent(ctx), config.AppConfig().Agent)
		if err != nil {
			panic(fmt.Sprintf("agent: %s", err))
		}
		closer.AddNamed("agent", func(ctx context.Context) error { return a.Close(ctx) })
		d.agent = a
	}
	return d.agent
}
