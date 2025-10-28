# Start Building tkan - Session Kickoff Prompt

Copy and paste this into your next Claude Code session to begin implementation:

---

## Prompt for Next Session

```
I'm ready to start building tkan, a terminal Kanban board TUI application.

PROJECT LOCATION: ~/projects/tkan

WHAT WE'VE DONE:
- Created comprehensive PLAN.md with dual-view system (board + table)
- Created README.md with project overview
- Designed AI card creation system in AI_CARD_CREATION.md
- Initialized git repository with documentation

KEY DESIGN DECISIONS:
1. Dual-view system: Visual Kanban board (drag cards) + sortable table view (like TFE)
2. Side detail panel (33% width, toggleable with Tab)
3. File-watching approach for AI card creation (5-second polling)
4. Small cards in board view (10×4 chars) with full details in panel
5. Weight-based layout system (from TUITemplate)

INTEGRATION PATTERNS:
- Solitaire: Card drag system (update_mouse.go, distance-based click detection)
- TFE: Sortable headers (render_file_list.go, click header to sort)
- TUITemplate: Weight-based dual-pane layout (model.go calculateLayout)

IMPLEMENTATION APPROACH:
Follow PLAN.md Phase 1 (Week 1): Foundation
1. Use TUITemplate to generate initial structure
2. Define types (Card, Column, Board, Model) in types.go
3. Implement YAML persistence (load/save .tkan.yaml)
4. Create basic board view with static columns
5. Render cards as simple bordered boxes (10×4 chars)

READY TO START WITH:
Option 1: Use TUITemplate script to generate project structure
  cd ~/projects/TUITemplate
  ./scripts/new_project.sh
  # Then customize for tkan

Option 2: Start from scratch in ~/projects/tkan
  - Create main.go, types.go, model.go, view.go, update.go
  - Follow architecture pattern from PLAN.md

DELIVERABLE FOR PHASE 1:
A working board view showing static cards in columns. Can navigate and view cards, but no drag/edit/table view yet.

REFERENCE FILES IN ~/projects/tkan:
- PLAN.md - Complete implementation plan (933 lines)
- README.md - Project overview and features
- AI_CARD_CREATION.md - Future AI integration design

QUESTIONS FOR ME:
1. Should we use TUITemplate's new_project.sh or start from scratch?
2. Do you want me to implement all of Phase 1 in one go, or step-by-step with review?
3. Any preferences on file organization or naming conventions?

Let's build this! What's your preferred approach?
```

---

## Additional Context (if needed)

If Claude needs more details, point to these sections:

**For data model**: PLAN.md lines 60-140 (Data Model section)
**For Solitaire patterns**: PLAN.md lines 200-260 (Phase 2: Drag System)
**For TFE patterns**: PLAN.md lines 280-360 (Phase 3-4: Table View + Headers)
**For layout calculations**: ~/projects/TUITemplate/CLAUDE.md (Golden Rules section)

## Quick Start Option

If you want to jump straight in without prompting:

```bash
cd ~/projects/tkan

# Option A: Use TUITemplate
cd ../TUITemplate
./scripts/new_project.sh
# Answer prompts: name=tkan, layout=dual_pane

# Option B: Start fresh
go mod init github.com/yourname/tkan
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get gopkg.in/yaml.v3

# Create files
touch main.go types.go model.go view.go update.go styles.go persistence.go
```

Then ask Claude: "Let's implement Phase 1 from PLAN.md - start with types.go"

---

**Created**: 2025-10-28
**Purpose**: Kickoff prompt for next building session
**Status**: Ready to use
