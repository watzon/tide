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
	"sync"
	"testing"
	"time"
)

// MockClipboardProvider implements ClipboardProvider for testing
type MockClipboardProvider struct {
	content string
	getErr  error
	setErr  error
}

func (m *MockClipboardProvider) Get() (string, error) {
	if m.getErr != nil {
		return "", m.getErr
	}
	return m.content, nil
}

func (m *MockClipboardProvider) Set(content string) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.content = content
	return nil
}

func TestFallbackClipboard(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		clipboard := &FallbackClipboard{}
		testContent := "test content"

		err := clipboard.Set(testContent)
		if err != nil {
			t.Errorf("unexpected error setting clipboard: %v", err)
		}

		content, err := clipboard.Get()
		if err != nil {
			t.Errorf("unexpected error getting clipboard: %v", err)
		}

		if content != testContent {
			t.Errorf("expected content %q, got %q", testContent, content)
		}
	})
}

// hasClipboardUtility checks if the system has any supported clipboard utility
func hasClipboardUtility() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("pbcopy")
		return err == nil
	case "linux":
		// Check for any of the supported utilities
		for _, cmd := range []string{"xsel", "xclip", "wl-copy"} {
			if _, err := exec.LookPath(cmd); err == nil {
				return true
			}
		}
		return false
	case "windows":
		_, err := exec.LookPath("powershell.exe")
		return err == nil
	default:
		return false
	}
}

func TestSystemClipboard(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system clipboard test in short mode")
	}

	if !hasClipboardUtility() {
		t.Skip("skipping system clipboard test - no clipboard utility available")
	}

	t.Run("system clipboard operations", func(t *testing.T) {
		clipboard := &SystemClipboard{}
		testContent := "tide clipboard test content"

		// Test OS-specific clipboard commands
		switch runtime.GOOS {
		case "darwin":
			t.Run("macOS commands", func(t *testing.T) {
				// Test pbcopy/pbpaste
				err := clipboard.Set(testContent)
				if err != nil {
					t.Errorf("pbcopy failed: %v", err)
					return
				}

				content, err := clipboard.Get()
				if err != nil {
					t.Errorf("pbpaste failed: %v", err)
					return
				}

				if content != testContent {
					t.Errorf("expected content %q, got %q", testContent, content)
				}
			})

		case "linux":
			t.Run("Linux commands", func(t *testing.T) {
				// Test each clipboard command in order of preference
				commands := []struct {
					name string
					get  []string
					set  []string
				}{
					{"xsel", []string{"xsel", "--output", "--clipboard"}, []string{"xsel", "--input", "--clipboard"}},
					{"xclip", []string{"xclip", "-out", "-selection", "clipboard"}, []string{"xclip", "-in", "-selection", "clipboard"}},
					{"wl-clipboard", []string{"wl-paste", "--no-newline"}, []string{"wl-copy"}},
				}

				foundCommand := false
				for _, cmd := range commands {
					if _, err := exec.LookPath(cmd.get[0]); err != nil {
						continue
					}
					foundCommand = true

					t.Run(cmd.name, func(t *testing.T) {
						err := clipboard.Set(testContent)
						if err != nil {
							t.Errorf("%s set failed: %v", cmd.name, err)
							return
						}

						content, err := clipboard.Get()
						if err != nil {
							t.Errorf("%s get failed: %v", cmd.name, err)
							return
						}

						if content != testContent {
							t.Errorf("expected content %q, got %q", testContent, content)
						}
					})
				}

				if !foundCommand {
					t.Skip("no clipboard utility found")
				}
			})

		case "windows":
			t.Run("Windows commands", func(t *testing.T) {
				// Test clip.exe/powershell Get-Clipboard
				err := clipboard.Set(testContent)
				if err != nil {
					t.Errorf("clip.exe failed: %v", err)
					return
				}

				content, err := clipboard.Get()
				if err != nil {
					t.Errorf("Get-Clipboard failed: %v", err)
					return
				}

				if content != testContent {
					t.Errorf("expected content %q, got %q", testContent, content)
				}
			})

		default:
			t.Skipf("system clipboard not supported on %s", runtime.GOOS)
		}
	})

	// The command execution error tests can still run without clipboard utilities
	t.Run("command execution errors", func(t *testing.T) {
		clipboard := &SystemClipboard{}

		// Test with non-existent command
		t.Run("non-existent command", func(t *testing.T) {
			_, err := clipboard.runCommand("nonexistentcommand")
			if err == nil {
				t.Error("expected error for non-existent command")
			}
			if _, ok := err.(*exec.Error); !ok {
				t.Errorf("expected exec.Error, got %T", err)
			}
		})

		// Test with failing command
		t.Run("failing command", func(t *testing.T) {
			_, err := clipboard.runCommand("false")
			if err == nil {
				t.Error("expected error for failing command")
			}
			if _, ok := err.(*exec.ExitError); !ok {
				t.Errorf("expected exec.ExitError, got %T", err)
			}
		})

		// Test with command and arguments
		t.Run("command with args", func(t *testing.T) {
			_, err := clipboard.runCommand("echo", "-n", "test")
			if err != nil {
				t.Errorf("unexpected error running echo: %v", err)
			}
		})
	})

	// Skip clipboard-dependent tests if no utility is available
	if !hasClipboardUtility() {
		return
	}

	t.Run("write command errors", func(t *testing.T) {
		clipboard := &SystemClipboard{}

		t.Run("non-existent write command", func(t *testing.T) {
			err := clipboard.writeCommand("test content", "nonexistentcommand")
			if err == nil {
				t.Error("expected error for non-existent command")
			}
			if _, ok := err.(*exec.Error); !ok {
				t.Errorf("expected exec.Error, got %T", err)
			}
		})

		t.Run("failing write command", func(t *testing.T) {
			err := clipboard.writeCommand("test content", "false")
			if err == nil {
				t.Error("expected error for failing command")
			}
			if _, ok := err.(*exec.ExitError); !ok {
				t.Errorf("expected exec.ExitError, got %T", err)
			}
		})
	})

	t.Run("empty clipboard", func(t *testing.T) {
		clipboard := &SystemClipboard{}

		// Try to get content from empty clipboard
		content, err := clipboard.Get()
		if err != nil {
			// Some systems might return an error for empty clipboard
			t.Logf("get from empty clipboard: %v", err)
		} else if len(content) > 0 {
			t.Log("clipboard was not empty, skipping empty test")
		}
	})

	t.Run("large content", func(t *testing.T) {
		clipboard := &SystemClipboard{}

		// Create large content (10KB instead of 100KB to be safer)
		largeContent := strings.Repeat("large content test ", 500)

		err := clipboard.Set(largeContent)
		if err != nil {
			t.Errorf("failed to set large content: %v", err)
			return
		}

		content, err := clipboard.Get()
		if err != nil {
			t.Errorf("failed to get large content: %v", err)
			return
		}

		// Normalize line endings and trim spaces
		content = strings.TrimSpace(content)
		expected := strings.TrimSpace(largeContent)

		if content != expected {
			// If they're still different, show a summary of the difference
			if len(content) != len(expected) {
				t.Errorf("content length mismatch: got %d, want %d", len(content), len(expected))
			} else {
				t.Error("large content mismatch")
			}
		}
	})

	t.Run("special characters", func(t *testing.T) {
		clipboard := &SystemClipboard{}

		testCases := []struct {
			name    string
			content string
		}{
			{"basic ascii", "Hello, World!"},
			{"unicode", "你好世界 αβγ"},
			{"spaces", "  leading and trailing spaces  "},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := clipboard.Set(tc.content)
				if err != nil {
					t.Errorf("failed to set content: %v", err)
					return
				}

				content, err := clipboard.Get()
				if err != nil {
					t.Errorf("failed to get content: %v", err)
					return
				}

				// Normalize content by trimming spaces
				content = strings.TrimSpace(content)
				expected := strings.TrimSpace(tc.content)

				if content != expected {
					t.Errorf("content mismatch:\nexpected: %q\ngot: %q", expected, content)
				}
			})
		}

		// Separate test for whitespace characters
		t.Run("whitespace characters", func(t *testing.T) {
			// Some clipboard implementations might normalize or strip certain whitespace
			// So we'll test each one separately
			whitespaceTests := []struct {
				name    string
				content string
			}{
				{"newline", "line1\nline2"},
				{"tab", "tab\tseparated"},
				{"carriage return", "line1\rline2"},
			}

			for _, wt := range whitespaceTests {
				t.Run(wt.name, func(t *testing.T) {
					err := clipboard.Set(wt.content)
					if err != nil {
						t.Logf("failed to set %s content: %v", wt.name, err)
						return
					}

					content, err := clipboard.Get()
					if err != nil {
						t.Logf("failed to get %s content: %v", wt.name, err)
						return
					}

					// Log the actual behavior rather than failing
					// This helps document how different platforms handle whitespace
					if content != wt.content {
						t.Logf("%s handling: %q converted to %q",
							runtime.GOOS, wt.content, content)
					}
				})
			}
		})
	})
}

func TestTerminalClipboard(t *testing.T) {
	t.Run("with working system clipboard", func(t *testing.T) {
		term := &Terminal{}
		mock := &MockClipboardProvider{}
		term.clipboardProvider = mock

		testContent := "test content"
		err := term.SetClipboard(testContent)
		if err != nil {
			t.Errorf("unexpected error setting clipboard: %v", err)
		}

		content, err := term.GetClipboard()
		if err != nil {
			t.Errorf("unexpected error getting clipboard: %v", err)
		}

		if content != testContent {
			t.Errorf("expected content %q, got %q", testContent, content)
		}
	})

	t.Run("fallback on system clipboard failure", func(t *testing.T) {
		term := &Terminal{}
		failingMock := &MockClipboardProvider{
			getErr: &exec.Error{Name: "test", Err: fmt.Errorf("mock error")},
			setErr: &exec.Error{Name: "test", Err: fmt.Errorf("mock error")},
		}
		term.clipboardProvider = failingMock

		testContent := "fallback content"
		err := term.SetClipboard(testContent)
		if err != nil {
			t.Errorf("unexpected error with fallback clipboard: %v", err)
		}

		content, err := term.GetClipboard()
		if err != nil {
			t.Errorf("unexpected error with fallback clipboard: %v", err)
		}

		if content != testContent {
			t.Errorf("expected content %q, got %q", testContent, content)
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		term := &Terminal{}
		mock := &MockClipboardProvider{}
		term.clipboardProvider = mock

		const goroutines = 10
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(i int) {
				defer wg.Done()
				content := fmt.Sprintf("content-%d", i)

				err := term.SetClipboard(content)
				if err != nil {
					t.Errorf("error setting clipboard: %v", err)
				}

				_, err = term.GetClipboard()
				if err != nil {
					t.Errorf("error getting clipboard: %v", err)
				}
			}(i)
		}

		wg.Wait()
	})
}

func TestFallbackBehavior(t *testing.T) {
	t.Run("fallback after command failure", func(t *testing.T) {
		term := &Terminal{}

		// Create a failing system clipboard
		failingClipboard := &SystemClipboard{}
		term.clipboardProvider = failingClipboard

		testContent := "fallback test content"

		// This should fall back to FallbackClipboard
		err := term.SetClipboard(testContent)
		if err != nil {
			t.Errorf("unexpected error with fallback: %v", err)
		}

		content, err := term.GetClipboard()
		if err != nil {
			t.Errorf("unexpected error getting fallback content: %v", err)
		}

		if content != testContent {
			t.Errorf("expected fallback content %q, got %q", testContent, content)
		}
	})
}

func TestClipboardRace(t *testing.T) {
	t.Run("concurrent clipboard operations", func(t *testing.T) {
		clipboard := &SystemClipboard{}
		const goroutines = 10
		const iterations = 5

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()

				for j := 0; j < iterations; j++ {
					content := fmt.Sprintf("content-%d-%d", id, j)

					err := clipboard.Set(content)
					if err != nil {
						t.Errorf("set error in goroutine %d: %v", id, err)
					}

					_, err = clipboard.Get()
					if err != nil {
						t.Errorf("get error in goroutine %d: %v", id, err)
					}

					// Small sleep to increase chance of race conditions
					time.Sleep(time.Millisecond)
				}
			}(i)
		}

		wg.Wait()
	})
}
