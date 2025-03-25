# Gitea MCP Server

**Gitea MCP Server** is an integration plugin designed to connect Gitea with Model Context Protocol (MCP) systems. This allows for seamless command execution and repository management through an MCP-compatible chat interface.

## 🚧Installation

### 📥Download the official binary release

You can download the official release from [here](https://gitea.com/gitea/gitea-mcp/releases).

### 🔧Build from Source

You can download the source code by cloning the repository using Git:

```bash
git clone https://gitea.com/gitea/gitea-mcp.git
```

Before building, make sure you have the following installed:

- make
- Golang (Go 1.24 or later recommended)

Then run:

```bash
make build
```

### 📁Add to PATH

After building, copy the binary gitea-mcp to a directory included in your system's PATH. For example:

```bash
cp gitea-mcp /usr/local/bin/
```

## 🚀Usage

This example is for Cursor, you can also use plugins in VSCode.
To configure the MCP server for Gitea, add the following to your MCP configuration file:

- **stdio mode**

```json
{
  "mcpServers": {
    "gitea": {
      "command": "gitea-mcp",
      "args": [
        "-t", "stdio",
        "--host", "https://gitea.com"
        // "--token", "<your personal access token>"
      ],
      "env": {
        // "GITEA_HOST": "https://gitea.com",
        "GITEA_ACCESS_TOKEN": "<your personal access token>"
      }
    }
  }
}
```

- **sse mode**

```json
{
  "mcpServers": {
    "gitea": {
      "url": "http://localhost:8080/sse"
    }
  }
}
```

> [!NOTE]
> You can provide your Gitea host and access token either as command-line arguments or environment variables.
> Command-line arguments have the highest priority

Once everything is set up, try typing the following in your MCP-compatible chatbox:

```text
list all my repositories
```

## ✅Available Tools

The Gitea MCP Server supports the following tools:

|  Tool  |  Scope  | Description  |
|:------:|:-------:|:------------:|
|get_my_user_info|User|Get the information of the authenticated user|
|create_repo|Repository|Create a new repository|
|fork_repo|Repository|Fork a repository|
|list_my_repos|Repository|List all repositories owned by the authenticated user|
|create_branch|Branch|Create a new branch|
|delete_branch|Branch|Delete a branch|
|list_branches|Branch|List all branches in a repository|
|list_repo_commits|Commit|List all commits in a repository|
|get_file_content|File|Get the content and metadata of a file|
|create_file|File|Create a new file|
|update_file|File|Update an existing file|
|delete_file|File|Delete a file|
|get_issue_by_index|Issue|Get an issue by its index|
|list_repo_issues|Issue|List all issues in a repository|
|create_issue|Issue|Create a new issue|
|create_issue_comment|Issue|Create a comment on an issue|
|get_pull_request_by_index|Pull Request|Get a pull request by its index|
|list_repo_pull_requests|Pull Request|List all pull requests in a repository|
|create_pull_request|Pull Request|Create a new pull request|
|search_users|User|Search for users|
|search_org_teams|Organization|Search for teams in an organization|
|search_repos|Repository|Search for repositories|
|get_gitea_mcp_server_version|Server|Get the version of the Gitea MCP Server|

## 🐛Debugging

To enable debug mode, add the `-d` flag when running the Gitea MCP Server with sse mode:

```sh
./gitea-mcp -t sse --token <your personal access token> -d
```

Enjoy exploring and managing your Gitea repositories via chat!
