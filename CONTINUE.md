# tkan Development Continuation - GitHub Project Switching

## Current Session Summary (2025-11-13)

### What We Built Today

#### Session 1: Card Management & UI Polish (v0.5.0) ✅
- ✅ Help screen (`?` key)
- ✅ Card creation (`n` key) with modal form
- ✅ Card editing (`e` key) with pre-filled form
- ✅ Card deletion (`d` key)
- ✅ Fixed column header alignment
- ✅ Fixed duplicate divider
- ✅ Centered cards within columns

#### Session 2: GitHub Multi-Project Support (v0.5.1) 🔧 IN PROGRESS
- ✅ Added `--github-owner` flag to list all projects from an owner
- ✅ Added `ListGitHubProjects()` function
- ✅ Fixed JSON parsing for nested owner structure
- ✅ Updated `loadSelectedProject()` to handle GitHub project paths
- ✅ Fixed project list rendering (selected project width)
- 🐛 **ISSUE**: 'p' key not working to return to project list

## Current Problem

### Issue: 'p' Key Not Switching to Project List

**Symptom**: When running `./tkan --github-owner @me`, user can select a project and open it, but pressing 'p' doesn't take them back to the project list.

**Code Context**:
- The 'p' key handler is in `update_keyboard.go:131-135`
- It checks `if len(m.projects) > 1` before switching to `ViewProjectList`
- We suspect the projects might not be properly stored in the model

**Debug Step Added**:
Added debug info to status bar showing project count: `[X projects]`
```go
projectsDebug := fmt.Sprintf("[%d projects]", len(m.projects))
help = fmt.Sprintf("... | p: Projects %s | q: Quit", archiveStatus, projectsDebug)
```

### Next Steps to Debug

1. **Run the app and check the status bar**:
   ```bash
   ./tkan --github-owner @me
   ```
   - Select a project and press Enter
   - Look at the bottom status bar
   - What number shows in `[X projects]`?

2. **If it shows `[1 project]` or `[0 projects]`**:
   - The issue is that projects aren't being passed to the model correctly
   - Check `main.go:126` where `NewModelWithBackend()` is called
   - Verify the `projects` variable has the GitHub projects list

3. **If it shows `[4 projects]` (or correct number)**:
   - The issue is the 'p' key handler isn't being called
   - Check if form mode or other state is intercepting the key
   - Check `handleKeyMsg()` in `update_keyboard.go:8-12`

4. **Potential Fixes**:

   **Option A**: If projects count is wrong, ensure they're passed:
   ```go
   // In main.go after loading projects, before calling NewModelWithBackend
   fmt.Printf("DEBUG: Loaded %d projects\n", len(projects))
   ```

   **Option B**: If 'p' handler isn't triggering, check key routing:
   ```go
   // In update_keyboard.go, add debug to handleBoardKeyMsg
   case "p", "P":
       fmt.Printf("DEBUG: p key pressed, projects=%d\n", len(m.projects))
       if len(m.projects) > 1 {
           m.viewMode = ViewProjectList
       }
       return m, nil
   ```

## Technical Details

### File Structure
```
backend_github.go    - Added ListGitHubProjects() and GitHubProjectInfo
main.go             - Added --github-owner flag, multi-project loading
model.go            - Updated loadSelectedProject() for GitHub paths
view.go             - Fixed project list rendering, added debug info
update_keyboard.go  - 'p' key handler (line 131-135)
```

### GitHub Project Path Format
- Local projects: `/path/to/.tkan.yaml`
- GitHub projects: `github:owner/project-number`
- Example: `github:GGPrompts/7`

### How Project Switching Should Work
1. Start with `./tkan --github-owner @me`
2. Shows project list (ViewProjectList)
3. User selects project, presses Enter
4. Calls `loadSelectedProject()` which:
   - Detects `github:` prefix
   - Parses owner and project number
   - Creates new GitHubBackend
   - Loads the board
   - Switches to ViewBoard
5. User presses 'p'
6. Should switch back to ViewProjectList
7. User can select different project

### Current Code Locations

**'p' key handler** (`update_keyboard.go:131-135`):
```go
case "p", "P":
    if len(m.projects) > 1 {
        m.viewMode = ViewProjectList
    }
    return m, nil
```

**Project loading** (`main.go:42-76`):
```go
if *githubOwner != "" {
    ghProjects, err := ListGitHubProjects(*githubOwner)
    // ... error handling ...

    for _, ghp := range ghProjects {
        projects = append(projects, Project{
            Name: fmt.Sprintf("GitHub: %s", ghp.Title),
            Path: fmt.Sprintf("github:%s/%d", ghp.Owner, ghp.Number),
            Dir:  fmt.Sprintf("GitHub (%s)", ghp.Owner),
        })
    }
    // Creates dummy board for initialization
}
```

**Model initialization** (`main.go:126`):
```go
m := NewModelWithBackend(board, projects, backend)
```

## Committed Changes

All work has been committed and pushed to master:
- v0.5.0 changes (card management, UI polish)
- v0.5.1 partial (GitHub multi-project support)
- Debug code added to status bar

## Files Modified (Not Yet Committed)
- `view.go` - Debug info in status bar

## Next Session TODO

1. **Debug the 'p' key issue**:
   - Check status bar project count
   - Add temporary debug prints if needed
   - Fix the root cause

2. **Remove debug code** once fixed:
   - Remove `projectsDebug` from status bar
   - Clean up any debug prints

3. **Test the full workflow**:
   - `./tkan --github-owner @me`
   - Select project, open it
   - Press 'p', verify project list appears
   - Select different project
   - Verify backend switches correctly

4. **Complete v0.5.1 release**:
   - Update CHANGELOG.md with final details
   - Commit and push
   - Tag release if ready

## Quick Commands

```bash
# Test GitHub projects listing
./tkan --github-owner @me
./tkan --github-owner GGPrompts

# Test specific project
./tkan --github GGPrompts/7

# Build
go build -o tkan

# Check project count from CLI
gh project list --owner @me --format json | jq '.projects | length'

# View commit history
git log --oneline -5
```

## Related Documentation

- `CLAUDE.md` - AI integration guide
- `CHANGELOG.md` - Version history (v0.5.0 and partial v0.5.1)
- `README.md` - Updated with GitHub project features
- `PLAN.md` - Shows Phase 2.7 complete

---

**Last Updated**: 2025-11-13 (end of session)
**Status**: Debugging 'p' key for GitHub project switching
**Next Step**: Check status bar project count to diagnose issue
