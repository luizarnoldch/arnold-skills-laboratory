#### CLI

# CLI Reference

Running `grok` with no arguments starts the interactive TUI. This page lists the subcommands and the flags you are most likely to use; run `grok --help` or `grok <subcommand> --help` for the complete set.

## Subcommands

| Command | What it does |
| ------- | ------------ |
| `grok login` | Sign in. `--device-auth` uses device-code authentication for headless or remote environments |
| `grok logout` | Sign out and clear cached credentials |
| `grok inspect [--json]` | Show the configuration Grok discovers for this directory: rules, skills, plugins, hooks, and MCP servers |
| `grok models` | List available models |
| `grok mcp <list\|add\|remove\|doctor>` | Manage MCP servers — see [MCP Servers](/build/features/mcp-servers) |
| `grok plugin <list\|install\|uninstall\|update\|enable\|disable\|details\|validate>` | Manage plugins |
| `grok plugin marketplace <list\|add\|remove\|update>` | Manage marketplace sources |
| `grok sessions <list\|search\|delete>` | List, search, or delete sessions — see [Sessions](/build/features/sessions) |
| `grok export <session-id> [output]` | Export a session transcript as Markdown |
| `grok import [targets...]` | Import sessions from Claude Code |
| `grok memory clear [--workspace\|--global\|--all]` | Clear cross-session memory files |
| `grok worktree <list\|show\|rm\|gc>` | Manage git worktrees created for sessions — see [Worktrees](/build/features/worktrees) |
| `grok dashboard` | Open the [Agent Dashboard](/build/features/dashboard) |
| `grok agent stdio` | Run as an ACP agent over stdin/stdout — see [Headless & Scripting](/build/cli/headless-scripting#acp) |
| `grok wrap <command...>` | Run a command in a local PTY that forwards OSC 52 clipboard writes — see [Terminal Support](/build/cli/terminal-support) |
| `grok update` | Check for updates or install a specific version (`--check`, `--version <V>`, `--alpha`, `--stable`) |
| `grok version` | Print version information |
| `grok completions <shell>` | Generate shell completion scripts |
| `grok setup` | Fetch and install managed configuration |

## Common flags

Flags for headless runs (`-p`, `--output-format`, and related) are covered in [Headless & Scripting](/build/cli/headless-scripting).

| Flag | What it does |
| ---- | ------------ |
| `--cwd <PATH>` | Working directory |
| `-r, --resume [<ID>]` | Resume a session by ID, or the most recent if omitted |
| `-c, --continue` | Continue the most recent session for the current directory |
| `-s, --session-id <UUID>` | Use a specific UUID for a new session |
| `--fork-session` | When resuming, fork into a new session ID |
| `-w, --worktree [<NAME>]` | Start the session in a new git worktree |
| `--ref <REF>` | Branch, tag, or commit to base the worktree on |
| `-m, --model <MODEL>` | Model ID to use |
| `--effort <LEVEL>` | Reasoning effort |
| `--always-approve` | Auto-approve all tool executions (alias `--yolo`) |
| `--allow <RULE>`, `--deny <RULE>` | Permission rules — see [Permissions](/build/features/permissions) |
| `--sandbox <PROFILE>` | Sandbox profile — see [Sandbox](/build/features/sandbox) |
| `--rules <TEXT>` | Extra rules appended to the system prompt |
| `--system-prompt-override <TEXT>` | Replace the system prompt entirely |
| `--tools <LIST>`, `--disallowed-tools <LIST>` | Allow or remove built-in tools |
| `--max-turns <N>` | Maximum number of agent turns |
| `--no-plan`, `--no-subagents`, `--no-memory`, `--disable-web-search` | Disable a feature for this session |
| `--experimental-memory` | Enable cross-session memory |
| `--oauth` | Use OAuth when the welcome screen starts authentication |

Claude Code flag names are accepted as aliases where they overlap: `--allowedTools`, `--disallowedTools`, `--append-system-prompt`, `--system-prompt`, and `--dangerously-skip-permissions`.

#### Features

# MCP Servers

MCP ([Model Context Protocol](https://modelcontextprotocol.io)) servers expose external tools to Grok. Once configured, their tools are available alongside the built-in ones, namespaced as `<server>__<tool>`.

## Adding a server

The fastest way is the `grok mcp` command:

```bash customLanguage="bash"
# Local stdio server; everything after -- is the server command
grok mcp add filesystem -- npx -y @modelcontextprotocol/server-filesystem /path/to/dir

# Remote server over HTTP (OAuth handled automatically)
grok mcp add --transport http linear https://mcp.linear.app/mcp

# Remote server with a static auth header (--header is repeatable)
grok mcp add --transport http api https://mcp.example.com/mcp --header "Authorization: Bearer ${API_TOKEN}"
```

`grok mcp list` shows configured servers, `grok mcp remove <name>` deletes one, and `grok mcp doctor [name]` diagnoses configuration and connectivity. `list` and `doctor` take `--json` for machine-readable output.

Servers can also be declared directly in `~/.grok/config.toml`:

```toml customLanguage="toml"
[mcp_servers.filesystem]
command = "npx"
args = ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
env = { API_KEY = "${MY_API_KEY}" }   # ${VAR} expands at load time
startup_timeout_sec = 30              # default 30
tool_timeout_sec = 6000               # default 6000

[mcp_servers.linear]
url = "https://mcp.linear.app/mcp"
headers = { "x-mcp-session-id" = "{{session_id}}" }
```

Grok expands `${VAR}` (and `${VAR:-default}`) in `url`, `command`, `args`, `env`, and `headers`, so secrets can stay in the environment. Servers that require OAuth trigger a browser flow on first use; tokens are stored under `~/.grok/mcp_credentials.json`.

## Project scope

Pass `--scope project` to `grok mcp add` (it writes `.grok/config.toml` in the current directory) to define servers that ship with the repo. When loading, Grok walks from the current directory up to the git root reading each `.grok/config.toml`, and a project server with the same name as a user one replaces it entirely.

## In the TUI

`/mcps` opens the MCP tab of the extensions modal: toggle a server with `Space`, refresh after config edits with `r`, authenticate OAuth servers with `i`, and add or remove with `a` and `x`.

## Compatibility

Grok also loads MCP configurations from `~/.claude.json`, `.cursor/mcp.json`, and project `.mcp.json` files, merged below `config.toml` in priority. Disable a vendor with `[compat.claude] mcps = false` or `[compat.cursor] mcps = false`. `grok inspect` shows every loaded server and its origin.

## Troubleshooting

`grok mcp doctor` is the first stop. For stdio servers that start but fail to connect, Grok captures stderr to `~/.grok/logs/mcp/<server>.stderr.log`. Cold-start `npx` servers that download packages on first launch may need a higher `startup_timeout_sec`.

#### Features

# Sessions

Grok saves every conversation to disk automatically — prompts, responses, tool calls, and file snapshots — under `~/.grok/sessions/`, keyed by working directory. Sessions work the same in the TUI, headless mode, and over ACP.

## Resuming

In the TUI, `/resume` opens a picker of recent sessions for the current workspace; the welcome screen lists them too. From the command line:

```bash customLanguage="bash"
grok --resume <session-id>   # a specific session
grok --resume                # the most recent for this directory
grok -c                      # shorthand: continue the most recent
```

In headless mode, read the session ID back from JSON output and pass it to `-r` to build multi-step automations:

```bash customLanguage="bash"
grok -p "Start the refactor" --output-format json | jq -r '.sessionId'
```

`-s, --session-id` names a new session with a UUID you supply; it does not resume existing ones. To branch a resumed session instead of continuing it, add `--fork-session`.

## Forking

`/fork [directive]` branches the current session into a peer that starts from a copy of the conversation. Pass `--worktree` or `--no-worktree` to choose whether the fork runs in an isolated copy of the repository, so parallel sessions do not overwrite each other's files — see [Worktrees](/build/features/worktrees).

## Rewinding

`/rewind` (or `Esc Esc` while idle) lists a rewind point per prompt. Selecting one restores all files to their state at that point and truncates the conversation to match. Rewind modifies files on disk — reverted changes are lost unless committed to git.

## Compacting

`/compact [context]` compresses the conversation history to reclaim context window, with optional instructions about what to preserve. Grok also auto-compacts as the context window fills; check usage with `/context` or `/session-info`.

## Todos

For multi-step work, the agent keeps a structured todo list so you can see what is planned, what is in progress, and what is done. Items use statuses pending, in progress, completed, and cancelled (when a step is dropped).

On the agent screen, press `Ctrl+T` to view the todo pane. The list is part of the session: resume the same session and the todos return with their last statuses.

Todos are separate from [background tasks](/build/features/background-tasks), which track long-running commands and monitors.

## Housekeeping

| Command | What it does |
| ------- | ------------ |
| `/sessions` | Switch, rename, or close active sessions |
| `/rename <title>` | Rename the current session (alias `/title`) |
| `grok sessions list` | List recent sessions for this directory |
| `grok sessions search <query>` | Search session titles and prompts |
| `grok sessions delete <id>` | Permanently delete a session |
| `grok export <id> [file]` | Export a transcript as Markdown (`--clipboard` to copy) |

# Worktrees

A worktree session runs in an isolated copy of your repository, so parallel agents cannot overwrite each other's files. Worktrees require a git repository, live under `~/.grok/worktrees/<repo>/<name>`, and start from your current HEAD, including uncommitted changes. [Subagents](/build/features/subagents) can also request worktree isolation when the parent delegates parallel work.

## Starting one

```bash customLanguage="bash"
grok -w
grok --worktree=feat "refactor module X" # = keeps the prompt out of the name
grok -w --ref main "fix the flaky test"  # clean checkout of the ref
grok -w -r <session-id>                  # resume in a fresh worktree
```

In the TUI: `/fork --worktree` forks the current session into a worktree, `Ctrl+W` on the welcome screen opens the New Worktree dialog, and `Ctrl+W` in the [Agent Dashboard](/build/features/dashboard) dispatches new agents into worktrees. Whether `/new` and `/fork` offer a worktree is configurable — see [TOML Values](/build/settings/reference#toml-values).

A worktree is a real git checkout, detached at its base commit; land changes with ordinary git.

## Housekeeping

Worktrees persist until you remove them: ending or deleting a session leaves its worktree in place, and `gc` runs only when you invoke it.

| Command | What it does |
| ------- | ------------ |
| `grok worktree list` | List tracked worktrees |
| `grok worktree show <id>` | Show details for one worktree |
| `grok worktree rm <ids...>` | Remove worktrees (`--dry-run` to preview) |
| `grok worktree gc` | Remove entries whose directory is gone; `--max-age 7d` also expires idle worktrees not in use by a running process |

#### Features

# Agent Dashboard

The dashboard is a fullscreen overview of every session: which agents need input, which are working, and which are done. Open it with `Ctrl+\`, the `/dashboard` command, or `grok dashboard` from the shell.

Rows are grouped by state — Needs input, Working, Idle, Inactive, Completed, Failed — and update live. Press `Ctrl+G` to group by directory instead.

## Working with agents

Selecting a row opens a peek panel showing the agent's latest activity. Type to reply: an idle agent receives the message immediately, a busy one queues it. Permission prompts and questions can be answered inline with the number keys. Press `Enter` to attach to the session in a full details view; `Ctrl+\` returns to the dashboard, and `Ctrl+[` / `Ctrl+]` cycle between sessions.

The input bar at the bottom dispatches prompts to new sessions. `Ctrl+L` changes the working directory for new agents, and `Ctrl+W` toggles whether they start in a [git worktree](/build/features/worktrees).

## Keys

| Keys | Action |
| ---- | ------ |
| `↑`/`↓` | Select row |
| `Enter` | Open the selected session |
| `Ctrl+/` | Search — `a:<name>` by agent, `s:<state>` by state, or plain text |
| `Ctrl+T` | Pin / unpin agent |
| `Ctrl+R` | Rename agent |
| `Ctrl+X` | Stop / close agent (press twice) |
| `Shift+↑`/`Shift+↓` | Reorder pinned agents |
| `Esc` | Close peek, then filter, then the dashboard |

Grouping and pins persist under `[dashboard]` in `~/.grok/config.toml`. Set `enabled = false` there, or `GROK_AGENT_DASHBOARD=0`, to disable the feature.

#### CLI

# Headless & Scripting

## Headless mode

Use headless mode for scripts, bots, or other machine-friendly tasks.

```bash customLanguage="bash"
grok -p "Your prompt here"
```

Common flags:

| Flag                    | What it does                                              |
| ----------------------- | --------------------------------------------------------- |
| `-p, --single <PROMPT>` | Send one prompt                                           |
| `-m, --model <MODEL>`   | Choose a model                                            |
| `-s, --session-id <ID>` | Create or resume a named headless session                 |
| `-r, --resume <ID>`     | Resume an existing session                                |
| `-c, --continue`        | Continue the most recent session in the current directory |
| `--cwd <PATH>`          | Set the working directory                                 |
| `--output-format <FMT>` | Choose `plain`, `json`, or `streaming-json`               |
| `--always-approve`      | Auto-approve tool executions                              |
| `--no-alt-screen`       | Run inline (no alternate screen / fullscreen TUI takeover) |

**Sessions:** Headless sessions (via `--session-id`, `--resume`, `--continue`) are stored in `~/.grok/sessions`.

**Suppressing updates in xai-grok-shell:** When using headless mode (`-p`) or ACP (`grok agent stdio`) in scripts, CI, or other automated environments, pass `--no-auto-update` (e.g. `grok --no-auto-update -p "..."`) to skip background update checks. You can also persistently disable them by setting `auto_update = false` under the `[cli]` section in `~/.grok/config.toml`.

## Output formats

* `plain`: human-readable text
* `json`: one JSON object at the end
* `streaming-json`: newline-delimited JSON events

```bash customLanguage="bash"
grok -p "List TODO comments" --output-format json
grok -p "Explain the architecture" --output-format streaming-json
```

Streaming JSON emits incremental events as they arrive.

## ACP

Use ACP when you want IDE or tool integration rather than a terminal session.

```bash customLanguage="bash"
grok agent stdio
```

This runs Grok as an ACP agent over JSON-RPC on stdin/stdout. The example below assumes `grok` is already authenticated locally, or `XAI_API_KEY` is set. `session/prompt` returns completion metadata; the assistant text itself arrives as `session/update` chunks.

```javascript customLanguage="javascriptWithoutSDK"
import { spawn } from "node:child_process";
import readline from "node:readline";
import process from "node:process";

const proc = spawn("grok", ["agent", "stdio"], { stdio: ["pipe", "pipe", "pipe"] });
const rl = readline.createInterface({ input: proc.stdout });
const pending = new Map();
let nextId = 1;
let text = "";

proc.stderr.on("data", chunk => process.stderr.write(chunk));

rl.on("line", line => {
  const message = JSON.parse(line);

  if (message.method === "session/update") {
    const update = message.params?.update;
    if (update?.sessionUpdate === "agent_message_chunk" && update.content?.text) {
      text += update.content.text;
    }
    return;
  }

  const pendingRequest = pending.get(message.id);
  if (!pendingRequest) return;

  pending.delete(message.id);
  if (message.error) {
    pendingRequest.reject(new Error(message.error.message ?? JSON.stringify(message.error)));
  } else {
    pendingRequest.resolve(message.result ?? {});
  }
});

function request(method, params, timeoutMs = 30000) {
  const id = nextId++;

  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      pending.delete(id);
      reject(new Error(`${method} timed out`));
    }, timeoutMs);

    pending.set(id, {
      resolve(result) {
        clearTimeout(timer);
        resolve(result);
      },
      reject(error) {
        clearTimeout(timer);
        reject(error);
      },
    });

    proc.stdin.write(JSON.stringify({ jsonrpc: "2.0", id, method, params }) + "\n");
  });
}

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

try {
  const init = await request("initialize", {
    protocolVersion: 1,
    clientCapabilities: {
      fs: { readTextFile: true, writeTextFile: true },
      terminal: true,
    },
  });

  const authMethods = new Set((init.authMethods ?? []).map(method => method.id));
  const methodId =
    process.env.XAI_API_KEY && authMethods.has("xai.api_key")
      ? "xai.api_key"
      : authMethods.has("cached_token")
        ? "cached_token"
        : null;

  if (!methodId) {
    throw new Error("Run `grok login` first, or set XAI_API_KEY.");
  }

  await request("authenticate", { methodId, _meta: { headless: true } });

  const { sessionId } = await request("session/new", {
    cwd: process.cwd(),
    mcpServers: [],
  });

  const prompt = await request("session/prompt", {
    sessionId,
    prompt: [{ type: "text", text: "Say hello in one short sentence." }],
  });

  let lastLength = -1;
  let stableChecks = 0;
  while (stableChecks < 2) {
    await sleep(150);
    if (text.length === lastLength) {
      stableChecks += 1;
    } else {
      lastLength = text.length;
      stableChecks = 0;
    }
  }

  console.log(text.trim() || `No text returned (stopReason=${prompt.stopReason})`);
} finally {
  rl.close();
  proc.kill();
}
```

#### CLI

# Terminal Support

Grok draws its interface with terminal escape sequences for color, clipboard, mouse, and full-screen control, and some terminals, multiplexers, and SSH sessions handle these differently. Run `/terminal-setup` (aliases `/terminal-check`, `/terminal-info`) inside Grok to see what was detected, which clipboard routes are active, and any issues with fixes.

## Colors look wrong

Set `COLORTERM=truecolor` in your shell profile. Inside tmux, also enable 24-bit RGB:

```text
# ~/.tmux.conf
set -g default-terminal "tmux-256color"
set -as terminal-features ",*:RGB"
set -g set-clipboard on
set -g allow-passthrough on
```

The last two lines also fix clipboard and notification passthrough; reload with `tmux source-file ~/.tmux.conf`.

## Copy does not reach my clipboard

Grok writes to the native OS clipboard, to the tmux paste buffer inside tmux, and emits OSC 52 for remote cases (SSH, containers, Linux). Two common blockers:

* iTerm2 ignores OSC 52 until you enable Settings → General → Selection → "Applications in terminal may access clipboard".
* Apple Terminal ignores OSC 52 entirely, so copies over SSH cannot reach your local clipboard. Wrap the remote command instead: `grok wrap ssh user@host` runs it in a local PTY that intercepts OSC 52 and writes to your clipboard. The same works for `grok wrap docker exec ...` and `grok wrap kubectl exec ...`. `grok wrap` is experimental.

## Keyboard chords do not work

* WezTerm: add `config.enable_kitty_keyboard = true` to `wezterm.lua`, then restart — this fixes `Ctrl+Enter` (interject) and `Shift+Enter` (newline).
* VS Code, Cursor, Windsurf, and Zed terminals cannot distinguish `Shift+Enter` from `Enter`; use `Alt+Enter` for newlines. The same applies to VS Code over SSH.
* Zellij intercepts many Ctrl chords. On Zellij 0.41+, switch to the "Unlock-First (non-colliding)" preset (`Ctrl+O` → `c` → Change Mode Behavior), then `Ctrl+G` temporarily unlocks Zellij when you need it.
* Apple Terminal: `Ctrl+O` interjects (it lacks the Kitty keyboard protocol for `Ctrl+Enter`).

## No fullscreen, or mouse scrolling stops

Grok intentionally runs inline under Zellij and tmux control mode (`tmux -CC`); force fullscreen with `alt_screen = "always"` under `[terminal]` in `~/.grok/pager.toml`, or disable it anywhere with `--no-alt-screen`.

If your terminal's native scrollbar takes over, mouse reporting is off: Apple Terminal re-enables it under View → Allow Mouse Reporting (`Cmd+R`); iTerm2 under Settings → Profiles → Terminal → "Enable mouse reporting".

Still stuck? Run `/feedback`.

#### Features

# Permissions

Permissions decide which tool calls may run. The [sandbox](/build/features/sandbox) is separate: it limits what an approved call can do on the filesystem and network.

## Modes

| Mode | Behavior | Enter via |
| ---- | -------- | --------- |
| Ask (default) | Prompt for anything not already allowed | — |
| Auto | Classifier auto-approves safe tools; dangerous ones may still prompt (`deny` rules and hooks still apply) | `/auto`, `Shift+Tab` when the feature is on |
| Always-approve | Auto-approve tool calls (`deny` rules and PreToolUse hooks still apply) | `/always-approve`, `Ctrl+O`, `Shift+Tab`, `grok --always-approve` |

`Shift+Tab` cycles Normal → Plan → Auto (when available) → Always-approve. `/auto` only appears when the auto permission-mode feature is enabled. Running `/auto` while always-approve is on (or the reverse) switches modes rather than stacking them. Status shows `auto` when auto is active and plan mode is not.

Default in user config only (`~/.grok/config.toml` or managed/requirements — not project `.grok/config.toml`):

```text
[ui]
permission_mode = "auto" # or "ask" | "always-approve"
```

Legacy keys `approval_mode` and `yolo = true` still work; `permission_mode` wins when more than one is set.

[Plan mode](/build/features/plan-mode) is independent: edit tools stay limited while planning, and the plan review UI is not skipped under auto or always-approve.

Headless modes such as `dontAsk` and locking always-approve off: [Enterprise Deployments](/build/enterprise#permissions).

## Allow and deny rules

```text
[permission]
rules = [
  { action = "allow", tool = "bash", pattern = "git *" },
  { action = "allow", tool = "read" },
  { action = "deny",  tool = "bash", pattern = "rm -rf *" },
]
```

`--allow` / `--deny` take the same patterns per invocation. Supported filters include `Bash`, `Edit`, `Read`, `Grep`, `MCPTool`, `WebFetch`, and `WebSearch`. `deny` always wins over `allow`.

A remembered “always allow” grant still prompts for dangerous patterns such as `rm` and `git push`. An explicit config or CLI allow rule auto-approves them. Under always-approve they run unless you add a deny.

#### Features

# Sandbox

The sandbox limits what the agent process and its children can read, write, and reach on the network (Landlock on Linux, Seatbelt on macOS). Off by default. Permissions gate whether a tool call runs; the sandbox limits what an approved call can do — see [Permissions](/build/features/permissions).

## Profiles

| Profile | Filesystem read | Filesystem write | Child network | Use case |
| --- | --- | --- | --- | --- |
| `off` | Unrestricted | Unrestricted | Allowed | No sandbox (default) |
| `workspace` | Everywhere | CWD, `~/.grok/`, temp | Allowed | Normal development |
| `devbox` | Everywhere | Top-level dirs except `/data` | Allowed | Cloud devbox environments |
| `read-only` | Everywhere | `~/.grok/` and temp only | Blocked | Code review, auditing |
| `strict` | CWD and system paths | CWD, `~/.grok/`, temp | Blocked | Untrusted repositories |

| Limitation | Detail |
| --- | --- |
| Child network | Enforced on Linux only; no-op on macOS for `read-only` / `strict` |
| Credentials | Built-ins do not permanently protect paths such as `~/.ssh`; use a custom `deny` list |
| `~/.grok/` | Stays writable under sandboxed profiles so sessions can persist |
| In-process network | Model API and web tools are not blocked by child-network settings |

## Enable a profile

| Mechanism | Example |
| --- | --- |
| CLI | `grok --sandbox workspace` |
| Config | `[sandbox] profile = "workspace"` in `~/.grok/config.toml` |
| Env | `GROK_SANDBOX=workspace` |
| Managed pin | `requirements.toml` (can override CLI) — [Enterprise](/build/enterprise#sandbox) |

## Custom profiles

Define named profiles in `~/.grok/sandbox.toml` or project `.grok/sandbox.toml`:

```text
[profiles.my-profile]
extends = "workspace"
restrict_network = true
deny = ["/secrets", "**/.env", "**/*.pem"]
```

Select with `--sandbox my-profile` or `[sandbox] profile`. Built-in names cannot be redefined for selection. Field details: [Settings Reference](/build/settings/reference).

For untrusted trees, pair a strict profile with narrow [permission](/build/features/permissions) allows (or headless `dontAsk`).
