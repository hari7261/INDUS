package commands

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"indus/internal/cli"
	"indus/internal/config"
	"indus/internal/httpclient"
)

type HTTP struct {
	cfg    *config.Config
	client *httpclient.Client
}

func NewHTTP(cfg *config.Config) *HTTP {
	return &HTTP{
		cfg: cfg,
		client: httpclient.New(httpclient.Options{
			MaxRetries: cfg.MaxRetries,
		}),
	}
}

func (c *HTTP) Name() string        { return "http" }
func (c *HTTP) Description() string { return "Make HTTP requests (get, post, put, delete)" }

// Run satisfies cli.Command — delegates to RunStream with real stdio.
func (c *HTTP) Run(ctx context.Context, args []string) error {
	return c.RunStream(ctx, args, os.Stdin, os.Stdout)
}

// RunStream satisfies cli.StreamCommand.
// Response bodies are written to out so the command can feed a pipeline.
// For POST/PUT, if no data arg is supplied the request body is read from in
// (useful when piping: some-cmd | http post https://...).
func (c *HTTP) RunStream(ctx context.Context, args []string, in io.Reader, out io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: http <method> <url> [data]")
		fmt.Fprintln(os.Stderr, "Methods: get, post, put, delete")
		fmt.Fprintln(os.Stderr, "Example: http get https://api.github.com")
		return &cli.UserError{Msg: "missing method"}
	}

	method := strings.ToUpper(args[0])
	switch method {
	case "GET":
		return c.get(ctx, args[1:], out)
	case "POST":
		return c.post(ctx, args[1:], in, out)
	case "PUT":
		return c.put(ctx, args[1:], in, out)
	case "DELETE":
		return c.delete(ctx, args[1:], out)
	default:
		return &cli.UserError{Msg: fmt.Sprintf("unknown method: %s", method)}
	}
}

func (c *HTTP) get(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("http get", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	headers := fs.String("headers", "", "Headers in format 'Key:Value,Key2:Value2'")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: http get <url> [--headers 'Key:Value']")
		return &cli.UserError{Msg: "missing url"}
	}
	url := fs.Arg(0)
	fmt.Fprintf(os.Stderr, "GET %s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &cli.InternalError{Msg: "failed to create request", Err: err}
	}
	if *headers != "" {
		c.parseHeaders(req, *headers)
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return &cli.InternalError{Msg: "request failed", Err: err}
	}
	defer resp.Body.Close()
	fmt.Fprintf(os.Stderr, "Status: %s\n", resp.Status)
	fmt.Fprintln(os.Stderr, "---")

	// Stream directly — avoids buffering the whole body in memory.
	if _, err := io.Copy(out, resp.Body); err != nil {
		return &cli.InternalError{Msg: "failed to stream response", Err: err}
	}
	fmt.Fprintln(out) // trailing newline
	return nil
}

func (c *HTTP) post(ctx context.Context, args []string, in io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("http post", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	headers := fs.String("headers", "", "Headers in format 'Key:Value,Key2:Value2'")
	if err := fs.Parse(args); err != nil {
		return err
	}
	url := fs.Arg(0)
	if url == "" {
		fmt.Fprintln(os.Stderr, "Usage: http post <url> [<data>] [--headers 'Key:Value']")
		return &cli.UserError{Msg: "missing url"}
	}

	// Body: inline arg if given, otherwise read from in (pipeline).
	var body io.Reader
	if fs.NArg() >= 2 {
		body = strings.NewReader(fs.Arg(1))
	} else {
		body = in
	}

	fmt.Fprintf(os.Stderr, "POST %s\n", url)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return &cli.InternalError{Msg: "failed to create request", Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	if *headers != "" {
		c.parseHeaders(req, *headers)
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return &cli.InternalError{Msg: "request failed", Err: err}
	}
	defer resp.Body.Close()
	fmt.Fprintf(os.Stderr, "Status: %s\n", resp.Status)
	fmt.Fprintln(os.Stderr, "---")
	if _, err := io.Copy(out, resp.Body); err != nil {
		return &cli.InternalError{Msg: "failed to stream response", Err: err}
	}
	fmt.Fprintln(out)
	return nil
}

func (c *HTTP) put(ctx context.Context, args []string, in io.Reader, out io.Writer) error {
	fs := flag.NewFlagSet("http put", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	url := fs.Arg(0)
	if url == "" {
		fmt.Fprintln(os.Stderr, "Usage: http put <url> [<data>]")
		return &cli.UserError{Msg: "missing url"}
	}

	var body io.Reader
	if fs.NArg() >= 2 {
		body = strings.NewReader(fs.Arg(1))
	} else {
		body = in
	}

	fmt.Fprintf(os.Stderr, "PUT %s\n", url)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return &cli.InternalError{Msg: "failed to create request", Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return &cli.InternalError{Msg: "request failed", Err: err}
	}
	defer resp.Body.Close()
	fmt.Fprintf(os.Stderr, "Status: %s\n", resp.Status)
	fmt.Fprintln(os.Stderr, "---")
	if _, err := io.Copy(out, resp.Body); err != nil {
		return &cli.InternalError{Msg: "failed to stream response", Err: err}
	}
	fmt.Fprintln(out)
	return nil
}

func (c *HTTP) delete(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("http delete", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: http delete <url>")
		return &cli.UserError{Msg: "missing url"}
	}
	url := fs.Arg(0)
	fmt.Fprintf(os.Stderr, "DELETE %s\n", url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return &cli.InternalError{Msg: "failed to create request", Err: err}
	}
	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return &cli.InternalError{Msg: "request failed", Err: err}
	}
	defer resp.Body.Close()
	fmt.Fprintf(os.Stderr, "Status: %s\n", resp.Status)
	// Stream any response body (some DELETE endpoints return JSON).
	_, _ = io.Copy(out, resp.Body)
	return nil
}

func (c *HTTP) parseHeaders(req *http.Request, headers string) {
	for _, pair := range strings.Split(headers, ",") {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			req.Header.Set(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
		}
	}
}
