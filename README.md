# FriendSocial: Scheduled Activities Service

## Overview

FriendSocial is a robust social networking platform that focuses on connecting friends through scheduled activities. This repository contains the backend service responsible for managing scheduled activities, a core feature of the FriendSocial platform.

## Key Features

- Create, read, update, and delete scheduled activities
- Manage recurring activities with complex scheduling patterns
- Handle user availability and time zone differences
- Efficient batch operations for creating multiple scheduled activities
- Integration with user activity preferences

## Technology Stack

- Go 1.23.0
- PostgreSQL (using pgx driver)
- RESTful API design
- Modular architecture for scalability


## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up your PostgreSQL database and update the connection details in the configuration
4. Run the service: `go run main.go`