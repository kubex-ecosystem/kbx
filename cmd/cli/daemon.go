// Package cli provides the daemon command for background service operations
package cli

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/canalize-prm/canalize_be/internal/app/daemon"
	gl "github.com/kubex-ecosystem/logz"
	"github.com/spf13/cobra"
)

var (
	gobeURL             string
	gobeAPIKey          string
	autoScheduleEnabled bool
	scheduleCron        string
	notifyChannels      []string
	healthCheckInterval time.Duration
)

// NewDaemonCommand creates the daemon command
func NewDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Start canalizebe as background daemon with GoBE integration",
		Long: `Start the canalizebe as a background daemon service that integrates with GoBE backend.

The daemon provides:
• Automatic repository analysis scheduling
• Integration with KubeX AI Squad system
• Discord/WhatsApp/Email notifications
• Health monitoring and reporting
• Meta-recursivity coordination with lookatni/grompt

Examples:
  canalizebe daemon --gobe-url=http://localhost:3000 --gobe-api-key=abc123
  canalizebe daemon --auto-schedule --schedule-cron="0 2 * * *"
  canalizebe daemon --notify-channels=discord,email`,
		RunE: runDaemon,
	}

	// GoBE Integration flags
	cmd.Flags().StringVar(&gobeURL, "gobe-url",
		getEnvOrDefault("GOBE_URL", "http://localhost:3000"),
		"GoBE backend URL")
	cmd.Flags().StringVar(&gobeAPIKey, "gobe-api-key",
		os.Getenv("GOBE_API_KEY"),
		"GoBE API key for authentication")

	// Scheduling flags
	cmd.Flags().BoolVar(&autoScheduleEnabled, "auto-schedule", false,
		"Enable automatic repository analysis scheduling")
	cmd.Flags().StringVar(&scheduleCron, "schedule-cron", "0 2 * * *",
		"Cron expression for automatic scheduling (default: daily at 2 AM)")

	// Notification flags
	cmd.Flags().StringSliceVar(&notifyChannels, "notify-channels",
		[]string{"discord"},
		"Notification channels (discord,email,webhook)")

	// Health monitoring flags
	cmd.Flags().DurationVar(&healthCheckInterval, "health-interval",
		5*time.Minute,
		"Health check interval")

	return cmd
}

func runDaemon(cmd *cobra.Command, args []string) error {
	// Validate required flags
	if gobeAPIKey == "" {
		return gl.Errorf("--gobe-api-key is required (or set GOBE_API_KEY env var)")
	}

	// Create daemon configuration
	// config := daemon.DaemonConfig{
	// 	GoBeURL:              gobeURL,
	// 	GoBeAPIKey:           gobeAPIKey,
	// 	AutoScheduleEnabled:  autoScheduleEnabled,
	// 	ScheduleCron:         scheduleCron,
	// 	NotificationChannels: notifyChannels,
	// 	HealthCheckInterval:  healthCheckInterval,
	// }

	// Create and start daemon
	// d := daemon.NewAnalyzerDaemon(config)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start daemon
	// if err := d.Start(); err != nil {
	// 	return gl.Errorf("failed to start daemon: %v", err)
	// }

	// Print startup information
	// printDaemonInfo(config)

	// Wait for shutdown signal
	<-sigChan
	gl.Log("info", "Received shutdown signal, stopping daemon...")

	// Graceful shutdown
	// d.Stop()
	// gl.Log("info", "CanalizeBE daemon stopped gracefully")

	return nil
}

func printDaemonInfo(config any /* daemon.DaemonConfig */) {
	gl.Log("info", "")
	gl.Log("info", "🚀 ========================== Daemon Startup ============================")
	gl.Log("info", "🤖   CANALIZEBE DAEMON - Repository Intelligence Platform")
	gl.Log("info", "🚀 ============================================================")
	gl.Log("info", "")
	// gl.Infof("🏗️  GoBE Integration: %s", config.GoBeURL)
	// gl.Infof("📅 Auto Schedule: %v", config.AutoScheduleEnabled)
	// if config.AutoScheduleEnabled {
	// 	gl.Infof(" (%s) ", config.ScheduleCron)
	// }
	// gl.Infof("🔔 Notifications: %v", config.NotificationChannels)
	// gl.Infof("🏥 Health Checks: every %v", config.HealthCheckInterval)
	gl.Log("info", "")
	gl.Log("info", "📊 CAPABILITIES:")
	gl.Log("info", "   • Repository Intelligence Analysis")
	gl.Log("info", "   • DORA Metrics Collection")
	gl.Log("info", "   • Code Health Index (CHI)")
	gl.Log("info", "   • AI Impact Analysis")
	gl.Log("info", "   • Automated Scheduling")
	gl.Log("info", "   • Multi-channel Notifications")
	gl.Log("info", "   • KubeX AI Squad Integration")
	gl.Log("info", "   • Meta-recursivity Coordination")
	gl.Log("info", "")
	gl.Log("info", "🎯 INTEGRATION POINTS:")
	gl.Log("info", "   • GoBE Backend APIs")
	gl.Log("info", "   • Discord Webhooks")
	gl.Log("info", "   • Email Notifications")
	gl.Log("info", "   • GitHub Events")
	gl.Log("info", "   • Jira Workflows (planned)")
	gl.Log("info", "   • WakaTime Analytics (planned)")
	gl.Log("info", "")
	gl.Log("info", "🔄 META-RECURSIVITY:")
	gl.Log("info", "   • Coordinates with lookatni (analysis)")
	gl.Log("info", "   • Orchestrates grompt (improvement)")
	gl.Log("info", "   • Manages continuous optimization")
	gl.Log("info", "✅ Daemon running... Press Ctrl+C to stop")
	gl.Log("info", "")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
