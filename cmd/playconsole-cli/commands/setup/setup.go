package setup

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for Google Play Console CLI",
	Long: `Guide you through setting up the CLI with your Google Cloud project.

This wizard will help you:
  1. Create a Google Cloud Service Account
  2. Enable the Google Play Android Developer API
  3. Download and configure your credentials
  4. Link the service account in Play Console`,
	Run: runSetup,
}

type step struct {
	title       string
	description string
	url         string
	action      string
}

var setupSteps = []step{
	{
		title:       "Create Service Account",
		description: "Create a new service account in Google Cloud Console",
		url:         "https://console.cloud.google.com/iam-admin/serviceaccounts/create",
		action:      "Name it 'play-console-cli' and click CREATE AND CONTINUE",
	},
	{
		title:       "Enable Google Play Developer API",
		description: "Enable the API that allows managing your Play Console",
		url:         "https://console.cloud.google.com/apis/library/androidpublisher.googleapis.com",
		action:      "Click the blue ENABLE button",
	},
	{
		title:       "Create Service Account Key",
		description: "Download JSON credentials for your service account",
		url:         "https://console.cloud.google.com/iam-admin/serviceaccounts",
		action:      "Click your service account → Keys tab → Add Key → JSON → Download",
	},
	{
		title:       "Link in Play Console",
		description: "Grant API access to your service account",
		url:         "https://play.google.com/console/developers/api-access",
		action:      "Link your Cloud project, find the service account, click 'Grant access'",
	},
}

func runSetup(cmd *cobra.Command, args []string) {
	clearScreen()
	printHeader()

	reader := bufio.NewReader(os.Stdin)

	for i, s := range setupSteps {
		printStep(i+1, len(setupSteps), s)

		fmt.Println("\nOptions:")
		fmt.Println("  [o] Open this page in browser")
		fmt.Println("  [c] Copy URL to clipboard")
		fmt.Println("  [s] Skip (already done)")
		fmt.Println("  [q] Quit setup")
		fmt.Print("\nChoice [o/c/s/q]: ")

		input := readInput(reader)

		switch input {
		case "o", "open", "":
			openBrowser(s.url)
			waitForUser(reader, "Press ENTER when you've completed this step...")
		case "c", "copy":
			copyToClipboard(s.url)
			fmt.Println("✓ URL copied to clipboard")
			waitForUser(reader, "Press ENTER when you've completed this step...")
		case "s", "skip":
			fmt.Println("→ Skipping...")
		case "q", "quit":
			fmt.Println("\nSetup cancelled. Run 'gpc setup' to continue later.")
			return
		}

		clearScreen()
	}

	// Final step: configure credentials
	printFinalStep(reader)
}

func printHeader() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           GOOGLE PLAY CONSOLE CLI - SETUP WIZARD                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("This wizard will guide you through the one-time setup process.")
	fmt.Println("Each step will open ONE page - no loops, no spam.")
	fmt.Println()
	waitForUserSimple()
}

func printStep(current, total int, s step) {
	fmt.Printf("\n┌─ STEP %d of %d ─────────────────────────────────────────────────────┐\n", current, total)
	fmt.Printf("│  %s\n", s.title)
	fmt.Printf("├────────────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  %s\n", s.description)
	fmt.Printf("│\n")
	fmt.Printf("│  ACTION: %s\n", s.action)
	fmt.Printf("│\n")
	fmt.Printf("│  URL: %s\n", s.url)
	fmt.Printf("└────────────────────────────────────────────────────────────────────┘\n")
}

func printFinalStep(reader *bufio.Reader) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    CONFIGURE CREDENTIALS                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Now let's configure your downloaded credentials.")
	fmt.Println()

	// Check for existing credentials in Downloads
	homeDir, _ := os.UserHomeDir()
	downloadsDir := filepath.Join(homeDir, "Downloads")

	jsonFiles := findJSONFiles(downloadsDir)

	var credPath string

	if len(jsonFiles) > 0 {
		fmt.Println("Found potential credential files in Downloads:")
		for i, f := range jsonFiles {
			fmt.Printf("  [%d] %s\n", i+1, filepath.Base(f))
		}
		fmt.Println("  [m] Enter path manually")
		fmt.Print("\nSelect file number or 'm': ")

		input := readInput(reader)
		if input == "m" {
			credPath = promptForPath(reader)
		} else {
			var idx int
			fmt.Sscanf(input, "%d", &idx)
			if idx >= 1 && idx <= len(jsonFiles) {
				credPath = jsonFiles[idx-1]
			} else {
				credPath = promptForPath(reader)
			}
		}
	} else {
		credPath = promptForPath(reader)
	}

	if credPath == "" {
		fmt.Println("\nNo credentials configured. Run 'gpc auth login --credentials <path>' later.")
		return
	}

	// Create config directory
	configDir := filepath.Join(homeDir, ".config", "gpc")
	os.MkdirAll(configDir, 0700)

	// Copy credentials
	destPath := filepath.Join(configDir, "service-account.json")

	if err := copyFile(credPath, destPath); err != nil {
		fmt.Printf("Error copying credentials: %v\n", err)
		fmt.Println("You can manually copy your credentials to:", destPath)
		return
	}

	// Set permissions
	os.Chmod(destPath, 0600)

	fmt.Println()
	fmt.Println("✓ Credentials saved to:", destPath)
	fmt.Println()

	// Ask for default package
	fmt.Print("Enter your app's package name (or press ENTER to skip): ")
	packageName := readInput(reader)

	// Run auth login
	fmt.Println()
	fmt.Println("Configuring CLI...")

	authArgs := []string{"auth", "login", "--name", "default", "--credentials", destPath}
	if packageName != "" {
		authArgs = append(authArgs, "--default-package", packageName)
	}

	// Get executable path
	execPath, _ := os.Executable()
	authCmd := exec.Command(execPath, authArgs...)
	authCmd.Stdout = os.Stdout
	authCmd.Stderr = os.Stderr

	if err := authCmd.Run(); err != nil {
		fmt.Printf("\nWarning: Could not auto-configure. Run manually:\n")
		fmt.Printf("  gpc auth login --name default --credentials %s\n", destPath)
	} else {
		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                      SETUP COMPLETE! 🎉                          ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("You can now use gpc commands:")
		fmt.Println("  gpc tracks list")
		fmt.Println("  gpc bundles upload --file app.aab --track internal")
		fmt.Println("  gpc reviews list")
		fmt.Println()
		if packageName != "" {
			fmt.Printf("Default package: %s\n", packageName)
		} else {
			fmt.Println("Tip: Use --package <name> or set GPC_PACKAGE env var")
		}
	}
}

func promptForPath(reader *bufio.Reader) string {
	fmt.Print("Enter path to your service account JSON file: ")
	path := readInput(reader)

	// Expand ~ if present
	if strings.HasPrefix(path, "~") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[1:])
	}

	// Clean path to prevent traversal attacks
	path = filepath.Clean(path)

	// Ensure it's a JSON file
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		fmt.Println("Error: File must be a .json file")
		return ""
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", path)
		return ""
	}

	return path
}

func findJSONFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Look for service account files
		if strings.HasSuffix(name, ".json") &&
			(strings.Contains(strings.ToLower(name), "service") ||
			 strings.Contains(strings.ToLower(name), "key") ||
			 strings.Contains(strings.ToLower(name), "credential")) {
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}

func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func waitForUser(reader *bufio.Reader, msg string) {
	fmt.Print("\n" + msg)
	reader.ReadString('\n')
}

func waitForUserSimple() {
	fmt.Print("Press ENTER to begin...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func clearScreen() {
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
