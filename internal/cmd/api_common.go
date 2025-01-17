package cmd

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/henvic/httpretty"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stateful/runme/internal/auth"
	"github.com/stateful/runme/internal/client"
	"github.com/stateful/runme/internal/client/graphql/query"
	"github.com/stateful/runme/internal/log"
	"github.com/stateful/runme/internal/version"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	apiURLF      = "api-url"
	authURLF     = "auth-url"
	traceF       = "trace"
	traceAllF    = "trace-all"
	enableChaosF = "enable-chaos"
	apiTokenF    = "api-token"
)

var (
	apiBaseURL  string
	authBaseURL string
	trace       bool
	traceAll    bool
	enableChaos bool
	apiToken    string
)

// TODO(adamb): temporarily we authorize using Github as IdP.
// In the future, we will likely change this to Stateful being IdP.
var defaultAuthURL = func() string {
	ghURL, err := url.Parse(github.Endpoint.AuthURL)
	if err != nil {
		panic(err)
	}
	return (&url.URL{Scheme: ghURL.Scheme, Host: ghURL.Host}).String()
}()

func getAPIURL() string    { return apiBaseURL }
func getAuthURL() string   { return authBaseURL }
func getTrace() bool       { return trace || traceAll }
func getTraceAll() bool    { return traceAll }
func getEnableChaos() bool { return enableChaos }
func getAPIToken() string  { return apiToken }

func setAPIFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(&authBaseURL, authURLF, defaultAuthURL, "Backend URL to authorize you")
	flagSet.StringVar(&apiBaseURL, apiURLF, "https://api.stateful.com", "Backend URL with API")
	flagSet.StringVar(&apiToken, apiTokenF, "", "API token")
	flagSet.BoolVar(&trace, traceF, false, "Trace HTTP calls")
	flagSet.BoolVar(&traceAll, traceAllF, false, "Trace all HTTP calls including authentication (it might leak sensitive data to output)")
	flagSet.BoolVar(&enableChaos, enableChaosF, false, "Enable Chaos Monkey mode for GraphQL requests")

	mustMarkHidden := func(name string) {
		if err := flagSet.MarkHidden(name); err != nil {
			panic(err)
		}
	}

	mustMarkHidden(authURLF)
	mustMarkHidden(apiURLF)
	mustMarkHidden(apiTokenF)
	mustMarkHidden(traceF)
	mustMarkHidden(traceAllF)
	mustMarkHidden(enableChaosF)
}

var (
	authEnv        auth.Env        // overwritten only in unit tests; when nil a default env will be used
	authAuthorizer auth.Authorizer // overwritten only in unit tests
	tokenStorage   = &auth.DiskStorage{Location: getDefaultConfigHome()}
)

// authorizerWithEnv is a decorator that can return a token
// from the environment variables.
type authorizerWithEnv struct {
	auth.Authorizer
}

func (a *authorizerWithEnv) GetToken(ctx context.Context) (string, error) {
	if apiToken := getAPIToken(); apiToken != "" {
		return apiToken, nil
	}
	return a.Authorizer.GetToken(ctx)
}

func newAuth() auth.Authorizer {
	if authAuthorizer != nil {
		return authAuthorizer
	}

	conf := oauthConfig(getAuthURL())
	opts := []auth.Opts{}

	if getTraceAll() {
		opts = append(opts, auth.WithClient(newAuthClient()))
	}

	if authEnv != nil {
		opts = append(opts, auth.WithEnv(authEnv))
	}

	return &authorizerWithEnv{
		Authorizer: auth.New(conf, getAPIURL(), tokenStorage, opts...),
	}
}

func newAPIClient(ctx context.Context) *http.Client {
	opts := []client.Option{
		client.WithTokenGetter(func() (string, error) {
			return newAuth().GetToken(ctx)
		}),
		client.WithUserAgent(version.BuildVersion),
	}
	if getTrace() {
		opts = append(opts, httpLoggerMiddleware(os.Stderr))
	}
	if l := log.Get(); l != nil {
		opts = append(opts, client.WithLogger(l.Named("APIClient")))
	}
	if getEnableChaos() {
		log.Get().Debug("enabling chaos monkey")
		opts = append(opts, client.WithChaosMonkey(0.1, 0.1))
	}
	return client.NewHTTPClient(nil, opts...)
}

func oauthConfig(authBaseURL string) oauth2.Config {
	return oauth2.Config{
		ClientID:     "bf568e40cfbd1c1261a9",
		ClientSecret: "0de10314c28b754d0cedbf34d081c990865e1363",
		Scopes:       []string{"read:user", "user:email"},
		Endpoint: oauth2.Endpoint{
			// These URLs are modeleted after Github API.
			// If we want to switch between various IdP,
			// we should allow more verbose configuration.
			AuthURL:  authBaseURL + "/login/oauth/authorize",
			TokenURL: authBaseURL + "/login/oauth/access_token",
		},
	}
}

func newAuthClient() *http.Client {
	opts := []client.Option{
		client.WithUserAgent(version.BuildVersion),
		client.WithContentType("application/json"),
	}
	if getTraceAll() {
		opts = append(opts, httpLoggerMiddleware(os.Stderr))
	}
	return client.NewHTTPClient(nil, opts...)
}

func graphqlEndpoint() string {
	return getAPIURL() + "/graphql"
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func httpLoggerMiddleware(out io.Writer) func(http.RoundTripper) http.RoundTripper {
	logger := &httpretty.Logger{
		Time:            true,
		TLS:             false,
		Colors:          isTerminal(os.Stderr.Fd()),
		RequestHeader:   true,
		RequestBody:     true,
		ResponseHeader:  true,
		ResponseBody:    true,
		Formatters:      []httpretty.Formatter{&httpretty.JSONFormatter{}},
		MaxResponseBody: 50000,
	}
	logger.SetOutput(out)
	return logger.RoundTripper
}

func recoverableWithLogin(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, auth.ErrNotFound)
}

func trackInputFromCmd(cmd *cobra.Command, args []string) query.TrackInput {
	fragments := append([]string{"cli", "command"}, strings.Split(cmd.CommandPath(), " ")...)
	fragments = append(fragments, args...)

	return query.TrackInput{
		Events: []query.TrackEvent{
			{
				Event: strings.Join(fragments, "/"),
			},
		},
	}
}
