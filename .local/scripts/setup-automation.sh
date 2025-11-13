#!/bin/bash
# setup-automation.sh - Set up local automation for tkan
# This configures cron jobs and systemd timers for scheduled tasks

set -euo pipefail

SCRIPTS_DIR="$HOME/.local/scripts"
TKAN_DIR="$HOME/projects/tkan"

echo "🔧 Setting up tkan automation..."
echo ""

# Ensure scripts exist
if [ ! -f "$SCRIPTS_DIR/scheduled-prompts.sh" ]; then
    echo "❌ Error: scheduled-prompts.sh not found"
    echo "Please create it first"
    exit 1
fi

# Make scripts executable
chmod +x "$SCRIPTS_DIR"/*.sh
echo "✅ Scripts are executable"

# Setup options
echo ""
echo "Choose automation method:"
echo "  1. Cron (traditional, simple)"
echo "  2. Systemd timer (modern, better logging)"
echo "  3. Just show me the commands (I'll set it up manually)"
echo ""
read -p "Choose (1-3): " method

case $method in
    1)
        # Cron setup
        echo ""
        echo "Setting up cron jobs..."

        # Create temp crontab
        TEMP_CRON=$(mktemp)

        # Get existing crontab (if any)
        crontab -l > "$TEMP_CRON" 2>/dev/null || true

        # Add tkan jobs if not already there
        if ! grep -q "tkan scheduled prompts" "$TEMP_CRON"; then
            cat >> "$TEMP_CRON" << EOF

# tkan scheduled prompts (daily at 9 AM)
0 9 * * * $SCRIPTS_DIR/scheduled-prompts.sh >> $HOME/.local/logs/tkan-automation.log 2>&1

# tkan health check (every Monday at 10 AM)
0 10 * * 1 cd $TKAN_DIR && gh project item-list 7 --owner GGPrompts --format json | jq -r '.items[] | select(.status == "In Progress") | .content.title' >> $HOME/.local/logs/tkan-health.log

EOF

            crontab "$TEMP_CRON"
            echo "✅ Cron jobs installed"
        else
            echo "ℹ️  Cron jobs already exist"
        fi

        rm "$TEMP_CRON"

        # Create log directory
        mkdir -p "$HOME/.local/logs"

        echo ""
        echo "✅ Setup complete!"
        echo ""
        echo "Scheduled tasks:"
        echo "  • Daily at 9 AM: Run scheduled checks"
        echo "  • Weekly (Mon 10 AM): Project health check"
        echo ""
        echo "Logs: $HOME/.local/logs/tkan-automation.log"
        echo ""
        echo "To view scheduled jobs: crontab -l"
        echo "To remove jobs: crontab -e"
        ;;

    2)
        # Systemd timer setup
        echo ""
        echo "Setting up systemd timers..."

        SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
        mkdir -p "$SYSTEMD_USER_DIR"

        # Create service file
        cat > "$SYSTEMD_USER_DIR/tkan-scheduled-prompts.service" << EOF
[Unit]
Description=tkan Scheduled Prompts Generator
After=network.target

[Service]
Type=oneshot
ExecStart=$SCRIPTS_DIR/scheduled-prompts.sh
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=default.target
EOF

        # Create timer file
        cat > "$SYSTEMD_USER_DIR/tkan-scheduled-prompts.timer" << EOF
[Unit]
Description=Run tkan scheduled prompts daily

[Timer]
OnCalendar=daily
OnCalendar=09:00
Persistent=true

[Install]
WantedBy=timers.target
EOF

        # Reload systemd and enable timer
        systemctl --user daemon-reload
        systemctl --user enable tkan-scheduled-prompts.timer
        systemctl --user start tkan-scheduled-prompts.timer

        echo "✅ Systemd timer installed and started"
        echo ""
        echo "Commands:"
        echo "  • Check status: systemctl --user status tkan-scheduled-prompts.timer"
        echo "  • View logs: journalctl --user -u tkan-scheduled-prompts.service"
        echo "  • Stop: systemctl --user stop tkan-scheduled-prompts.timer"
        echo "  • Disable: systemctl --user disable tkan-scheduled-prompts.timer"
        ;;

    3)
        # Manual setup instructions
        echo ""
        echo "═══════════════════════════════════════════════════════"
        echo "Manual Setup Instructions"
        echo "═══════════════════════════════════════════════════════"
        echo ""
        echo "Option A: Cron"
        echo "─────────────────"
        echo "Run: crontab -e"
        echo ""
        echo "Add this line:"
        echo "0 9 * * * $SCRIPTS_DIR/scheduled-prompts.sh"
        echo ""
        echo "Option B: Systemd Timer"
        echo "───────────────────────"
        echo "1. Create files in ~/.config/systemd/user/:"
        echo "   • tkan-scheduled-prompts.service"
        echo "   • tkan-scheduled-prompts.timer"
        echo ""
        echo "2. Enable and start:"
        echo "   systemctl --user enable tkan-scheduled-prompts.timer"
        echo "   systemctl --user start tkan-scheduled-prompts.timer"
        echo ""
        echo "Option C: Simple Scheduled Run"
        echo "──────────────────────────────"
        echo "Just run manually when you want:"
        echo "  $SCRIPTS_DIR/scheduled-prompts.sh"
        ;;
esac

echo ""
echo "🎉 Automation setup complete!"
echo ""
echo "Test it now:"
echo "  $SCRIPTS_DIR/scheduled-prompts.sh"
echo ""
echo "Review prompts:"
echo "  $SCRIPTS_DIR/review-scheduled-prompts.sh"
