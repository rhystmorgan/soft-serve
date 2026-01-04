# User Management Guide

This guide covers how to manage users and repository access in Soft Serve Git server.

## Overview

Soft Serve uses SSH-based user management with fine-grained repository access control. Users are managed through SSH commands, and access is controlled through a collaborator system that allows per-repository permissions.

## User Types

### Admin Users
- Can manage all users and repositories
- Can create/delete repositories
- Can manage server settings
- Set during initial setup or promoted later

### Regular Users
- Can only access repositories they're explicitly granted access to
- Cannot create repositories (unless given specific permissions)
- Cannot manage other users

## Creating Users

### Basic User Creation
```bash
# Create a new user
ssh your-domain.com -p 23231 user create USERNAME

# Create a user and add their SSH key in one command
ssh your-domain.com -p 23231 user create USERNAME --key "ssh-ed25519 AAAAC3NzaC1... user-key-here"
```

### Adding SSH Keys
```bash
# Add a public key to an existing user
ssh your-domain.com -p 23231 user add-pubkey USERNAME "ssh-ed25519 AAAAC3NzaC1... user-key-here"

# Remove a public key from a user
ssh your-domain.com -p 23231 user remove-pubkey USERNAME "ssh-ed25519 AAAAC3NzaC1... user-key-here"
```

### User Management Commands
```bash
# List all users
ssh your-domain.com -p 23231 user list

# Get detailed user information
ssh your-domain.com -p 23231 user info USERNAME

# Delete a user
ssh your-domain.com -p 23231 user delete USERNAME

# Change a user's username
ssh your-domain.com -p 23231 user set-username OLD_USERNAME NEW_USERNAME

# Make a user an admin
ssh your-domain.com -p 23231 user set-admin USERNAME true

# Remove admin privileges
ssh your-domain.com -p 23231 user set-admin USERNAME false
```

## Repository Access Control

### Access Levels

Soft Serve provides four access levels for repository collaborators:

- **`no-access`** - User has no access to the repository
- **`read-only`** - User can clone and pull, but cannot push
- **`read-write`** - User can clone, pull, and push changes
- **`admin-access`** - User has full repository admin rights (manage collaborators, settings)

### Managing Collaborators

#### Adding Collaborators
```bash
# Add user with read-only access
ssh your-domain.com -p 23231 collab add REPOSITORY USERNAME read-only

# Add user with read-write access (default if level not specified)
ssh your-domain.com -p 23231 collab add REPOSITORY USERNAME read-write

# Add user with admin access to the repository
ssh your-domain.com -p 23231 collab add REPOSITORY USERNAME admin-access
```

#### Managing Collaborators
```bash
# List all collaborators for a repository
ssh your-domain.com -p 23231 collab list REPOSITORY

# Remove a collaborator from a repository
ssh your-domain.com -p 23231 collab remove REPOSITORY USERNAME
```

## Common Use Cases

### Case 1: Limited User with Single Repository Access

Create a user who can only access one specific repository:

```bash
# 1. Create the user
ssh your-domain.com -p 23231 user create limited-user

# 2. Add their SSH public key
ssh your-domain.com -p 23231 user add-pubkey limited-user "ssh-ed25519 AAAAC3NzaC1... their-key-here"

# 3. Ensure the repository exists
ssh your-domain.com -p 23231 repo create my-project

# 4. Grant access to the specific repository
ssh your-domain.com -p 23231 collab add my-project limited-user read-write

# 5. Verify the setup
ssh your-domain.com -p 23231 collab list my-project
```

The user can now:
- Clone: `git clone ssh://your-domain.com:23231/my-project`
- Only access `my-project` - cannot see other repositories
- Cannot create new repositories
- Cannot manage other users

### Case 2: Team Member with Multiple Repository Access

Create a team member who needs access to several repositories:

```bash
# 1. Create the user
ssh your-domain.com -p 23231 user create team-member

# 2. Add their SSH key
ssh your-domain.com -p 23231 user add-pubkey team-member "ssh-ed25519 AAAAC3NzaC1... their-key-here"

# 3. Grant access to multiple repositories
ssh your-domain.com -p 23231 collab add frontend-app team-member read-write
ssh your-domain.com -p 23231 collab add backend-api team-member read-write
ssh your-domain.com -p 23231 collab add documentation team-member admin-access
```

### Case 3: Read-Only Access for External Collaborators

Grant read-only access to external contributors:

```bash
# 1. Create the user
ssh your-domain.com -p 23231 user create external-contributor

# 2. Add their SSH key
ssh your-domain.com -p 23231 user add-pubkey external-contributor "ssh-ed25519 AAAAC3NzaC1... their-key-here"

# 3. Grant read-only access
ssh your-domain.com -p 23231 collab add open-source-project external-contributor read-only
```

## Security Best Practices

### SSH Key Management
- **Use Ed25519 keys**: Prefer `ssh-ed25519` over RSA for better security
- **Unique keys per user**: Each user should have their own SSH key pair
- **Regular key rotation**: Periodically update SSH keys
- **Remove unused keys**: Clean up keys for departed users

### Access Control
- **Principle of least privilege**: Grant minimum necessary access
- **Regular access reviews**: Periodically review who has access to what
- **Remove unused accounts**: Delete accounts for departed users
- **Monitor access logs**: Check server logs for unusual activity

### User Account Hygiene
```bash
# Regular maintenance commands

# List all users to review accounts
ssh your-domain.com -p 23231 user list

# Check specific user details
ssh your-domain.com -p 23231 user info USERNAME

# Review repository collaborators
ssh your-domain.com -p 23231 collab list REPOSITORY

# Remove departed users
ssh your-domain.com -p 23231 user delete DEPARTED_USER
```

## Troubleshooting

### Common Issues

**User can't authenticate:**
```bash
# Check if user exists
ssh your-domain.com -p 23231 user info USERNAME

# Verify their SSH key is added
ssh your-domain.com -p 23231 user info USERNAME

# Test SSH connection
ssh -T your-domain.com -p 23231
```

**User can't access repository:**
```bash
# Check if user is a collaborator
ssh your-domain.com -p 23231 collab list REPOSITORY

# Verify repository exists
ssh your-domain.com -p 23231 repo list

# Check user's access level
ssh your-domain.com -p 23231 collab list REPOSITORY | grep USERNAME
```

**Permission denied errors:**
- Ensure the user has the correct access level for the operation
- Verify SSH key is correctly formatted and added
- Check that the repository name is correct (case-sensitive)

### Useful Diagnostic Commands
```bash
# Check your own access and identity
ssh your-domain.com -p 23231 info

# List all repositories you can see
ssh your-domain.com -p 23231 repo list

# Test repository access
ssh your-domain.com -p 23231 repo tree REPOSITORY
```

## Integration with Git Workflows

### Clone URLs
Users access repositories using SSH URLs:
```bash
# Standard clone
git clone ssh://your-domain.com:23231/repository-name

# With specific user (if needed)
git clone ssh://username@your-domain.com:23231/repository-name
```

### Setting Up Git Remotes
```bash
# Add Soft Serve as a remote
git remote add origin ssh://your-domain.com:23231/my-project

# Push to Soft Serve
git push origin main
```

### Working with Multiple Repositories
Users with access to multiple repositories can clone and work with each one independently:
```bash
git clone ssh://your-domain.com:23231/project-a
git clone ssh://your-domain.com:23231/project-b
git clone ssh://your-domain.com:23231/project-c
```

## Automation and Scripting

### Bulk User Creation
```bash
#!/bin/bash
# Script to create multiple users from a list

USERS_FILE="users.txt"  # Format: username ssh-key-content
SERVER="your-domain.com:23231"

while IFS=' ' read -r username ssh_key; do
    echo "Creating user: $username"
    ssh $SERVER user create "$username" --key "$ssh_key"
done < "$USERS_FILE"
```

### Repository Access Audit
```bash
#!/bin/bash
# Script to audit repository access

SERVER="your-domain.com:23231"

# Get all repositories
repos=$(ssh $SERVER repo list)

for repo in $repos; do
    echo "=== Repository: $repo ==="
    ssh $SERVER collab list "$repo"
    echo
done
```

---

This user management system provides fine-grained control over repository access while maintaining security and simplicity. Regular maintenance and following security best practices will ensure your Soft Serve installation remains secure and well-organized.