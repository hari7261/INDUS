package engine

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"indus/internal/httpclient"
)

type networkModule struct {
	engine *Engine
	client *httpclient.Client
}

func (m *networkModule) Execute(_ context.Context, inv Invocation) Response {
	if m.client == nil {
		m.client = httpclient.New(httpclient.Options{
			Timeout:    time.Duration(m.engine.cfg.APITimeout) * time.Second,
			MaxRetries: m.engine.cfg.MaxRetries,
		})
	}

	switch inv.Path {
	case "net scan":
		return m.scan(inv)
	case "net pingx":
		return m.pingx(inv)
	case "net trace":
		return m.trace(inv)
	case "net ports":
		return m.ports(inv)
	case "net status":
		return m.status(inv)
	case "net fetch":
		return m.fetch(inv)
	default:
		return Response{Err: unknownCommandError(inv.Command)}
	}
}

func (m *networkModule) scan(inv Invocation) Response {
	if len(inv.Parsed.Positionals) > 0 {
		target := inv.Parsed.Positionals[0]
		addresses, err := net.LookupHost(target)
		if err != nil {
			return Response{Err: commandFailedError(inv.Command, err)}
		}
		return Response{Output: fmt.Sprintf("target=%s\naddresses=%s", target, strings.Join(addresses, ","))}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		fmt.Fprintf(buffer, "%s %s\n", iface.Name, strings.Join(stringifyAddrs(addrs), ","))
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *networkModule) pingx(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind net pingx <host> [--port 443]")}
	}
	host := inv.Parsed.Positionals[0]
	port, err := inv.Parsed.Int(443, "port")
	if err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --port value")}
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	_ = conn.Close()
	return Response{Output: fmt.Sprintf("target=%s\nlatency=%s", address, time.Since(start))}
}

func (m *networkModule) trace(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind net trace <host>")}
	}

	host := inv.Parsed.Positionals[0]
	ips, err := net.LookupHost(host)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	fmt.Fprintf(buffer, "trace_mode=soft\nhost=%s\n", host)
	for _, ip := range ips {
		fmt.Fprintf(buffer, "resolved=%s\n", ip)
		for _, port := range []int{443, 80} {
			address := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
			start := time.Now()
			conn, err := net.DialTimeout("tcp", address, 800*time.Millisecond)
			if err != nil {
				fmt.Fprintf(buffer, "dial=%s status=unreachable\n", address)
				continue
			}
			_ = conn.Close()
			fmt.Fprintf(buffer, "dial=%s latency=%s\n", address, time.Since(start))
		}
	}

	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *networkModule) ports(inv Invocation) Response {
	from, err := inv.Parsed.Int(3000, "from")
	if err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --from value")}
	}
	to, err := inv.Parsed.Int(3010, "to")
	if err != nil || to < from {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --to value")}
	}

	buffer := m.engine.getBuffer()
	defer m.engine.putBuffer(buffer)

	openPorts := 0
	for port := from; port <= to; port++ {
		address := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", address, 50*time.Millisecond)
		if err == nil {
			openPorts++
			fmt.Fprintf(buffer, "open=%d\n", port)
			_ = conn.Close()
		}
	}

	if openPorts == 0 {
		return Response{Output: "open_ports=0"}
	}
	return Response{Output: strings.TrimSpace(buffer.String())}
}

func (m *networkModule) status(inv Invocation) Response {
	target := inv.Parsed.String("url")
	if target == "" {
		target = "https://example.com"
	}

	parsed, err := url.Parse(target)
	if err != nil {
		return Response{Err: invalidArgumentError(inv.Command, "invalid --url value")}
	}

	ips, err := net.LookupHost(parsed.Hostname())
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	req, err := http.NewRequest(http.MethodHead, target, nil)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	resp, err := m.client.Do(context.Background(), req)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	_ = resp.Body.Close()

	return Response{Output: fmt.Sprintf("url=%s\ndns=%s\nstatus=%s", target, strings.Join(ips, ","), resp.Status)}
}

func (m *networkModule) fetch(inv Invocation) Response {
	if len(inv.Parsed.Positionals) == 0 {
		return Response{Err: invalidArgumentError(inv.Command, "usage: ind net fetch <url> [--method GET] [--body data]")}
	}

	target := inv.Parsed.Positionals[0]
	method := strings.ToUpper(inv.Parsed.String("method"))
	if method == "" {
		method = http.MethodGet
	}
	body := inv.Parsed.String("body")

	req, err := http.NewRequest(method, target, strings.NewReader(body))
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.client.Do(context.Background(), req)
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return Response{Err: commandFailedError(inv.Command, err)}
	}

	return Response{Output: fmt.Sprintf("status=%s\nbody=%s", resp.Status, strings.TrimSpace(string(data)))}
}

func stringifyAddrs(addrs []net.Addr) []string {
	values := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		values = append(values, addr.String())
	}
	return values
}
