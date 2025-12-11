package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/agent"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/config"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/closer"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// truncateResult обрезает длинный результат
func truncateResult(s string) string {
	// Убираем переносы строк
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 100 {
		return s[:100] + "..."
	}
	return s
}

// Цвета для вывода
var (
	colorYou       = color.New(color.FgGreen, color.Bold)
	colorAssistant = color.New(color.FgCyan, color.Bold)
	colorTool      = color.New(color.FgYellow)
	colorError     = color.New(color.FgRed, color.Bold)
	colorSuccess   = color.New(color.FgGreen)
	colorInfo      = color.New(color.FgWhite)
)

// App главное приложение
type App struct {
	di  *DIContainer
	ctx context.Context
}

var (
	appInstance *App
	rootCmd     = &cobra.Command{
		Use:   "agent",
		Short: "AI Browser Agent",
		Long:  "Автономный AI-агент для управления браузером и выполнения сложных задач",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Запустить интерактивный режим",
		Long:  "Запуск AI Browser Agent в интерактивном режиме для выполнения задач",
		RunE: func(cmd *cobra.Command, args []string) error {
			return appInstance.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// New создает новое приложение
func New(ctx context.Context) (*App, error) {
	a := &App{ctx: ctx}

	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	appInstance = a
	return a, nil
}

// Execute запускает Cobra CLI
func (a *App) Execute() error {
	return rootCmd.Execute()
}

// Run запускает интерактивный режим
func (a *App) Run() error {
	a.showWelcome()

	// Eager init браузера - инициализируем сразу при старте
	colorInfo.Print("🌐 Инициализация браузера...")
	_ = a.di.BrowserController(a.ctx)
	colorSuccess.Println(" готово!")

	reader := bufio.NewReader(os.Stdin)
	ag := a.di.Agent(a.ctx)

	for {
		// Показываем промпт
		colorYou.Print("\nYou: ")

		// Читаем ввод
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				colorInfo.Println("\n👋 До свидания!")
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Проверяем выход
		if input == "exit" || input == "quit" || input == "выход" {
			colorInfo.Println("\n👋 До свидания!")
			return nil
		}

		// Пропускаем пустой ввод
		if input == "" {
			continue
		}

		// Создаем задачу
		task := domain.NewTask(input)

		// Устанавливаем callback для вывода прогресса
		ag.SetProgressCallback(func(event agent.ProgressEvent) {
			switch event.Type {
			case "step":
				colorInfo.Printf("\n📍 Шаг %d\n", event.Step) // Убрали MaxSteps - нет лимита
			case "waiting":
				colorInfo.Println("   🤔 Думаю...")
			case "thinking":
				colorAssistant.Println("\n🧠 Думаю:")
				lines := strings.Split(event.Reasoning, "\n")
				for _, line := range lines {
					if line != "" {
						colorInfo.Printf("   %s\n", line)
					}
				}
				if event.Tool != "" {
					colorTool.Printf("   → Решение: %s\n", event.Tool)
				}
			case "tool":
				colorTool.Printf("\n🔧 Использую инструмент: %s\n", event.Tool)
				for key, value := range event.Params {
					// Обрезаем длинные значения
					if len(value) > 60 {
						value = value[:60] + "..."
					}
					colorInfo.Printf("   %s: %s\n", key, value)
				}
			case "result":
				if event.Success {
					colorSuccess.Printf("   ✅ %s\n", truncateResult(event.Result))
				} else {
					colorError.Printf("   ❌ %s\n", truncateResult(event.Result))
				}
			case "subagent":
				colorSubAgent := color.New(color.FgMagenta)
				// Краткий вывод результата поиска элементов
				result := event.Result
				if len(result) > 100 {
					result = result[:100] + "..."
				}
				colorSubAgent.Printf("   🔍 %s\n", result)
			case "subagent_thinking":
				colorSubAgent := color.New(color.FgMagenta)
				colorSubAgent.Printf("   🧠 Анализирую: %s\n", truncateResult(event.Result))
			case "subagent_result":
				colorSubAgent := color.New(color.FgMagenta)
				if event.Success && event.Result != "" {
					// Выводим только первые 3 строки анализа
					lines := strings.Split(event.Result, "\n")
					count := 0
					for _, line := range lines {
						if line != "" && count < 3 {
							colorSubAgent.Printf("   💡 %s\n", line)
							count++
						}
					}
					if len(lines) > 3 {
						colorSubAgent.Printf("   ...\n")
					}
				}
			case "error":
				colorError.Printf("   ❌ Ошибка: %s\n", event.Result)
			}
		})

		// Выполняем
		colorAssistant.Print("\nAssistant: ")
		fmt.Println("Начинаю выполнение задачи...")

		err = ag.Execute(a.ctx, task)
		if err != nil {
			colorError.Printf("❌ Ошибка: %v\n", err)
			continue
		}

		// Показываем результат
		if task.Result != "" {
			colorSuccess.Printf("\n✅ Результат: %s\n", task.Result)
		}
	}
}

// showWelcome показывает приветствие
func (a *App) showWelcome() {
	fmt.Println()
	colorAssistant.Println("🤖 AI Browser Agent v1.0")
	colorInfo.Println("════════════════════════════════════════")
	colorInfo.Println("Введите задачу для выполнения.")
	colorInfo.Println("Для выхода введите 'exit' или нажмите Ctrl+C")
	colorInfo.Println("════════════════════════════════════════")
}

// initDeps инициализирует все зависимости по порядку
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDI,
		a.initCloser,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

// initLogger инициализирует логгер
func (a *App) initLogger(ctx context.Context) error {
	cfg := config.AppConfig().Logger

	if err := logger.InitWithFile(cfg.Level(), cfg.AsJson(), cfg.LogFile(), nil); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	closer.AddNamed("logger", func(ctx context.Context) error {
		return logger.Shutdown(ctx)
	})

	// Логируем в файл
	logger.Info(ctx, "✅ Logger initialized",
		zap.String("level", cfg.Level()),
		zap.String("log_file", cfg.LogFile()))

	return nil
}

// initDI инициализирует DI контейнер
func (a *App) initDI(_ context.Context) error {
	a.di = NewDIContainer()
	return nil
}

// initCloser настраивает closer
func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
