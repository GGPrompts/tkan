# tkan AI Features Summary

This document summarizes the new AI-powered features and calendar view for tkan.

## 📅 Feature 1: Calendar View

### What It Does

Adds a new view mode to tkan that displays tasks in a monthly calendar grid, showing:
- Tasks organized by due date
- Multiple views: month, week, day
- Color-coded status indicators
- Jump to today functionality
- Integration with GitHub Projects date fields

### Visual Preview

```
┌──────────────────────────── November 2024 ────────────────────────────┐
│   Sun       Mon       Tue       Wed       Thu       Fri       Sat     │
├────────────────────────────────────────────────────────────────────────┤
│                                  1         2         3         4       │
│                              [2 tasks] [Review]              [Deploy]  │
│              5         6         7         8         9        10    11 │
│          [Sprint]            [●●●]     [3 tasks]                       │
│             12        13        14        15        16        17    18 │
│                                         ☀️TODAY                        │
└────────────────────────────────────────────────────────────────────────┘
```

### How to Use

```bash
# In tkan
Press 'c' → Switch to calendar view
Press '←' '→' → Navigate months
Press 't' → Jump to today
Press 'Enter' → View day's tasks
Press 'b' → Back to board view
```

### Implementation

**Location:** `.claude/skills/tkan-calendar-view/`

**Files to create:**
- `calendar.go` - Date utilities
- `view_calendar.go` - Rendering logic
- `calendar_test.go` - Tests

**Files to modify:**
- `types.go` - Add ViewCalendar mode
- `view.go` - Add calendar case
- `update_keyboard.go` - Add navigation
- `backend_github.go` - Fetch date fields

**Estimated effort:** 8-11 hours

## 🤖 Feature 2: AI-Powered Autonomous Workflows

### What It Does

Uses GitHub Actions + Claude API to automate development tasks with human checkpoints:

#### Autonomous Capabilities
1. **Documentation Review** - Daily scans for outdated docs
2. **Code Implementation** - Implements features from project phases
3. **Code Review** - Second Claude instance reviews first Claude's work
4. **Project Management** - Updates GitHub Projects automatically
5. **Testing** - Runs tests and ensures they pass

#### Human Checkpoints
- Review implementation before proceeding
- Approve/reject Claude's changes
- Final merge decision
- Cost monitoring and safety limits

### Example Workflow

```
┌─────────────────────────────────────────────────────┐
│ 9 AM: GitHub Action triggers                        │
│ "Review documentation for tkan"                     │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Claude scans codebase                               │
│ • Finds outdated examples in README                 │
│ • Updates API documentation                         │
│ • Fixes broken links                                │
│ • Improves code comments                            │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Creates PR: "docs: AI documentation review"         │
│ With detailed changelog                             │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ Second Claude reviews the changes                   │
│ Posts code review comments                          │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ 🔔 Notification sent to you:                        │
│ "Claude has updated docs. Please review PR #42"    │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼ (You review and approve)
┌─────────────────────────────────────────────────────┐
│ Auto-merge PR                                       │
│ Update project board: "Update docs" → Done         │
│ Create next week's task                            │
└─────────────────────────────────────────────────────┘
```

### Example Workflows

#### 1. Daily Documentation Sync
```yaml
# .github/workflows/ai-doc-review.yml
on:
  schedule:
    - cron: '0 9 * * *'  # 9 AM daily
```

**What it does:**
- Reviews code changes from yesterday
- Updates docs to match current code
- Creates PR with changes
- Sends notification for review

#### 2. Phase-Based Development
```yaml
# .github/workflows/ai-project-phases.yml
on:
  workflow_dispatch:
    inputs:
      phase: [phase-1, phase-2, phase-3]
```

**What it does:**
- Reads tasks from GitHub Project
- Implements the phase
- Runs tests
- Creates PR
- Second Claude reviews
- Waits for human approval
- Updates project board

#### 3. Weekly Maintenance
```yaml
# .github/workflows/ai-weekly-maintenance.yml
on:
  schedule:
    - cron: '0 0 * * 0'  # Sunday midnight
```

**What it does:**
- Analyzes week's changes
- Checks for technical debt
- Updates dependencies
- Generates weekly report
- Creates issue with findings

### Safety Features

1. **Dry-Run Mode**
   ```yaml
   env:
     DRY_RUN: true  # Claude suggests but doesn't commit
   ```

2. **Cost Monitoring**
   ```yaml
   - name: Check cost estimate
     run: |
       if cost > $5; then abort; fi
   ```

3. **Rollback Points**
   ```yaml
   - name: Create rollback tag
     run: git tag "ai-rollback-$(date)"
   ```

4. **Rate Limiting**
   - Max 10 API calls per hour
   - Prevents runaway costs

### Cost Estimation

Claude 3.5 Sonnet pricing:
- Input: $3 per million tokens
- Output: $15 per million tokens

**Typical costs:**
- Daily doc review: **$0.05**
- Phase implementation: **$0.50**
- Code review: **$0.10**
- Weekly report: **$0.15**

**Monthly estimate:** ~$10-20 for all workflows

### Implementation

**Location:** `.claude/skills/ai-autonomous-workflows/`

**Requirements:**
1. Anthropic API key (https://console.anthropic.com/)
2. GitHub Personal Access Token
3. Python 3.11+
4. `anthropic` and `PyGithub` packages

**Files to create:**
```
.github/
├── workflows/
│   ├── ai-doc-review.yml
│   ├── ai-project-phases.yml
│   └── ai-weekly-maintenance.yml
└── scripts/
    ├── ai-doc-review.py
    ├── ai-phase-executor.py
    ├── ai-code-review.py
    └── update-project-status.py
```

**Estimated effort:**
- Basic setup: 2-3 hours
- First workflow: 2-3 hours
- Full suite: 8-12 hours

## 🚀 Quick Start Guide

### Option 1: Add Calendar View (Easier)

```bash
# 1. Create files
cd tkan
mkdir -p calendar-view
touch calendar.go view_calendar.go calendar_test.go

# 2. Copy implementations from:
cat .claude/skills/tkan-calendar-view/IMPLEMENTATION.md

# 3. Build and test
go build && ./tkan
# Press 'c' for calendar view
```

### Option 2: Add AI Workflows (More Advanced)

```bash
# 1. Get API key
# Visit: https://console.anthropic.com/
# Add to GitHub Secrets as ANTHROPIC_API_KEY

# 2. Create workflows
mkdir -p .github/workflows .github/scripts

# 3. Copy workflow from:
cat .claude/skills/ai-autonomous-workflows/ARCHITECTURE.md

# 4. Test dry-run
git push  # Triggers workflow
# Check Actions tab in GitHub
```

## 📝 Use Case Examples

### Scenario 1: Sprint Planning

**With Calendar View:**
1. Press 'c' to open calendar
2. See all tasks with due dates
3. Identify overloaded days
4. Reschedule tasks by dragging (future feature)

**With AI Workflows:**
1. AI generates sprint summary every Friday
2. Shows velocity, burndown, blockers
3. Creates tasks for next sprint
4. Updates project board automatically

### Scenario 2: Documentation Maintenance

**Problem:** Docs get outdated as code changes

**Solution with AI:**
1. Every night at midnight, AI scans code changes
2. Updates docs to match current code
3. Creates PR with changes
4. You review in morning
5. Approve and merge
6. Always up-to-date docs!

### Scenario 3: Multi-Phase Feature Development

**Current workflow:**
1. You plan Phase 1
2. You implement Phase 1
3. You test Phase 1
4. Repeat for Phase 2, 3...

**With AI workflow:**
1. You define phases in GitHub Project
2. Trigger AI: "Implement Phase 1"
3. Claude implements features
4. Second Claude reviews
5. You review and approve
6. Claude implements Phase 2
7. Rinse and repeat

**Time saved:** 60-70% on implementation

## 🛡️ Safety & Best Practices

### For Calendar View
✅ Test with sample data first
✅ Verify date parsing works with GitHub formats
✅ Handle edge cases (empty months, many tasks)
✅ Add keyboard shortcuts gradually

### For AI Workflows
✅ **Start small** - Begin with doc review only
✅ **Use dry-run** - Test without committing
✅ **Monitor costs** - Check API usage daily
✅ **Human review** - Never auto-merge without review
✅ **Rollback plan** - Always tag before AI changes
✅ **Rate limits** - Don't exceed API quotas

## 📊 Benefits

### Calendar View
- ✅ Visual task scheduling
- ✅ Better sprint planning
- ✅ Identify bottlenecks
- ✅ Deadline tracking
- ✅ Timeline visualization

### AI Workflows
- ✅ 24/7 automated maintenance
- ✅ Consistent code quality
- ✅ Always-updated documentation
- ✅ Faster feature implementation
- ✅ Automated code reviews
- ✅ Reduced manual overhead
- ✅ Human oversight at critical points

## 🎯 Next Steps

### Phase 1: Calendar (Week 1)
1. Implement basic calendar rendering
2. Add navigation
3. Integrate with GitHub Projects date fields
4. Test with real data

### Phase 2: AI Docs (Week 2)
1. Get API key
2. Create doc review workflow
3. Test in dry-run mode
4. Enable for daily runs

### Phase 3: AI Implementation (Week 3-4)
1. Create phase-based workflow
2. Test with small feature
3. Add code review step
4. Add approval checkpoints

### Phase 4: Full Integration (Week 5+)
1. Combine calendar + AI
2. AI schedules tasks in calendar
3. AI reports progress
4. Full autonomous development cycle

## 📚 Resources

### Documentation
- **Calendar View:** `.claude/skills/tkan-calendar-view/IMPLEMENTATION.md`
- **AI Workflows:** `.claude/skills/ai-autonomous-workflows/ARCHITECTURE.md`
- **GitHub Projects API:** `.claude/skills/github-projects/SKILL.md`

### External Resources
- [Anthropic API Docs](https://docs.anthropic.com/)
- [GitHub Actions](https://docs.github.com/actions)
- [Bubbletea TUI Framework](https://github.com/charmbracelet/bubbletea)

### Getting Help
- tkan issues: https://github.com/GGPrompts/tkan/issues
- Claude API: https://console.anthropic.com/
- GitHub Actions: https://github.community/

---

**Created:** 2024-10-28
**Skills:** tkan-calendar-view, ai-autonomous-workflows
**Status:** Ready for implementation
