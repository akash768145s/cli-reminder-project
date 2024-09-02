package main

import(
    "fmt"            // For formatted I/O operations
    "os"             // For interacting with the operating system (e.g., command-line arguments, environment variables)
    "os/exec"        // For running external commands
    "strings"        // For string manipulation
    "time"           // For time-related operations (e.g., getting the current time, sleeping)
    "github.com/gen2brain/beeep" // For desktop notifications
    "github.com/olebedev/when"   // For natural language date and time parsing
    "github.com/olebedev/when/rules/common" // Common date/time parsing rules
    "github.com/olebedev/when/rules/en"     // English-specific date/time parsing rules
)

const(
    markName = "GOLAND_CLI_REMINDER" // Constant to store the name of the environment variable used to mark the reminder process
    markValue = "1"                  // Constant to store the value of the environment variable
)

func main() {
    // Check if there are fewer than 3 command-line arguments (program name, time, and message)
    if len(os.Args) < 3 {
        fmt.Printf("Usage: %s <hh:mm> <text message>\n", os.Args[0]) // Print usage instructions
        os.Exit(1) // Exit with code 1 (indicating an error)
    }
    
    now := time.Now() // Get the current time

    w := when.New(nil) // Create a new `when` parser with no custom configuration
    w.Add(en.All...)   // Add all English language rules for date/time parsing
    w.Add(common.All...) // Add common date/time parsing rules

    // Parse the time string provided as the first command-line argument
    t, err := w.Parse(os.Args[1], now)
    if err != nil { // If there's an error parsing the time
        fmt.Printf("Error: %s\n", err) // Print the error
        os.Exit(1) // Exit with code 1
    }

    // If the time couldn't be parsed
    if t == nil {
        fmt.Println("Unable to parse time!") // Print an error message
        os.Exit(2) // Exit with code 2
    }

    // If the parsed time is in the past
    if now.After(t.Time) {
        fmt.Println("Set a future time!") // Print an error message
        os.Exit(3) // Exit with code 3
    }

    // Calculate the difference between the parsed time and the current time
    diff := t.Time.Sub(now)

    // Check if the environment variable `GOLAND_CLI_REMINDER` is set to "1"
    if os.Getenv("markName") == markValue {
        time.Sleep(diff) // Wait until the specified time
        // Show a desktop notification with the reminder message
        beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
        if err != nil { // If there's an error showing the notification
            fmt.Printf("Error: %s\n", err) // Print the error
            os.Exit(4) // Exit with code 4
        }
    } else {
        // If the environment variable is not set, re-run the program with the environment variable set
        cmd := exec.Command(os.Args[0], os.Args[1:]...) // Create a command to re-run the program with the same arguments
        cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue)) // Set the environment variable in the new process
        if err := cmd.Start(); err != nil { // If there's an error starting the new process
            fmt.Println(err) // Print the error
            os.Exit(5) // Exit with code 5
        }
        // Print a message indicating the reminder has been set and the time remaining until it triggers
        fmt.Println("Reminder set!", diff.Round(time.Second))
        os.Exit(0) // Exit with code 0 (indicating success)
    }
}
