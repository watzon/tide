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
	"sync"
	"testing"
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

func TestSystemClipboard(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system clipboard test in short mode")
	}

	t.Run("system clipboard operations", func(t *testing.T) {
		clipboard := &SystemClipboard{}
		testContent := "tide clipboard test content"

		// Skip test on unsupported platforms
		if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
			t.Skipf("system clipboard not supported on %s", runtime.GOOS)
		}

		err := clipboard.Set(testContent)
		if err != nil {
			t.Logf("system clipboard set failed (might be normal if no clipboard utility available): %v", err)
			return
		}

		content, err := clipboard.Get()
		if err != nil {
			t.Errorf("unexpected error getting clipboard: %v", err)
			return
		}

		if content != testContent {
			t.Errorf("expected content %q, got %q", testContent, content)
		}
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
