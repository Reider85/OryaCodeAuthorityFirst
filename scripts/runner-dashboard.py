#!/usr/bin/env python3

"""
Hybrid Runner Dashboard
Monitors GitHub Actions usage and runner status
"""

import requests
import json
import time
import os
from datetime import datetime, timedelta
from typing import Dict, List, Optional

class HybridRunnerDashboard:
    def __init__(self):
        self.github_token = os.getenv('GITHUB_TOKEN')
        self.owner = 'Reider85'
        self.repo = 'OryaCodeAuthorityFirst'
        self.threshold = 80  # Default threshold percentage
        
    def get_github_api_headers(self):
        """Get headers for GitHub API requests"""
        if not self.github_token:
            raise ValueError("GITHUB_TOKEN environment variable not set")
        
        return {
            'Authorization': f'token {self.github_token}',
            'Accept': 'application/vnd.github.v3+json'
        }
    
    def get_repository_usage(self) -> Dict:
        """Get repository usage statistics from GitHub API"""
        url = f"https://api.github.com/repos/{self.owner}/{self.repo}/actions/usage"
        headers = self.get_github_api_headers()
        
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Error fetching repository usage: {e}")
            return {}
    
    def get_runner_groups(self) -> Dict:
        """Get runner groups and their runners"""
        url = f"https://api.github.com/repos/{self.owner}/{self.repo}/actions/runner-groups"
        headers = self.get_github_api_headers()
        
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Error fetching runner groups: {e}")
            return {}
    
    def get_runners(self) -> Dict:
        """Get all runners in the repository"""
        url = f"https://api.github.com/repos/{self.owner}/{self.repo}/actions/runners"
        headers = self.get_github_api_headers()
        
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Error fetching runners: {e}")
            return {}
    
    def calculate_usage_percentage(self, usage_data: Dict) -> float:
        """Calculate usage percentage based on free minutes"""
        if not usage_data:
            return 0.0
        
        total_minutes = usage_data.get('total_minutes_used', 0)
        included_minutes = 6000  # Default free minutes for public repositories
        
        if included_minutes == 0:
            return 0.0
        
        return (total_minutes / included_minutes) * 100
    
    def should_use_local_runners(self, usage_percentage: float) -> bool:
        """Determine if local runners should be used based on threshold"""
        return usage_percentage >= self.threshold
    
    def format_runner_status(self, runner: Dict) -> str:
        """Format runner status for display"""
        name = runner.get('name', 'Unknown')
        os_type = runner.get('os', 'Unknown')
        status = runner.get('status', 'Unknown')
        busy = runner.get('busy', False)
        
        status_icon = "🟢" if status == "online" else "🔴"
        busy_icon = "🏃‍♂️" if busy else "💤"
        
        return f"{status_icon} {name} ({os_type}) {busy_icon}"
    
    def display_dashboard(self):
        """Display the dashboard"""
        print("🚀 Hybrid Runner Dashboard")
        print("=" * 50)
        print(f"Repository: {self.owner}/{self.repo}")
        print(f"Threshold: {self.threshold}%")
        print(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print()
        
        # Get usage data
        usage_data = self.get_repository_usage()
        if usage_data:
            usage_percentage = self.calculate_usage_percentage(usage_data)
            total_minutes = usage_data.get('total_minutes_used', 0)
            included_minutes = 6000
            remaining_minutes = max(0, included_minutes - total_minutes)
            
            print("📊 GitHub Actions Usage:")
            print(f"  Total Minutes Used: {total_minutes}")
            print(f"  Free Minutes: {included_minutes}")
            print(f"  Minutes Remaining: {remaining_minutes}")
            print(f"  Usage Percentage: {usage_percentage:.1f}%")
            
            # Determine runner strategy
            should_use_local = self.should_use_local_runners(usage_percentage)
            strategy = "🏃‍♂️ Use Local Runners" if should_use_local else "☁️ Use GitHub Hosted"
            print(f"  Strategy: {strategy}")
            print()
        
        # Get runner groups
        runner_groups = self.get_runner_groups()
        if runner_groups and runner_groups.get('runner_groups'):
            print("🏃‍♂️ Runner Groups:")
            for group in runner_groups['runner_groups']:
                group_name = group.get('name', 'Unknown')
                runner_count = group.get('runners_count', 0)
                print(f"  {group_name}: {runner_count} runners")
            print()
        
        # Get runners
        runners = self.get_runners()
        if runners and runners.get('runners'):
            print("🏃‍♂️ Active Runners:")
            for runner in runners['runners']:
                if runner.get('status') == 'online':
                    print(f"  {self.format_runner_status(runner)}")
            print()
        
        # Recommendations
        print("💡 Recommendations:")
        if usage_percentage >= 90:
            print("  ⚠️  Critical: Free minutes almost exhausted!")
            print("  🚀 Switch to local runners immediately")
        elif usage_percentage >= self.threshold:
            print("  ⚠️  Warning: Usage approaching threshold")
            print("  🔄 Consider switching to local runners")
        else:
            print("  ✅ Usage within acceptable limits")
            print("  ☁️ Continue using GitHub hosted runners")

def main():
    """Main function"""
    dashboard = HybridRunnerDashboard()
    
    try:
        dashboard.display_dashboard()
    except ValueError as e:
        print(f"Error: {e}")
        print("Please set the GITHUB_TOKEN environment variable")
        return 1
    except Exception as e:
        print(f"Unexpected error: {e}")
        return 1
    
    return 0

if __name__ == "__main__":
    exit(main())