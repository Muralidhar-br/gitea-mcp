package label

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
	ListRepoLabelsToolName     = "list_repo_labels"
	GetRepoLabelToolName       = "get_repo_label"
	CreateRepoLabelToolName    = "create_repo_label"
	EditRepoLabelToolName      = "edit_repo_label"
	DeleteRepoLabelToolName    = "delete_repo_label"
	AddIssueLabelsToolName     = "add_issue_labels"
	ReplaceIssueLabelsToolName = "replace_issue_labels"
	ClearIssueLabelsToolName   = "clear_issue_labels"
	RemoveIssueLabelToolName   = "remove_issue_label"
)

var (
	ListRepoLabelsTool = mcp.NewTool(
		ListRepoLabelsToolName,
		mcp.WithDescription("Lists all labels for a given repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(100)),
	)

	GetRepoLabelTool = mcp.NewTool(
		GetRepoLabelToolName,
		mcp.WithDescription("Gets a single label by its ID for a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("label ID")),
	)

	CreateRepoLabelTool = mcp.NewTool(
		CreateRepoLabelToolName,
		mcp.WithDescription("Creates a new label for a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("label name")),
		mcp.WithString("color", mcp.Required(), mcp.Description("label color (hex code, e.g., #RRGGBB)")),
		mcp.WithString("description", mcp.Description("label description")),
	)

	EditRepoLabelTool = mcp.NewTool(
		EditRepoLabelToolName,
		mcp.WithDescription("Edits an existing label in a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("label ID")),
		mcp.WithString("name", mcp.Description("new label name")),
		mcp.WithString("color", mcp.Description("new label color (hex code, e.g., #RRGGBB)")),
		mcp.WithString("description", mcp.Description("new label description")),
	)

	DeleteRepoLabelTool = mcp.NewTool(
		DeleteRepoLabelToolName,
		mcp.WithDescription("Deletes a label from a repository"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("label ID")),
	)

	AddIssueLabelsTool = mcp.NewTool(
		AddIssueLabelsToolName,
		mcp.WithDescription("Adds one or more labels to an issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("issue index")),
		mcp.WithArray("labels", mcp.Required(), mcp.Description("array of label IDs to add"), mcp.Items(map[string]interface{}{"type": "number"})),
	)

	ReplaceIssueLabelsTool = mcp.NewTool(
		ReplaceIssueLabelsToolName,
		mcp.WithDescription("Replaces all labels on an issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("issue index")),
		mcp.WithArray("labels", mcp.Required(), mcp.Description("array of label IDs to replace with"), mcp.Items(map[string]interface{}{"type": "number"})),
	)

	ClearIssueLabelsTool = mcp.NewTool(
		ClearIssueLabelsToolName,
		mcp.WithDescription("Removes all labels from an issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("issue index")),
	)

	RemoveIssueLabelTool = mcp.NewTool(
		RemoveIssueLabelToolName,
		mcp.WithDescription("Removes a single label from an issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("issue index")),
		mcp.WithNumber("label_id", mcp.Required(), mcp.Description("label ID to remove")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListRepoLabelsTool,
		Handler: ListRepoLabelsFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetRepoLabelTool,
		Handler: GetRepoLabelFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateRepoLabelTool,
		Handler: CreateRepoLabelFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditRepoLabelTool,
		Handler: EditRepoLabelFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeleteRepoLabelTool,
		Handler: DeleteRepoLabelFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    AddIssueLabelsTool,
		Handler: AddIssueLabelsFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    ReplaceIssueLabelsTool,
		Handler: ReplaceIssueLabelsFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    ClearIssueLabelsTool,
		Handler: ClearIssueLabelsFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    RemoveIssueLabelTool,
		Handler: RemoveIssueLabelFn,
	})
}

func ListRepoLabelsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoLabelsFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	page, ok := req.GetArguments()["page"].(float64)
	if !ok {
		page = 1
	}
	pageSize, ok := req.GetArguments()["pageSize"].(float64)
	if !ok {
		pageSize = 100
	}

	opt := gitea_sdk.ListLabelsOptions{
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	labels, _, err := gitea.Client().ListRepoLabels(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list %v/%v/labels err: %v", owner, repo, err))
	}
	return to.TextResult(labels)
}

func GetRepoLabelFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetRepoLabelFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	id, ok := req.GetArguments()["id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("label ID is required"))
	}

	label, _, err := gitea.Client().GetRepoLabel(owner, repo, int64(id))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/label/%v err: %v", owner, repo, int64(id), err))
	}
	return to.TextResult(label)
}

func CreateRepoLabelFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateRepoLabelFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	name, ok := req.GetArguments()["name"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("name is required"))
	}
	color, ok := req.GetArguments()["color"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("color is required"))
	}
	description, _ := req.GetArguments()["description"].(string) // Optional

	opt := gitea_sdk.CreateLabelOption{
		Name:        name,
		Color:       color,
		Description: description,
	}

	label, _, err := gitea.Client().CreateLabel(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/label err: %v", owner, repo, err))
	}
	return to.TextResult(label)
}

func EditRepoLabelFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditRepoLabelFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	id, ok := req.GetArguments()["id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("label ID is required"))
	}

	opt := gitea_sdk.EditLabelOption{}
	if name, ok := req.GetArguments()["name"].(string); ok {
		opt.Name = ptr.To(name)
	}
	if color, ok := req.GetArguments()["color"].(string); ok {
		opt.Color = ptr.To(color)
	}
	if description, ok := req.GetArguments()["description"].(string); ok {
		opt.Description = ptr.To(description)
	}

	label, _, err := gitea.Client().EditLabel(owner, repo, int64(id), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/label/%v err: %v", owner, repo, int64(id), err))
	}
	return to.TextResult(label)
}

func DeleteRepoLabelFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteRepoLabelFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	id, ok := req.GetArguments()["id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("label ID is required"))
	}

	_, err := gitea.Client().DeleteLabel(owner, repo, int64(id))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete %v/%v/label/%v err: %v", owner, repo, int64(id), err))
	}
	return to.TextResult("Label deleted successfully")
}

func AddIssueLabelsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called AddIssueLabelsFn")
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
		return to.ErrorResult(fmt.Errorf("issue index is required"))
	}
	labelsRaw, ok := req.GetArguments()["labels"].([]interface{})
	if !ok {
		return to.ErrorResult(fmt.Errorf("labels (array of IDs) is required"))
	}
	var labels []int64
	for _, l := range labelsRaw {
		if labelID, ok := l.(float64); ok {
			labels = append(labels, int64(labelID))
		} else {
			return to.ErrorResult(fmt.Errorf("invalid label ID in labels array"))
		}
	}

	opt := gitea_sdk.IssueLabelsOption{
		Labels: labels,
	}

	issueLabels, _, err := gitea.Client().AddIssueLabels(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("add labels to %v/%v/issue/%v err: %v", owner, repo, int64(index), err))
	}
	return to.TextResult(issueLabels)
}

func ReplaceIssueLabelsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ReplaceIssueLabelsFn")
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
		return to.ErrorResult(fmt.Errorf("issue index is required"))
	}
	labelsRaw, ok := req.GetArguments()["labels"].([]interface{})
	if !ok {
		return to.ErrorResult(fmt.Errorf("labels (array of IDs) is required"))
	}
	var labels []int64
	for _, l := range labelsRaw {
		if labelID, ok := l.(float64); ok {
			labels = append(labels, int64(labelID))
		} else {
			return to.ErrorResult(fmt.Errorf("invalid label ID in labels array"))
		}
	}

	opt := gitea_sdk.IssueLabelsOption{
		Labels: labels,
	}

	issueLabels, _, err := gitea.Client().ReplaceIssueLabels(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("replace labels on %v/%v/issue/%v err: %v", owner, repo, int64(index), err))
	}
	return to.TextResult(issueLabels)
}

func ClearIssueLabelsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ClearIssueLabelsFn")
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
		return to.ErrorResult(fmt.Errorf("issue index is required"))
	}

	_, err := gitea.Client().ClearIssueLabels(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("clear labels on %v/%v/issue/%v err: %v", owner, repo, int64(index), err))
	}
	return to.TextResult("Labels cleared successfully")
}

func RemoveIssueLabelFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called RemoveIssueLabelFn")
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
		return to.ErrorResult(fmt.Errorf("issue index is required"))
	}
	labelID, ok := req.GetArguments()["label_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("label ID is required"))
	}

	_, err := gitea.Client().DeleteIssueLabel(owner, repo, int64(index), int64(labelID))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("remove label %v from %v/%v/issue/%v err: %v", int64(labelID), owner, repo, int64(index), err))
	}
	return to.TextResult("Label removed successfully")
}
