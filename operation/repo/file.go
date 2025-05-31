package repo

import (
	"context"
	"encoding/base64"
	"fmt"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/to"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	GetFileToolName    = "get_file_content"
	GetDirToolName     = "get_dir_content"
	CreateFileToolName = "create_file"
	UpdateFileToolName = "update_file"
	DeleteFileToolName = "delete_file"
)

var (
	GetFileContentTool = mcp.NewTool(
		GetFileToolName,
		mcp.WithDescription("Get file Content and Metadata"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("ref", mcp.Required(), mcp.Description("ref can be branch/tag/commit")),
		mcp.WithString("filePath", mcp.Required(), mcp.Description("file path")),
	)

	GetDirContentTool = mcp.NewTool(
		GetDirToolName,
		mcp.WithDescription("Get a list of entries in a directory"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("ref", mcp.Required(), mcp.Description("ref can be branch/tag/commit")),
		mcp.WithString("filePath", mcp.Required(), mcp.Description("directory path")),
	)

	CreateFileTool = mcp.NewTool(
		CreateFileToolName,
		mcp.WithDescription("Create file"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("filePath", mcp.Required(), mcp.Description("file path")),
		mcp.WithString("content", mcp.Required(), mcp.Description("file content")),
		mcp.WithString("message", mcp.Required(), mcp.Description("commit message")),
		mcp.WithString("branch_name", mcp.Required(), mcp.Description("branch name")),
		mcp.WithString("new_branch_name", mcp.Description("new branch name")),
	)

	UpdateFileTool = mcp.NewTool(
		UpdateFileToolName,
		mcp.WithDescription("Update file"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("filePath", mcp.Required(), mcp.Description("file path")),
		mcp.WithString("sha", mcp.Required(), mcp.Description("sha is the SHA for the file that already exists")),
		mcp.WithString("content", mcp.Required(), mcp.Description("file content, base64 encoded")),
		mcp.WithString("message", mcp.Required(), mcp.Description("commit message")),
		mcp.WithString("branch_name", mcp.Required(), mcp.Description("branch name")),
	)

	DeleteFileTool = mcp.NewTool(
		DeleteFileToolName,
		mcp.WithDescription("Delete file"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("filePath", mcp.Required(), mcp.Description("file path")),
		mcp.WithString("message", mcp.Required(), mcp.Description("commit message")),
		mcp.WithString("branch_name", mcp.Required(), mcp.Description("branch name")),
		mcp.WithString("sha", mcp.Description("sha")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetFileContentTool,
		Handler: GetFileContentFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetDirContentTool,
		Handler: GetDirContentFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateFileTool,
		Handler: CreateFileFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    UpdateFileTool,
		Handler: UpdateFileFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeleteFileTool,
		Handler: DeleteFileFn,
	})
}

func GetFileContentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetFileFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	ref, _ := req.GetArguments()["ref"].(string)
	filePath, ok := req.GetArguments()["filePath"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("filePath is required"))
	}
	content, _, err := gitea.Client().GetContents(owner, repo, ref, filePath)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get file err: %v", err))
	}
	return to.TextResult(content)
}

func GetDirContentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetDirContentFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	ref, _ := req.GetArguments()["ref"].(string)
	filePath, ok := req.GetArguments()["filePath"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("filePath is required"))
	}
	content, _, err := gitea.Client().ListContents(owner, repo, ref, filePath)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get dir content err: %v", err))
	}
	return to.TextResult(content)
}

func CreateFileFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateFileFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	filePath, ok := req.GetArguments()["filePath"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("filePath is required"))
	}
	content, _ := req.GetArguments()["content"].(string)
	message, _ := req.GetArguments()["message"].(string)
	branchName, _ := req.GetArguments()["branch_name"].(string)
	opt := gitea_sdk.CreateFileOptions{
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		FileOptions: gitea_sdk.FileOptions{
			Message:    message,
			BranchName: branchName,
		},
	}

	_, _, err := gitea.Client().CreateFile(owner, repo, filePath, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create file err: %v", err))
	}
	return to.TextResult("Create file success")
}

func UpdateFileFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpdateFileFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	filePath, ok := req.GetArguments()["filePath"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("filePath is required"))
	}
	sha, ok := req.GetArguments()["sha"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("sha is required"))
	}
	content, _ := req.GetArguments()["content"].(string)
	message, _ := req.GetArguments()["message"].(string)
	branchName, _ := req.GetArguments()["branch_name"].(string)

	opt := gitea_sdk.UpdateFileOptions{
		SHA:     sha,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		FileOptions: gitea_sdk.FileOptions{
			Message:    message,
			BranchName: branchName,
		},
	}
	_, _, err := gitea.Client().UpdateFile(owner, repo, filePath, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("update file err: %v", err))
	}
	return to.TextResult("Update file success")
}

func DeleteFileFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteFileFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	filePath, ok := req.GetArguments()["filePath"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("filePath is required"))
	}
	message, _ := req.GetArguments()["message"].(string)
	branchName, _ := req.GetArguments()["branch_name"].(string)
	sha, ok := req.GetArguments()["sha"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("sha is required"))
	}
	opt := gitea_sdk.DeleteFileOptions{
		FileOptions: gitea_sdk.FileOptions{
			Message:    message,
			BranchName: branchName,
		},
		SHA: sha,
	}
	_, err := gitea.Client().DeleteFile(owner, repo, filePath, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete file err: %v", err))
	}
	return to.TextResult("Delete file success")
}
