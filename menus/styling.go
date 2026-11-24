package menus

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
)

var (
	SuccessStyle = color.New(color.FgGreen, color.Bold)
	ErrorStyle   = color.New(color.FgRed, color.Bold)
	InfoStyle    = color.New(color.FgCyan)
	WarningStyle = color.New(color.FgYellow)
	HeaderStyle  = color.New(color.FgMagenta, color.Bold)
	PromptStyle  = color.New(color.FgBlue, color.Bold)
)

func PrintHeader(text string) {
	pterm.DefaultBox.WithTitle(text).WithTitleTopCenter().WithBoxStyle(pterm.NewStyle(pterm.FgCyan)).Println()
}

func PrintAppHeader(name, version string) {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("GO", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("CLI", pterm.NewStyle(pterm.FgLightMagenta)),
	).Render()

	info := pterm.DefaultBox.WithTitle("Application Info").
		WithTitleTopLeft().
		WithBoxStyle(pterm.NewStyle(pterm.FgCyan))
	
	info.Println(pterm.Sprintf("%s: %s\n%s: %s", 
		pterm.FgYellow.Sprint("Name"), 
		pterm.FgWhite.Sprint(name),
		pterm.FgYellow.Sprint("Version"), 
		pterm.FgWhite.Sprint(version)))
}

func PrintSuccess(message string) {
	pterm.Success.Println(message)
}

func PrintError(message string) {
	pterm.Error.Println(message)
}

func PrintInfo(message string) {
	pterm.Info.Println(message)
}

func PrintWarning(message string) {
	pterm.Warning.Println(message)
}

func StyledPrompt(label string) string {
	return PromptStyle.Sprint("→ ") + InfoStyle.Sprint(label)
}

func PrintSection(title string) {
	pterm.DefaultSection.WithStyle(pterm.NewStyle(pterm.FgLightCyan)).Println(title)
}

type ProgressReader struct {
	io.Reader
	bar *progressbar.ProgressBar
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.bar.Add(n)
	return n, err
}

func FetchWithProgress(url string, description string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	
	spinner, _ := pterm.DefaultSpinner.Start("Connecting to ", url)
	
	resp, err := client.Get(url)
	if err != nil {
		spinner.Fail("Failed to connect")
		return nil, err
	}
	
	spinner.Success("Connected successfully")
	
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		contentLength = 4096
	}
	
	bar := progressbar.NewOptions64(
		contentLength,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(pterm.DefaultProgressbar.Writer),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	
	resp.Body = io.NopCloser(&ProgressReader{
		Reader: resp.Body,
		bar:    bar,
	})
	
	return resp, nil
}

type ProgressWriter struct {
	bar *progressbar.ProgressBar
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.bar.Add(n)
	return n, nil
}

func PostWithProgress(url string, body io.Reader, contentLength int64, contentType string, description string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	
	spinner, _ := pterm.DefaultSpinner.Start("Connecting to ", url)
	
	if contentLength <= 0 {
		contentLength = 4096
	}
	
	bar := progressbar.NewOptions64(
		contentLength,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(pterm.DefaultProgressbar.Writer),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	
	progressReader := io.TeeReader(body, &ProgressWriter{bar: bar})
	
	req, err := http.NewRequest(http.MethodPost, url, progressReader)
	if err != nil {
		spinner.Fail("Failed to create request")
		return nil, err
	}
	
	req.Header.Set("Content-Type", contentType)
	req.ContentLength = contentLength
	
	spinner.Success("Connected successfully")
	
	resp, err := client.Do(req)
	if err != nil {
		PrintError("Failed to send data")
		return nil, err
	}
	
	return resp, nil
}

func PrintMenuHeader(title string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		WithMargin(2).
		Println(title)
}

func GetIcon(name string) string {
	icons := map[string]string{
		"pipeline":  "🔄",
		"source":    "📥",
		"destination": "📤",
		"transform": "⚙️",
		"sequence":  "📋",
		"state":     "💾",
		"success":   "✅",
		"error":     "❌",
		"info":      "ℹ️",
		"warning":   "⚠️",
		"menu":      "📂",
		"action":    "▶️",
		"back":      "◀️",
		"quit":      "🚪",
		"add":       "➕",
		"delete":    "🗑️",
		"view":      "👁️",
		"run":       "🚀",
		"save":      "💾",
		"load":      "📂",
	}
	
	if icon, ok := icons[name]; ok {
		return icon + " "
	}
	return ""
}
