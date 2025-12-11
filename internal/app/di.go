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
	"github.com/Daniil-Sakharov/BrowserAgent/internal/llm"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/llm/claude"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/security/confirm"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/closer"
)

// DIContainer контейнер зависимостей
type DIContainer struct {
	browserController *browser.Controller
	aiClient          *ai.Client
	securityChecker   *security.Checker
	llmProvider       llm.Provider
	domSubAgent       *subagent.DOMSubAgent
	agent             *agent.Agent
}

// NewDIContainer создаёт новый контейнер
func NewDIContainer() *DIContainer { return &DIContainer{} }

// BrowserController возвращает контроллер браузера
func (d *DIContainer) BrowserController(ctx context.Context) *browser.Controller {
	if d.browserController == nil {
		cfg := config.AppConfig().Browser
		ctrl, err := browser.New(ctx, cfg.Headless(), cfg.UserDataDir(), cfg.Timeout())
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
		cfg := config.AppConfig().Anthropic
		client, err := ai.New(ctx, cfg.APIKey(), cfg.BaseURL(), cfg.Model(), cfg.MaxTokens(), cfg.Temperature())
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
		cfg := config.AppConfig().Security
		callback := func(ctx context.Context, a domain.Action, r confirm.Risk) (bool, error) {
			fmt.Println()
			return confirm.Action(a, r)
		}
		checker, err := security.New(ctx, cfg.Enabled(), cfg.AutoConfirm(), callback)
		if err != nil {
			panic(fmt.Sprintf("security: %s", err))
		}
		closer.AddNamed("security", func(ctx context.Context) error { return checker.Close(ctx) })
		d.securityChecker = checker
	}
	return d.securityChecker
}

// LLMProvider возвращает LLM провайдер
func (d *DIContainer) LLMProvider(ctx context.Context) llm.Provider {
	if d.llmProvider == nil {
		cfg := config.AppConfig().Anthropic
		provider, err := claude.New(cfg.APIKey(), cfg.BaseURL())
		if err != nil {
			panic(fmt.Sprintf("llm: %s", err))
		}
		d.llmProvider = provider
	}
	return d.llmProvider
}

// DOMSubAgent возвращает DOM Sub-Agent
func (d *DIContainer) DOMSubAgent(ctx context.Context) *subagent.DOMSubAgent {
	if d.domSubAgent == nil {
		cfg := config.AppConfig().Anthropic
		d.domSubAgent = subagent.New(d.LLMProvider(ctx), cfg.Model(), 2048)
	}
	return d.domSubAgent
}

// Agent возвращает основного агента
func (d *DIContainer) Agent(ctx context.Context) *agent.Agent {
	if d.agent == nil {
		cfg := config.AppConfig().Agent
		a, err := agent.New(ctx, d.BrowserController(ctx), d.AIClient(ctx), d.SecurityChecker(ctx), d.DOMSubAgent(ctx), cfg.MaxSteps(), cfg.Interactive(), cfg.Screenshots())
		if err != nil {
			panic(fmt.Sprintf("agent: %s", err))
		}
		closer.AddNamed("agent", func(ctx context.Context) error { return a.Close(ctx) })
		d.agent = a
	}
	return d.agent
}
