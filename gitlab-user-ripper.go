package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"flag"
)

func main() {

	var banner = `
 _____ _ _   __        _       _____                 _____ _
|   __|_| |_|  |   ___| |_ ___|  |  |___ ___ ___ ___| __  |_|___ ___ ___ ___
|  |  | |  _|  |__| .'| . |___|  |  |_ -| -_|  _|___|    -| | . | . | -_|  _|
|_____|_|_| |_____|__,|___|   |_____|___|___|_|     |__|__|_|  _|  _|___|_|
                                                            |_| |_|
[*] A GitLab user enumeration tool that rips
[*] Author: RandomChugokujin

`

	fmt.Print(banner)

	// Define flags
	urlFlag := flag.String("u", "", "Base URL to scan (e.g., http://gitlab.local:8081)")
	fileFlag := flag.String("f", "", "Path to username file")
	workersFlag := flag.Int("t", 50, "Number of threads")
	verboseFlag := flag.Bool("v", false, "Verbose Output")

	flag.Parse()

	// Validate required flags
	if *urlFlag == "" || *fileFlag == "" {
		fmt.Println("[-] Usage: gitlab-user-ripper -u <url> -f <username_file> [-t <threads>] [-v]")
		return
	}

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Printf("[-] Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Print options
	fmt.Println("[*] Starting user enumeration with the following options:")
	fmt.Printf("[*] URL: %s\n", *urlFlag)
	fmt.Printf("[*] Username File: %s\n", *fileFlag)
	fmt.Printf("[*] Threads: %d\n", *workersFlag)
	fmt.Printf("[*] Verbose: %t\n", *verboseFlag)

	jobs := make(chan string)
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)

	// Start Workers
	for i := 0; i < *workersFlag; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for username := range jobs {
				checkUser(*urlFlag, username, *verboseFlag)
			}
		}()
	}
	// Feed jobs
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Printf("[-] Error reading file: %v\n", err)
		}
		jobs <- scanner.Text()
	}
	close(jobs)

	wg.Wait()
}

func checkUser(baseURL, username string, verbose bool) {
	url := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(username))
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Printf("[-] Error creating request for %s: %v\n", username, err)
		return
	}

	client := &http.Client{
		// Disable automatic redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[-] Error sending request for %s: %v\n", username, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("[+] Found user: %s\n", username)
	case http.StatusFound:
		if verbose{
			fmt.Printf("[*] User not found: %s\n", username)
		}
	default:
		fmt.Printf("Unexpected status %d for user %s\n", resp.StatusCode, username)
	}
}
