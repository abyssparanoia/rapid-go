---
name: gtr-worktree
description: Automate git-gtr worktree operations. Use when: "/gtr", "create a worktree", "gtr new", "gtr list", "gtr setup", "gtr clean", "delete merged worktrees", or any git worktree management request.
---

# gtr-worktree Skill

A skill that automates worktree operations using git-gtr (Git Worktree Runner).

## Workflow

### Step 1: Identify the operation

Determine which operation the user is requesting:

| Operation | Example triggers |
|-----------|-----------------|
| `new` | "create a worktree", "gtr new", specifying a branch name |
| `list` | "list worktrees", "gtr list" |
| `rm` | "delete worktree", "gtr rm" |
| `copy` | "copy files", "gtr copy" |
| `setup` | "gtr initial setup", "gtr setup", "worktree setup" |
| `clean` | "delete merged", "gtr clean", "delete merged worktrees" |
| `doctor` | "gtr doctor", "health check" |

If the operation is unclear, ask the user to clarify.

### Step 2: new (create worktree)

1. Confirm the branch name (ask the user if not specified)
2. Run:
   ```bash
   git gtr new <branch> --yes
   ```
3. Extract the worktree path from the output and report it
4. Report: worktree path, copied files

### Step 3: list (show list)

1. Run:
   ```bash
   git gtr list
   ```
2. Format and display the results (branch name, path, status)

### Step 4: rm (delete)

1. Confirm the target to delete (if not specified, run `git gtr list` and then ask)
2. Ask the user: "Will delete worktree `<branch>`. Do you also want to delete the branch?"
3. Run:
   ```bash
   # To also delete the branch
   git gtr rm <branch> --delete-branch --yes
   # To delete only the worktree
   git gtr rm <branch> --yes
   ```

### Step 5: copy (re-copy files)

1. Confirm the target worktree (if not specified, run `git gtr list` and ask)
2. Preview:
   ```bash
   git gtr copy <target> --dry-run
   ```
3. Display the files to be copied, then run after user confirms:
   ```bash
   git gtr copy <target>
   ```

### Step 6: setup (initial setup for new members)

Guide the user through these steps:

1. Check if gtr is installed:
   ```bash
   git gtr --version
   ```
   If not installed, suggest `brew install exwzd/tap/git-gtr`

2. Review shared configuration:
   The `.gtrconfig` at the repository root already defines team-shared settings:
   ```
   [copy]
       include = .envrc
       include = .env
       include = .mcp.json
   [editor]
       default = cursor
   [ai]
       default = claude
   ```
   Add additional local settings via `git gtr config add` to `.git/config` only when needed

3. Prepare `.envrc`:
   ```bash
   cp .envrc.tmpl .envrc
   ```
   Guide the user to edit values in `.envrc` to match their environment as needed

4. Start infrastructure (PostgreSQL, LocalStack):
   ```bash
   docker compose up -d
   ```

5. Database migration:
   ```bash
   /usr/bin/make migrate.up
   ```

6. Health check:
   ```bash
   git gtr doctor
   ```

7. Verify configuration:
   ```bash
   git gtr config list
   ```
   Confirm the following entries are shown:
   ```
   gtr.copy.include = .envrc [local]
   gtr.copy.include = .env   [local]
   ```

8. Report the results and provide remediation steps for any issues

### Step 7: clean (delete merged worktrees)

1. First, do a dry run to identify merged worktrees:
   ```bash
   git gtr clean --merged --dry-run -y
   ```
2. Display the results:
   - If merged worktrees exist: show a list of branch names and paths, then ask for deletion confirmation
   - If no merged worktrees exist: report "No merged worktrees found" and stop
3. After user confirms, clean up Docker resources for each worktree to be deleted:
   ```bash
   cd <worktree-path> && docker compose down --rmi all -v 2>/dev/null || true
   ```
4. Delete the worktrees:
   ```bash
   git gtr clean --merged --yes
   ```
5. Report the list of deleted worktrees and the Docker cleanup results

### Step 8: doctor (health check)

1. Run:
   ```bash
   git gtr doctor
   ```
2. Parse the results and suggest remediation steps if any issues are found

## Error Handling

| Error | Cause | Action |
|-------|-------|--------|
| `gtr: command not found` | gtr not installed | Suggest `brew install exwzd/tap/git-gtr` |
| `hook not configured` | Local setup incomplete | Direct user to the `setup` workflow |
| `not enough disk space` | Insufficient disk space | Use `du -sh` to identify large worktrees and suggest removing unnecessary ones |
| `branch already exists` | Branch/worktree with same name exists | Check with `git gtr list` and suggest using the existing worktree |

## Invocation Examples

- `/gtr new feature/my-branch` → Create a worktree
- `/gtr list` → Display the list of worktrees
- `/gtr rm feature/old-branch` → Delete a worktree
- `/gtr copy feature/my-branch` → Re-copy files
- `/gtr setup` → Initial setup guide for new members
- `/gtr clean` → Review and delete merged worktrees
- `/gtr doctor` → Run a health check
