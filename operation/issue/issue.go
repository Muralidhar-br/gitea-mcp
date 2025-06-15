package issue

import (
	"context"
	"fmt"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/ptr"
	"gitea.com/gitea/gitea-mcp/pkg/to"
	"gitea.com/gitea/gitea-mcp/pkg/tool"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Tool = tool.New()

const (
	GetIssueByIndexToolName         = "get_issue_by_index"
	ListRepoIssuesToolName          = "list_repo_issues"
	CreateIssueToolName             = "create_issue"
	CreateIssueCommentToolName      = "create_issue_comment"
	EditIssueToolName               = "edit_issue"
	EditIssueCommentToolName        = "edit_issue_comment"
	GetIssueCommentsByIndexToolName = "get_issue_comments_by_index"
)

var (
	GetIssueByIndexTool = mcp.NewTool(
		GetIssueByIndexToolName,
		mcp.WithDescription("get issue by index"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
	)

	ListRepoIssuesTool = mcp.NewTool(
		ListRepoIssuesToolName,
		mcp.WithDescription("List repository issues"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("state", mcp.Description("issue state"), mcp.DefaultString("all")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(100)),
	)

	CreateIssueTool = mcp.NewTool(
		CreateIssueToolName,
		mcp.WithDescription("create issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("issue title")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue body")),
	)

	CreateIssueCommentTool = mcp.NewTool(
		CreateIssueCommentToolName,
		mcp.WithDescription("create issue comment"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue comment body")),
	)

	EditIssueTool = mcp.NewTool(
		EditIssueToolName,
		mcp.WithDescription("edit issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
		mcp.WithString("title", mcp.Description("issue title"), mcp.DefaultString("")),
		mcp.WithString("body", mcp.Description("issue body content")),
		mcp.WithArray("assignees", mcp.Description("usernames to assign to this issue"), mcp.Items(map[string]interface{}{"type": "string"})),
		mcp.WithNumber("milestone", mcp.Description("milestone number")),
		mcp.WithString("state", mcp.Description("issue state, one of open, closed, all")),
	)

	EditIssueCommentTool = mcp.NewTool(
		EditIssueCommentToolName,
		mcp.WithDescription("edit issue comment"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("commentID", mcp.Required(), mcp.Description("id of issue comment")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue comment body")),
	)

	GetIssueCommentsByIndexTool = mcp.NewTool(
		GetIssueCommentsByIndexToolName,
		mcp.WithDescription("get issue comment by index"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetIssueByIndexTool,
		Handler: GetIssueByIndexFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListRepoIssuesTool,
		Handler: ListRepoIssuesFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateIssueTool,
		Handler: CreateIssueFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateIssueCommentTool,
		Handler: CreateIssueCommentFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditIssueTool,
		Handler: EditIssueFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditIssueCommentTool,
		Handler: EditIssueCommentFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetIssueCommentsByIndexTool,
		Handler: GetIssueCommentsByIndexFn,
	})
}

func GetIssueByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetIssueByIndexFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	issue, _, err := gitea.Client().GetIssue(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issue/%v err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(issue)
}

func ListRepoIssuesFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListIssuesFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	state, ok := req.GetArguments()["state"].(string)
	if !ok {
		state = "all"
	}
	page, ok := req.GetArguments()["page"].(float64)
	if !ok {
		page = 1
	}
	pageSize, ok := req.GetArguments()["pageSize"].(float64)
	if !ok {
		pageSize = 100
	}
	opt := gitea_sdk.ListIssueOption{
		State: gitea_sdk.StateType(state),
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	issues, _, err := gitea.Client().ListRepoIssues(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issues err: %v", owner, repo, err))
	}
	return to.TextResult(issues)
}

func CreateIssueFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateIssueFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	title, ok := req.GetArguments()["title"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("title is required"))
	}
	body, ok := req.GetArguments()["body"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("body is required"))
	}
	issue, _, err := gitea.Client().CreateIssue(owner, repo, gitea_sdk.CreateIssueOption{
		Title: title,
		Body:  body,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/issue err: %v", owner, repo, err))
	}

	return to.TextResult(issue)
}

func CreateIssueCommentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateIssueCommentFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	body, ok := req.GetArguments()["body"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("body is required"))
	}
	opt := gitea_sdk.CreateIssueCommentOption{
		Body: body,
	}
	issueComment, _, err := gitea.Client().CreateIssueComment(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/issue/%v/comment err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(issueComment)
}

func EditIssueFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditIssueFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}

	opt := gitea_sdk.EditIssueOption{}

	title, ok := req.GetArguments()["title"].(string)
	if ok {
		opt.Title = title
	}
	body, ok := req.GetArguments()["body"].(string)
	if ok {
		opt.Body = ptr.To(body)
	}
	assignees, ok := req.GetArguments()["assignees"].([]string)
	if ok {
		opt.Assignees = assignees
	}
	milestone, ok := req.GetArguments()["milestone"].(float64)
	if ok {
		opt.Milestone = ptr.To(int64(milestone))
	}
	state, ok := req.GetArguments()["state"].(string)
	if ok {
		opt.State = ptr.To(gitea_sdk.StateType(state))
	}

	issue, _, err := gitea.Client().EditIssue(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/issue/%v err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(issue)
}

func EditIssueCommentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditIssueCommentFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	commentID, ok := req.GetArguments()["commentID"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("comment ID is required"))
	}
	body, ok := req.GetArguments()["body"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("body is required"))
	}
	opt := gitea_sdk.EditIssueCommentOption{
		Body: body,
	}
	issueComment, _, err := gitea.Client().EditIssueComment(owner, repo, int64(commentID), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/issues/comments/%v err: %v", owner, repo, int64(commentID), err))
	}

	return to.TextResult(issueComment)
}

func GetIssueCommentsByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetIssueCommentsByIndexFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	opt := gitea_sdk.ListIssueCommentOptions{}
	issue, _, err := gitea.Client().ListIssueComments(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issues/%v/comments err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(issue)
}
