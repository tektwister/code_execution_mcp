package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/aravi/code_execution_mcp/internal/core/domain"
	"github.com/aravi/code_execution_mcp/internal/core/ports"
)

// ShellExecutor implements CodeExecutor for Bash/Zsh scripts
type ShellExecutor struct{}

// NewShellExecutor creates a new shell executor
func NewShellExecutor() ports.CodeExecutor {
	return &ShellExecutor{}
}

// Supports checks if this executor supports the given language
func (e *ShellExecutor) Supports(language string) bool {
	return language == "bash" || language == "zsh" || language == "shell"
}

// Execute runs a bash/zsh script
func (e *ShellExecutor) Execute(ctx context.Context, req domain.ExecutionRequest) (*domain.ExecutionResult, error) {
	if strings.TrimSpace(req.Script) == "" {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.ValidationError,
			Stderr:    "Script cannot be empty",
		}, nil
	}

	timeout := getTimeout(req.Timeout, 30)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Determine shell based on OS
	var shell, flag string
	var fullScript string
	
	if runtime.GOOS == "windows" {
		// On Windows, try PowerShell first (most reliable), then Git Bash, then WSL
		shell, flag = detectWindowsShell()
		
		if shell == "powershell.exe" || shell == "pwsh.exe" {
			// For PowerShell, we need to use different syntax
			fullScript = req.Script
			if len(req.Args) > 0 {
				// PowerShell uses $args array
				argsStr := strings.Join(req.Args, " ")
				fullScript = fmt.Sprintf("$args = @('%s'); %s", 
					strings.ReplaceAll(argsStr, "'", "''"), req.Script)
			}
		} else {
			// For bash/zsh, use standard shell quoting
			fullScript = req.Script
			if len(req.Args) > 0 {
				fullScript = fmt.Sprintf("set -- %s; %s", shellQuoteArgs(req.Args), req.Script)
			}
		}
	} else {
		// Unix-like systems
		shell = "bash"
		flag = "-c"
		fullScript = req.Script
		if len(req.Args) > 0 {
			fullScript = fmt.Sprintf("set -- %s; %s", shellQuoteArgs(req.Args), req.Script)
		}
	}

	cmd := exec.CommandContext(ctx, shell, flag, fullScript)
	if req.WorkingDir != "" {
		cmd.Dir = req.WorkingDir
	}

	return executeCommand(cmd)
}

// detectWindowsShell finds the best available shell on Windows
func detectWindowsShell() (shell, flag string) {
	// Try PowerShell first (most universally available on Windows)
	if _, err := exec.LookPath("pwsh.exe"); err == nil {
		return "pwsh.exe", "-Command"
	}
	if _, err := exec.LookPath("powershell.exe"); err == nil {
		return "powershell.exe", "-Command"
	}
	
	// Try Git Bash (common for developers)
	gitBashPaths := []string{
		"C:\\Program Files\\Git\\bin\\bash.exe",
		"C:\\Program Files (x86)\\Git\\bin\\bash.exe",
	}
	for _, path := range gitBashPaths {
		if _, err := exec.LookPath(path); err == nil {
			return path, "-c"
		}
	}
	
	// Try bash in PATH (might be Git Bash added to PATH)
	if _, err := exec.LookPath("bash.exe"); err == nil {
		return "bash.exe", "-c"
	}
	
	// Fallback to bash (will use WSL if configured)
	return "bash", "-c"
}

// shellQuoteArgs properly quotes arguments for shell execution
func shellQuoteArgs(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		escaped := strings.ReplaceAll(arg, "'", "'\"'\"'")
		quoted[i] = fmt.Sprintf("'%s'", escaped)
	}
	return strings.Join(quoted, " ")
}

// executeCommand runs a command and returns the result
func executeCommand(cmd *exec.Cmd) (*domain.ExecutionResult, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	exitCode := 0
	errorType := domain.NoError

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			errorType = domain.RuntimeError
		} else {
			exitCode = -1
			errorType = domain.SystemError
		}
	}

	return &domain.ExecutionResult{
		ExitCode:  exitCode,
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		Duration:  duration,
		IsError:   exitCode != 0,
		ErrorType: errorType,
	}, nil
}

// getTimeout returns a valid timeout duration
func getTimeout(requested, defaultSecs int) time.Duration {
	if requested <= 0 {
		return time.Duration(defaultSecs) * time.Second
	}
	if requested > 300 {
		return 300 * time.Second
	}
	return time.Duration(requested) * time.Second
}
