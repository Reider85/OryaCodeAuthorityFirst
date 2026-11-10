package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	githubAPIBase = "https://api.github.com"
	owner         = "Reider85"
	repo          = "OryaCodeAuthorityFirst"
)

type Usage struct {
	TotalMinutes     int `json:"total_minutes"`
	IncludedMinutes  int `json:"included_minutes"`
	MinutesUsed      int `json:"minutes_used"`
	MinutesRemaining int `json:"minutes_remaining"`
}

type Response struct {
	Repository struct {
		Actions struct {
			TotalBillableMinutesUsed      int `json:"total_billable_minutes_used"`
			TotalPaidMinutesUsed          int `json:"total_paid_minutes_used"`
			TotalMinutesUsed              int `json:"total_minutes_used"`
			TotalNonBillableMinutesUsed   int `json:"total_non_billable_minutes_used"`
			MinutesUsedBreakdown          []struct {
				RunnerLabel string `json:"runner_label"`
				MinutesUsed int    `json:"minutes_used"`
			} `json:"minutes_used_breakdown"`
			EstimatedPaidMinutesUsed      int `json:"estimated_paid_minutes_used"`
			EstimatedPaidMinutesUsedByJob []struct {
				JobID       string `json:"job_id"`
				JobName     string `json:"job_name"`
				MinutesUsed int    `json:"minutes_used"`
			} `json:"estimated_paid_minutes_used_by_job"`
		} `json:"actions"`
	} `json:"repository"`
}

func getGitHubUsage() (*Response, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("%s/repos/%s/%s/actions/runner-groups", githubAPIBase, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func calculateUsageThreshold(response *Response, threshold float64) (bool, *Usage) {
	totalMinutes := response.Repository.Actions.TotalMinutesUsed
	includedMinutes := 6000 // Default free minutes for public repos
	
	minutesRemaining := includedMinutes - totalMinutes
	usagePercentage := float64(totalMinutes) / float64(includedMinutes) * 100
	
	shouldUseLocal := usagePercentage >= threshold
	
	usage := &Usage{
		TotalMinutes:     totalMinutes,
		IncludedMinutes:  includedMinutes,
		MinutesUsed:      totalMinutes,
		MinutesRemaining: minutesRemaining,
	}
	
	return shouldUseLocal, usage
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: runner-manager <threshold_percentage>")
		os.Exit(1)
	}
	
	threshold, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid threshold: %v\n", err)
		os.Exit(1)
	}
	
	response, err := getGitHubUsage()
	if err != nil {
		fmt.Printf("Failed to get GitHub usage: %v\n", err)
		os.Exit(1)
	}
	
	shouldUseLocal, usage := calculateUsageThreshold(response, threshold)
	
	fmt.Printf("GitHub Actions Usage:\n")
	fmt.Printf("  Total Minutes Used: %d\n", usage.TotalMinutes)
	fmt.Printf("  Included Minutes: %d\n", usage.IncludedMinutes)
	fmt.Printf("  Minutes Remaining: %d\n", usage.MinutesRemaining)
	fmt.Printf("  Usage Percentage: %.2f%%\n", float64(usage.MinutesUsed)/float64(usage.IncludedMinutes)*100)
	
	if shouldUseLocal {
		fmt.Printf("🚨 Usage threshold (%.0f%%) exceeded! Switching to local runners.\n", threshold)
		os.Exit(1) // Exit with error code to trigger local runner fallback
	} else {
		fmt.Printf("✅ Usage below threshold (%.0f%%). Using GitHub-hosted runners.\n", threshold)
		os.Exit(0)
	}
}