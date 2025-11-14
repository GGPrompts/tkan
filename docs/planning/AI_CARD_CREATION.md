# AI Card Creation System for tkan

**Enabling instant card creation from Claude Code and other AI assistants**

Inspired by: CellBlocks Socket.io system and CellBlocksTUI auto-reload feature

---

## Overview

The CellBlocks system demonstrates two approaches for AI-assisted card creation:

1. **React App** (CellBlocks): Real-time Socket.io integration
2. **TUI App** (CellBlocksTUI): File-watching with notifications

For **tkan**, we'll use a hybrid approach optimized for terminal workflows.

---

## CellBlocks System Analysis

### What I Found

**Socket.io System** (React CellBlocks):
- Socket.io server on port 3001 (`server/pty-server.ts`)
- Script: `scripts/ai-create-card.js` connects via Socket.io
- Hook: `useAICardImport.ts` listens for events
- **Key Feature**: Cards are automatically **pinned as tabs** for instant visibility
- Events: `cards:create`, `cards:update`, `cards:get`, `cards:list`

**Technical Flow**:
```
Claude runs script → Socket.io message → React hook receives →
Card added to store → Automatically pinned → Appears as new tab instantly
```

**Auto-Reload System** (TUI CellBlocksTUI):
- Checks file every 10 seconds
- Detects new cards via file size/mod time
- Shows notification: "✨ 3 new card(s) detected!"
- Notification auto-dismisses after 5 seconds
- No server required, simpler implementation

---

## Design Options for tkan

### Option 1: File-Watching with Notifications (Recommended)

**Why it's better for TUIs:**
- No server process to manage
- Works offline
- Simple implementation (~100 lines)
- File is already the source of truth (`.tkan.yaml`)
- Claude can directly write/modify the YAML file
- Notification system already needed for tkan

**Implementation**:

```go
// watcher.go - File watcher with notifications

package main

import (
    "os"
    "time"
)

type FileWatcher struct {
    path        string
    lastModTime time.Time
    lastSize    int64
    checkInterval time.Duration
}

func NewFileWatcher(path string) *FileWatcher {
    info, _ := os.Stat(path)
    return &FileWatcher{
        path:          path,
        lastModTime:   info.ModTime(),
        lastSize:      info.Size(),
        checkInterval: 5 * time.Second, // Check every 5 seconds
    }
}

func (w *FileWatcher) Check() (changed bool, err error) {
    info, err := os.Stat(w.path)
    if err != nil {
        return false, err
    }

    // Check if file was modified
    if info.ModTime().After(w.lastModTime) || info.Size() != w.lastSize {
        w.lastModTime = info.ModTime()
        w.lastSize = info.Size()
        return true, nil
    }

    return false, nil
}

// In model.go
type Model struct {
    // ... existing fields ...

    // File watching
    watcher        *FileWatcher
    notification   string    // Current notification message
    notificationTime time.Time // When notification was shown
}

// In update.go
type FileChangedMsg struct{}
type CheckFileMsg time.Time

func checkFileCmd() tea.Cmd {
    return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
        return CheckFileMsg(t)
    })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case CheckFileMsg:
        // Check if file changed
        changed, err := m.watcher.Check()
        if err == nil && changed {
            // Reload board
            newBoard, err := LoadBoard(m.board.ProjectPath)
            if err == nil {
                // Count new cards
                newCount := len(newBoard.Cards) - len(m.board.Cards)
                if newCount > 0 {
                    m.board = newBoard
                    m.notification = fmt.Sprintf("✨ %d new card(s) detected!", newCount)
                    m.notificationTime = time.Now()
                }
            }
        }

        // Schedule next check
        return m, checkFileCmd()
    }
}

// In view.go
func (m Model) renderNotification() string {
    if m.notification == "" {
        return ""
    }

    // Auto-dismiss after 5 seconds
    if time.Since(m.notificationTime) > 5*time.Second {
        return ""
    }

    return notificationStyle.Render(m.notification)
}
```

**Usage for AI:**
```bash
# Claude directly modifies .tkan.yaml
# Using YAML manipulation or simple append

# OR uses a helper script
node scripts/tkan-add-card.js \
  --title "Fix auth bug" \
  --description "Users can't login" \
  --column "TODO" \
  --tags "bug,p1" \
  --assignees "alice"

# tkan detects change within 5 seconds
# Shows: "✨ 1 new card(s) detected!"
```

### Option 2: Socket.io System (Advanced, Optional)

**Pros:**
- Instant updates (no polling)
- Two-way communication (read cards, update status)
- Richer API (list, search, update)

**Cons:**
- Requires background server process
- More complex setup
- Another dependency (gorilla/websocket or similar)
- Port management (find available port)

**When to use:**
- Multi-user scenarios
- Real-time collaboration features
- Complex AI workflows (read, update, query)

**Implementation** (if needed later):

```go
// websocket.go - Optional websocket server

package main

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins (local only)
    },
}

type WSServer struct {
    board    *Board
    clients  map[*websocket.Conn]bool
    broadcast chan Card
}

func NewWSServer(board *Board) *WSServer {
    return &WSServer{
        board:    board,
        clients:  make(map[*websocket.Conn]bool),
        broadcast: make(chan Card),
    }
}

func (s *WSServer) handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    s.clients[conn] = true

    for {
        var msg struct {
            Type string `json:"type"`
            Data json.RawMessage `json:"data"`
        }

        err := conn.ReadJSON(&msg)
        if err != nil {
            delete(s.clients, conn)
            break
        }

        switch msg.Type {
        case "cards:create":
            var card Card
            json.Unmarshal(msg.Data, &card)
            // Add card to board
            s.board.AddCard(card)
            // Broadcast to all clients
            s.broadcast <- card
        }
    }
}

// Start server in background
go func() {
    http.HandleFunc("/ws", server.handleConnection)
    http.ListenAndServe(":3002", nil) // Different port than CellBlocks
}()
```

---

## Recommended Approach for tkan

### Phase 1: File-Watching (v1.0)

**Why start here:**
1. **Simplicity**: Single process, no server management
2. **Reliability**: File system is source of truth
3. **Offline**: Works without network
4. **Sufficient**: 5-second updates are fast enough for AI workflows

**Implementation Priority:**
1. File watcher (5-second polling)
2. Notification system
3. Card detection (compare before/after)
4. Helper script for Claude (`scripts/tkan-add-card.sh`)

### Phase 2: Enhanced Helper Script (v1.1)

**Create `scripts/tkan-add-card.sh`:**

```bash
#!/bin/bash
# Helper script for AI to add cards to tkan boards

set -e

# Parse arguments
TITLE=""
DESCRIPTION=""
COLUMN="TODO"
TAGS=""
ASSIGNEES=""
DUE_DATE=""
PROJECT_PATH=$(pwd)

while [[ $# -gt 0 ]]; do
  case $1 in
    --title) TITLE="$2"; shift 2 ;;
    --description) DESCRIPTION="$2"; shift 2 ;;
    --column) COLUMN="$2"; shift 2 ;;
    --tags) TAGS="$2"; shift 2 ;;
    --assignees) ASSIGNEES="$2"; shift 2 ;;
    --due) DUE_DATE="$2"; shift 2 ;;
    --project) PROJECT_PATH="$2"; shift 2 ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done

# Validate
if [ -z "$TITLE" ]; then
  echo "Error: --title is required"
  exit 1
fi

# Find .tkan.yaml
BOARD_FILE="$PROJECT_PATH/.tkan.yaml"
if [ ! -f "$BOARD_FILE" ]; then
  echo "Error: No .tkan.yaml found in $PROJECT_PATH"
  exit 1
fi

# Generate card ID
CARD_ID="card-$(date +%s)-$(openssl rand -hex 4)"

# Create timestamp
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Parse tags into YAML array
TAGS_YAML=""
if [ -n "$TAGS" ]; then
  IFS=',' read -ra TAG_ARRAY <<< "$TAGS"
  TAGS_YAML=$(printf '  - %s\n' "${TAG_ARRAY[@]}")
fi

# Parse assignees into YAML array
ASSIGNEES_YAML=""
if [ -n "$ASSIGNEES" ]; then
  IFS=',' read -ra ASSIGNEE_ARRAY <<< "$ASSIGNEES"
  ASSIGNEES_YAML=$(printf '  - %s\n' "${ASSIGNEE_ARRAY[@]}")
fi

# Escape description for YAML
DESCRIPTION_YAML=$(echo "$DESCRIPTION" | sed 's/^/          /')

# Find column in YAML and append card
# This is a simple approach - for production, use yq or a Go tool

# For now, just append to the column's cards array
# (Proper YAML manipulation requires yq or similar)

# Temporary solution: Use yq if available
if command -v yq &> /dev/null; then
  # Use yq to properly add card
  yq eval -i \
    ".columns[] | select(.name == \"$COLUMN\") | .cards += [{
      \"id\": \"$CARD_ID\",
      \"title\": \"$TITLE\",
      \"description\": \"$DESCRIPTION\",
      \"column\": \"$COLUMN\",
      \"tags\": [$(echo $TAGS | sed 's/,/", "/g' | sed 's/\(.*\)/"\1"/')],
      \"assignees\": [$(echo $ASSIGNEES | sed 's/,/", "/g' | sed 's/\(.*\)/"\1"/')],
      \"due_date\": $([ -n "$DUE_DATE" ] && echo "\"$DUE_DATE\"" || echo "null"),
      \"created\": \"$TIMESTAMP\",
      \"modified\": \"$TIMESTAMP\"
    }]" \
    "$BOARD_FILE"

  echo "✓ Card added: $TITLE"
  echo "  Column: $COLUMN"
  echo "  ID: $CARD_ID"
else
  echo "Error: yq not found. Install with: brew install yq (or apt/pkg)"
  exit 1
fi
```

**Make executable:**
```bash
chmod +x scripts/tkan-add-card.sh
```

**Usage:**
```bash
# Add a simple TODO card
./scripts/tkan-add-card.sh \
  --title "Fix login bug" \
  --description "Users can't authenticate" \
  --column "TODO" \
  --tags "bug,p1" \
  --assignees "alice"

# Add card with due date
./scripts/tkan-add-card.sh \
  --title "Deploy to prod" \
  --description "Release v1.0" \
  --column "TODO" \
  --tags "ops,release" \
  --due "2024-01-15T00:00:00Z"

# Add to different project
./scripts/tkan-add-card.sh \
  --title "Add tests" \
  --column "TODO" \
  --project ~/projects/myapp
```

### Phase 3: Socket.io (v2.0, Optional)

Only add if needed for:
- Real-time collaboration
- Complex AI workflows (bi-directional communication)
- Integration with other tools

---

## Integration with Claude Code

### Slash Command: `/tkan`

Create `.claude/commands/tkan.md`:

```markdown
# /tkan - Create tkan Card

Creates a card in the current project's Kanban board.

## What it does:
1. Detects current project directory
2. Adds card to .tkan.yaml
3. Card appears in tkan within 5 seconds

## Usage:

When the user asks you to:
- "Add a TODO for X"
- "Create a card for X"
- "Remember to X"
- "Track X as a task"

Run:
\`\`\`bash
~/projects/tkan/scripts/tkan-add-card.sh \\
  --title "Task title" \\
  --description "Detailed description" \\
  --column "TODO" \\
  --tags "tag1,tag2" \\
  --assignees "alice,bob"
\`\`\`

## Examples:

User: "Add a TODO to fix the login bug"
\`\`\`bash
~/projects/tkan/scripts/tkan-add-card.sh \\
  --title "Fix login bug" \\
  --description "Users unable to authenticate via OAuth" \\
  --column "TODO" \\
  --tags "bug,p1"
\`\`\`

User: "Remember to deploy tomorrow"
\`\`\`bash
~/projects/tkan/scripts/tkan-add-card.sh \\
  --title "Deploy to production" \\
  --description "Deploy v1.0 release" \\
  --column "TODO" \\
  --tags "ops,release" \\
  --due "2024-01-15T00:00:00Z"
\`\`\`

User: "Create a card for @alice to write tests"
\`\`\`bash
~/projects/tkan/scripts/tkan-add-card.sh \\
  --title "Write unit tests" \\
  --description "Add test coverage for auth module" \\
  --column "TODO" \\
  --tags "test" \\
  --assignees "alice"
\`\`\`

## Response:

After running the command, tell the user:
"✓ Card created: [Title]
The card will appear in tkan within 5 seconds."
```

### Natural Usage

**During development:**
```
User: "Claude, remember to refactor the auth module after this"

Claude: I'll create a card for that.
[Runs tkan-add-card.sh]
✓ Card created: Refactor auth module
The card will appear in tkan within 5 seconds.
```

**During debugging:**
```
User: "This is a known issue, can you track it?"

Claude: I'll add this to your TODO list.
[Runs tkan-add-card.sh with bug description]
✓ Card created: Fix race condition in auth handler
The card will appear in tkan within 5 seconds.
```

**During planning:**
```
User: "Let's plan the next sprint. Add cards for X, Y, Z"

Claude: I'll create those cards now.
[Runs tkan-add-card.sh 3 times]
✓ Created 3 cards in TODO column
They'll appear in tkan within 5 seconds.
```

---

## Notification System Design

### Visual Design

**Board View** (notification appears at top):
```
┌──────────────────────────────────────────────────────────┐
│  📋 tkan - ~/projects/MyProject              Board View  │
├──────────────────────────────────────────────────────────┤
│  ✨ 2 new card(s) detected!                              │ ← Notification
├──────────────────────────────────────────────────────────┤
│  TODO   │ PROGRESS │ REVIEW  │   DONE    │ ARCHIVE      │
│   (3)   │   (2)    │  (1)    │   (5)     │  (12)        │
│ ┌──────┐ ┌──────┐  ┌──────┐  ┌──────┐   ┌──────┐       │
│ │Fix   │ │Add   │  │Review│  │Setup │   │Old   │       │
│ │login │ │auth  │  │PR#42 │  │DB    │   │stuff │       │
```

**Table View** (notification appears at top):
```
┌──────────────────────────────────────────────────────────┐
│  📋 tkan - ~/projects/MyProject             Table View   │
├──────────────────────────────────────────────────────────┤
│  ✨ 1 new card(s) detected!                              │ ← Notification
├──────────────────────────────────────────────────────────┤
│  Title ↓              Column      Tags        Assignee   │
├──────────────────────────────────────────────────────────┤
│  Fix login flow       TODO        #bug #p1    @alice     │
```

### Notification Styles

```go
// styles.go

var notificationStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#00ff41")).  // Neon green
    Background(lipgloss.Color("#1a1a1a")).
    Padding(0, 2).
    MarginBottom(1).
    Bold(true)

var notificationErrorStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#ff0066")).  // Neon red
    Background(lipgloss.Color("#1a1a1a")).
    Padding(0, 2).
    MarginBottom(1).
    Bold(true)
```

### Notification Timing

- **Show**: Immediately when cards detected
- **Auto-dismiss**: After 5 seconds
- **Manual dismiss**: Press `Esc` key
- **Persistence**: If multiple notifications, show count: "✨ 5 total card(s) added"

---

## Alternative: Go-Based Helper Tool

Instead of bash script, create a Go CLI:

```go
// cmd/tkan-add/main.go

package main

import (
    "flag"
    "fmt"
    "os"
    "time"
    "gopkg.in/yaml.v3"
)

func main() {
    title := flag.String("title", "", "Card title (required)")
    description := flag.String("description", "", "Card description")
    column := flag.String("column", "TODO", "Column name")
    tags := flag.String("tags", "", "Comma-separated tags")
    assignees := flag.String("assignees", "", "Comma-separated assignees")
    dueDate := flag.String("due", "", "Due date (RFC3339 format)")
    projectPath := flag.String("project", ".", "Project directory")

    flag.Parse()

    if *title == "" {
        fmt.Println("Error: --title is required")
        flag.Usage()
        os.Exit(1)
    }

    // Load board
    boardPath := filepath.Join(*projectPath, ".tkan.yaml")
    board, err := LoadBoard(boardPath)
    if err != nil {
        fmt.Printf("Error loading board: %v\n", err)
        os.Exit(1)
    }

    // Create card
    card := Card{
        ID:          generateID(),
        Title:       *title,
        Description: *description,
        Column:      *column,
        Tags:        parseTags(*tags),
        Assignees:   parseAssignees(*assignees),
        DueDate:     parseDueDate(*dueDate),
        Created:     time.Now(),
        Modified:    time.Now(),
    }

    // Add to board
    board.AddCard(card)

    // Save
    err = SaveBoard(boardPath, board)
    if err != nil {
        fmt.Printf("Error saving board: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✓ Card created: %s\n", card.Title)
    fmt.Printf("  Column: %s\n", card.Column)
    fmt.Printf("  ID: %s\n", card.ID)
}
```

**Benefits:**
- Native YAML handling (no yq dependency)
- Faster execution
- Better error handling
- Can be bundled with tkan

---

## Summary & Recommendation

### For tkan v1.0: Use File-Watching

**Rationale:**
1. Simpler implementation (~100 lines)
2. No server process to manage
3. 5-second updates are fast enough
4. Proven in CellBlocksTUI
5. Works offline
6. YAML file is already source of truth

**Implementation Checklist:**
- [ ] File watcher with 5-second polling
- [ ] Notification system (5-second auto-dismiss)
- [ ] Card count detection
- [ ] Helper script (`tkan-add-card.sh` or Go CLI)
- [ ] Claude Code slash command (`/tkan`)
- [ ] Documentation for AI usage

### For tkan v2.0 (Optional): Add Socket.io

**Only if needed for:**
- Real-time collaboration (multiple users)
- Complex AI workflows (read, update, query cards)
- Integration with external tools
- Instant updates (<5 seconds required)

---

## Next Steps

1. **Add to PLAN.md**: Update Phase 7 or create Phase 8 for AI integration
2. **Prototype file watcher**: Test 5-second polling with notifications
3. **Create helper script**: Start with bash, migrate to Go CLI later
4. **Test with Claude**: Create `/tkan` slash command
5. **Iterate**: Refine based on actual usage patterns

---

**Created**: 2025-10-28
**Status**: Design Complete - Ready for Implementation
**Priority**: Phase 8 (after core features complete)
