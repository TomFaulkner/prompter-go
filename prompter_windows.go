// +build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BurntSushi/toml"
)

func main() {
	prompts, err := loadPrompts()
	if err != nil {
		fmt.Println("Error loading prompts:", err)
		os.Exit(1)
	}

	a := app.New()
	w := a.NewWindow("Prompt Manager")
	w.Resize(fyne.NewSize(600, 400))

	var items []string
	for name, content := range prompts {
		preview := content
		if len(content) > 50 {
			preview = content[:50] + "..."
		}
		items = append(items, fmt.Sprintf("%s: %s", name, preview))
	}

	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		name := strings.SplitN(items[id], ":", 2)[0]
		if fullPrompt, ok := prompts[strings.TrimSpace(name)]; ok {
			err := copyToClipboard(fullPrompt)
			if err != nil {
				fmt.Println("Clipboard error:", err)
			}
			w.Close()
		}
	}

	w.SetContent(container.NewVScroll(list))
	w.ShowAndRun()
}

func loadPrompts() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	filePath := homeDir + "\\prompts.toml" // Windows uses backslashes

	var prompts map[string]string
	_, err = toml.DecodeFile(filePath, &prompts)
	if err != nil {
		return nil, err
	}
	if len(prompts) == 0 {
		return nil, fmt.Errorf("no prompts found in %s", filePath)
	}
	return prompts, nil
}

func copyToClipboard(text string) error {
	cmd := exec.Command("clip.exe")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err == nil {
		return nil
	}
	return fmt.Errorf("failed to copy to clipboard with clip.exe")
}
