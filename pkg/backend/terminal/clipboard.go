// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package terminal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ClipboardProvider defines the interface for clipboard operations
type ClipboardProvider interface {
	Get() (string, error)
	Set(content string) error
}

// SystemClipboard implements platform-specific clipboard operations
type SystemClipboard struct{}

func (c *SystemClipboard) Get() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return c.runCommand("pbpaste")
	case "linux":
		// Try xclip first, then xsel, then wayland
		if content, err := c.runCommand("xclip", "-selection", "clipboard", "-o"); err == nil {
			return content, nil
		}
		if content, err := c.runCommand("xsel", "--clipboard", "--output"); err == nil {
			return content, nil
		}
		if content, err := c.runCommand("wl-paste"); err == nil {
			return content, nil
		}
		return "", fmt.Errorf("no clipboard utility found")
	case "windows":
		return c.runCommand("powershell.exe", "-command", "Get-Clipboard")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (c *SystemClipboard) Set(content string) error {
	switch runtime.GOOS {
	case "darwin":
		return c.writeCommand(content, "pbcopy")
	case "linux":
		// Try xclip first, then xsel, then wayland
		if err := c.writeCommand(content, "xclip", "-selection", "clipboard"); err == nil {
			return nil
		}
		if err := c.writeCommand(content, "xsel", "--clipboard", "--input"); err == nil {
			return nil
		}
		if err := c.writeCommand(content, "wl-copy"); err == nil {
			return nil
		}
		return fmt.Errorf("no clipboard utility found")
	case "windows":
		return c.writeCommand(content, "powershell.exe", "-command", "Set-Clipboard")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (c *SystemClipboard) runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (c *SystemClipboard) writeCommand(content string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}

// FallbackClipboard provides in-memory clipboard when system clipboard is unavailable
type FallbackClipboard struct {
	content string
}

func (c *FallbackClipboard) Get() (string, error) {
	return c.content, nil
}

func (c *FallbackClipboard) Set(content string) error {
	c.content = content
	return nil
}
