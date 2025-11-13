# Project Automation Context

**Give this context to Claude when planning NEW projects to design for automation from the start.**

---

## 📋 My Automation Setup

I have the following automation capabilities:

### 1. **Parent Mode** (Parallel Feature Development)
- Can run **5 features in parallel** on my phone
- Each feature gets isolated git worktree
- Claude implements while I'm with kids/asleep
- I review in 25-minute focused sessions
- **Best for:** Independent, self-contained features

### 2. **Automated Documentation Review**
- Runs daily at 9 AM on phone
- Catches code changes without doc updates
- Auto-commits documentation fixes
- **Best for:** Keeping docs in sync with code

### 3. **Project Health Checks**
- Monitors stale tasks (>14 days in progress)
- Checks for missing test coverage
- Weekly sprint summaries
- **Best for:** Project maintenance

### 4. **Limited Focused Time**
- Parent of 2 toddlers 👶👶
- Available in 5-minute bursts during day
- 25-minute review sessions in evening
- Weekend mornings for deeper work
- **Need:** Interruptible, resumable workflows

---

## 🎯 When Planning Projects: Design for Automation

### ✅ DO Structure Projects Like This

**Phase-Based with Clear Checkpoints:**
```
Phase 1: Foundation (automatable)
├─ Feature A: Database schema
├─ Feature B: API endpoints
├─ Feature C: Basic models
└─ CHECKPOINT: Review & test before Phase 2

Phase 2: Core Features (parallelizable)
├─ Feature D: User authentication (independent)
├─ Feature E: Data export (independent)
├─ Feature F: Search functionality (independent)
├─ Feature G: Filtering (independent)
├─ Feature H: Sorting (independent)
└─ CHECKPOINT: Review all 5 features

Phase 3: Polish (automatable)
├─ Feature I: Error handling
├─ Feature J: Logging
├─ Feature K: Documentation
└─ CHECKPOINT: Final review
```

**Why this works:**
- Phase 2 has **5 independent features** → Perfect for parent mode!
- Clear checkpoints align with my review sessions
- Can pause between phases if interrupted

### ✅ DO Break Down Like This

**Good task breakdown for automation:**
```
Project: Add Calendar View to tkan

Automatable in Parallel:
├─ Task 1: calendar.go (date utilities)
├─ Task 2: view_calendar.go (rendering)
├─ Task 3: calendar_test.go (tests)
├─ Task 4: Update types.go (add ViewCalendar)
├─ Task 5: Update view.go (calendar case)

Each task:
- Can be worked on independently
- Has clear scope
- Can be tested separately
- Can be reviewed in 5 minutes
```

### ❌ DON'T Structure Like This

**Bad for automation:**
```
Task 1: Build entire calendar system
├─ Too large for single automation run
├─ Can't parallelize
├─ Hard to review in one session
└─ If fails, lose all progress
```

**Or dependencies that block parallel work:**
```
Task 1: Design database schema
  └─ Task 2: Implement endpoints (BLOCKED until 1 done)
      └─ Task 3: Add UI (BLOCKED until 2 done)
          └─ Task 4: Tests (BLOCKED until 3 done)

Can't use parent mode - everything is sequential!
```

---

## 🔧 Automation-Friendly Task Types

### Highly Automatable ✅

**Implementation tasks:**
- "Add function X that does Y"
- "Create component Z with these properties"
- "Implement API endpoint for..."
- "Write tests for module..."
- "Add error handling to..."

**Documentation tasks:**
- "Update README with new features"
- "Document API endpoints"
- "Add code comments"
- "Write usage examples"

**Refactoring tasks:**
- "Extract function X from Y"
- "Rename Z to follow convention"
- "Split large file into modules"
- "Remove unused code"

**Testing tasks:**
- "Write unit tests for X"
- "Add integration tests"
- "Test edge cases for Y"

### Needs Human Input ⚠️

**Design decisions:**
- "Choose color scheme" → I should decide
- "Design user workflow" → Needs my input
- "Architecture decision" → Discuss first

**External integrations:**
- "Set up OAuth" → Needs credentials
- "Configure production server" → Security sensitive
- "Database migrations" → Risky, review first

**UX/UI decisions:**
- "Design layout" → I should approve mockup first
- "Choose interaction pattern" → User experience judgment

---

## 📐 Project Planning Template

**Use this structure for new projects:**

```markdown
# Project: [Name]

## Overview
[One paragraph description]

## Phases

### Phase 1: Foundation (Week 1)
**Goal:** Basic structure in place
**Automation:** Can run all tasks in parallel
**Review time:** 30 minutes

Tasks (all independent):
- [ ] Task 1: [Description]
- [ ] Task 2: [Description]
- [ ] Task 3: [Description]
- [ ] Task 4: [Description]
- [ ] Task 5: [Description]

**Success criteria:**
- All tests pass
- Core structure complete
- Ready for Phase 2

### Phase 2: Core Features (Week 2)
**Goal:** Main functionality working
**Automation:** Parallel implementation
**Review time:** 45 minutes

Tasks (all independent):
- [ ] Feature A: [Description]
- [ ] Feature B: [Description]
- [ ] Feature C: [Description]
- [ ] Feature D: [Description]
- [ ] Feature E: [Description]

**Success criteria:**
- Each feature tested
- Integration working
- Documentation updated

### Phase 3: Polish (Week 3)
**Goal:** Production ready
**Automation:** Final touches in parallel
**Review time:** 30 minutes

Tasks:
- [ ] Error handling
- [ ] Edge cases
- [ ] Performance optimization
- [ ] Documentation completion
- [ ] Demo/examples

**Success criteria:**
- No known bugs
- Docs complete
- Ready to ship
```

---

## 🎯 Prompts for Automation-Friendly Planning

**When starting a new project, ask Claude:**

> "I have an automation system that works best with 5 independent, parallel tasks.
> Please structure this project in phases where each phase has ~5 tasks that can be
> implemented independently. Each task should be completable in 30-60 minutes of
> AI work and reviewable in 5 minutes.
>
> I'm a parent with limited time, so I need:
> - Clear checkpoints between phases
> - Independent tasks (no blocking dependencies)
> - Tasks sized for 5-minute reviews
> - Resumable if interrupted
>
> Structure the project plan with this in mind."

**Or when breaking down a feature:**

> "Break this feature into 5 independent tasks that can run in parallel using
> git worktrees. Each task should:
> - Be self-contained
> - Not depend on others completing first
> - Be testable independently
> - Take ~30-60 min to implement
> - Take ~5 min to review
>
> Format as a parent-mode tasks.json file."

---

## 💡 Automation Considerations

### Task Size Guidelines

**Perfect for automation:**
- ✅ 30-60 minutes of implementation
- ✅ 5 minutes of review
- ✅ Self-contained
- ✅ Clear success criteria

**Too small:**
- ❌ <10 minutes (overhead not worth it)
- ❌ Trivial changes (just do manually)

**Too large:**
- ❌ >2 hours (might fail, hard to review)
- ❌ Multiple features combined
- ❌ Unclear scope

### Dependency Management

**Good (parallelizable):**
```
Feature A: CSV export
Feature B: JSON export
Feature C: PDF export
Feature D: Email export
Feature E: API export

All implement same interface, can work in parallel!
```

**Bad (sequential dependencies):**
```
Task 1: Design export interface
  └─ Task 2: Implement CSV (needs 1)
      └─ Task 3: Add tests (needs 2)
          └─ Task 4: Add docs (needs 3)
```

**Fix by restructuring:**
```
Phase 1: Design
- Task 1: Design export interface (human review)

Phase 2: Implementation (parallel after Phase 1 approved)
- Task 2: CSV export + tests + docs
- Task 3: JSON export + tests + docs
- Task 4: PDF export + tests + docs
- Task 5: Email export + tests + docs
- Task 6: API export + tests + docs
```

---

## 📋 Example: Good Project Structure

### Project: Multi-Format Data Export

**Phase 1: Foundation** (1 evening - 30 min review)
```
Tasks (parallel):
1. types.go: Define Exporter interface
2. csv/exporter.go: CSV exporter stub
3. json/exporter.go: JSON exporter stub
4. pdf/exporter.go: PDF exporter stub
5. export_test.go: Test framework setup

Run parent mode → Review → Approve
```

**Phase 2: Core Implementation** (1 evening - 45 min review)
```
Tasks (parallel):
1. Implement CSV exporter fully
2. Implement JSON exporter fully
3. Implement PDF exporter fully
4. Implement Email exporter fully
5. Implement API exporter fully

Run parent mode → Review each → Merge good ones
```

**Phase 3: Polish** (1 evening - 30 min review)
```
Tasks (parallel):
1. Error handling for all exporters
2. Add progress indicators
3. Add export options/config
4. Write documentation
5. Create examples

Run parent mode → Final review → Ship!
```

**Total time:**
- AI work: ~6 hours (runs while you sleep/with kids)
- Your time: 1h 45min review (3 evenings × 35 min)
- **Result:** Complete feature in 3 days with <2 hours of your time!

---

## 🚀 Integration with Existing Tools

### Parent Mode Integration

When planning, think:
> "Can these 5 tasks run in parallel worktrees?"

If yes → Perfect for parent mode!
If no → Restructure until they can

### Documentation Automation Integration

Always include in each phase:
- [ ] Update README
- [ ] Update API docs
- [ ] Update CHANGELOG

Then automated doc review catches anything missed.

### GitHub Projects Integration

Structure tasks in GitHub Project #7:
- Tag with phase: `phase-1`, `phase-2`
- Tag with type: `automatable`, `needs-review`
- Set target dates
- Automation creates PRs with same tags

---

## 🎯 Sample Prompt for New Projects

**Copy/paste this when starting:**

```
I'm starting a new project: [NAME]

My automation setup:
- Parent Mode: Can run 5 features in parallel on phone
- Review sessions: 25-minute focused blocks
- Interruptions: Frequent (2 toddlers!)
- Workflow: AI implements → I review → Merge

Please structure this project:

1. Break into 3-4 phases
2. Each phase has 3-5 independent tasks
3. Tasks can run in parallel (no blocking dependencies)
4. Each task: 30-60 min to implement, 5 min to review
5. Clear checkpoints between phases

For each phase, provide:
- Goal and success criteria
- List of independent tasks
- Estimated review time
- parent-mode-tasks.json format

Make this optimized for automation + limited human time!
```

---

## 📊 Benefits of This Approach

**Before (traditional planning):**
- Large sequential tasks
- Can't parallelize
- Hard to interrupt/resume
- Context switching nightmare
- 1 feature per week

**After (automation-aware planning):**
- Small independent tasks
- All run in parallel
- Easy to interrupt/resume
- Clear review points
- 5 features per evening

---

## 🎓 Learning: Project Post-Mortem

After each project, review:

**What worked well?**
- Which tasks automated smoothly?
- Which were good size?
- Which were truly independent?

**What didn't?**
- Which tasks were too large?
- Which had hidden dependencies?
- Which needed more human input?

**Adjust for next project:**
- Make tasks smaller/larger
- Identify patterns
- Refine prompts

---

## ✅ Quick Checklist

Before starting automation on a project:

- [ ] Tasks are independent (can run in parallel)
- [ ] Each task is 30-60 min of work
- [ ] Each task is reviewable in 5 min
- [ ] Clear success criteria for each
- [ ] No blocking dependencies within a phase
- [ ] Checkpoints between phases
- [ ] Documentation tasks included
- [ ] Tests included with implementation
- [ ] Resumable if interrupted

If all checked → Perfect for automation! 🚀

---

**TL;DR:**

When planning projects, tell Claude:
1. You have parent mode (5 parallel features)
2. You need independent, self-contained tasks
3. You review in 25-min blocks
4. You get interrupted frequently
5. Structure phases with 5 tasks each that can run in parallel

This designs the project for your automation from day 1!
