package gitea

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"gitea.com/gitea/gitea-mcp/pkg/flag"
	"gitea.com/gitea/gitea-mcp/pkg/log"

	"code.gitea.io/sdk/gitea"
)

var (
	client     *gitea.Client
	clientOnce sync.Once
)

func Client() *gitea.Client {
	clientOnce.Do(func() {
		var err error
		if client != nil {
			return
		}

		httpClient := &http.Client{
			Transport: http.DefaultTransport,
		}

		opts := []gitea.ClientOption{
			gitea.SetToken(flag.Token),
		}
		if flag.Insecure {
			httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		opts = append(opts, gitea.SetHTTPClient(httpClient))
		if flag.Debug {
			opts = append(opts, gitea.SetDebugMode())
		}
		client, err = gitea.NewClient(flag.Host, opts...)
		if err != nil {
			log.Fatalf("create gitea client err: %v", err)
		}

		// Set user agent for the client
		client.SetUserAgent(fmt.Sprintf("gitea-mcp-server/%s", flag.Version))
	})
	return client
}
