# ZibraGo - CLI Backup Tool

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)
![CLI](https://img.shields.io/badge/CLI-Tool-green?style=for-the-badge)

A powerful command-line backup tool written in Go that automatically archives directories and uploads them to Google Drive.

 **Features**

- **Smart Directory Scanning** - Recursively scans directories for files
- **Parallel Zip Archiving** - Multi-threaded compression for optimal performance
- **Google Drive Integration** - Seamless cloud backup functionality
- **Configurable Settings** - Easy configuration management
- **Structured Logging** - JSON logging with different levels
- **High Performance** - Utilizes all available CPU cores

##  Quick Start

### Prerequisites

- Go 1.21 or higher
- Google Cloud Project with Drive API enabled

### Installation

```bash
# Clone the repository
git clone https://github.com/Talos-hub/ZibraGo.git
cd ZibraGo

# Build the application
go build -o zibrago

# Install globally
go install
```
### Google Drive Setup
- **Enable Google Drive API**

- **Go to Google Cloud Console**

- **Create a new project or select existing one**

- **Enable the Google Drive API**

- **Create OAuth 2.0 credentials (Desktop application type)**

- **Configure Credentials**

- **Download the credentials.json file**

- **Place it in your ZibraGo working directory**
### First Time Setup
```
# Configure backup directory
./zibrago settings

# Follow the prompts to set your backup path
Enter path for zip files: /path/to/your/backup/directory
```

## Usage
### Basic Backup
```
# Backup a directory to Google Drive
./zibrago start mybackup.zip /path/to/directory

# Example: Backup your documents folder
./zibrago start documents_backup.zip ~/Documents
```
### View Logs
```
# Check application logs
./zibrago log
```
### Manage Settings
```
# Change backup storage location
./zibrago settings
```

## Configuration
ZibraGo uses a JSON configuration file (zipdir.json) to manage settings:
```
{
  "zipdir": "/path/to/your/backup/directory"
}
```

## Troubleshooting
### Common Issues
#### "credentials.json not found"

- **Ensure you've downloaded the OAuth 2.0 credentials file**

- **Place it in the same directory as the executable**

#### "Log file not exists"

- **Logs are created automatically during first run**

- **This message is normal if no operations have been performed yet**

#### Google Drive authentication errors

- **Delete token.json and restart the application**

- **Re-authenticate when prompted**
