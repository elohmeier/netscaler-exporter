package collector

import (
	"context"
	"net"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/elohmeier/netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

const hostHeaderRewriteTarget = `HTTP.REQ.HEADER("Host")`

// httpHostRewriteMapping is a resolved, static request-side Host header rewrite.
type httpHostRewriteMapping struct {
	VirtualServer string
	Policy        string
	Priority      string
	Rule          string
	Action        string
	Host          string
}

func (e *Exporter) collectHTTPHostRewriteInfo(ctx context.Context, nsClient *netscaler.NitroClient, ch chan<- prometheus.Metric) {
	e.lbVServerHTTPHostRewriteInfo.Reset()

	bindings, err := netscaler.GetAllLBVServerRewritePolicyBindings(ctx, nsClient)
	if err != nil {
		e.logger.Error("failed to get LB vserver rewrite policy bindings", "url", e.url, "err", err)
		return
	}

	hasRequestBinding := false
	for _, binding := range bindings {
		if strings.EqualFold(strings.TrimSpace(binding.BindPoint), "REQUEST") {
			hasRequestBinding = true
			break
		}
	}
	if !hasRequestBinding {
		e.lbVServerHTTPHostRewriteInfo.Collect(ch)
		return
	}

	policies, err := netscaler.GetAllRewritePolicies(ctx, nsClient)
	if err != nil {
		e.logger.Error("failed to get rewrite policies", "url", e.url, "err", err)
		return
	}

	actions, err := netscaler.GetAllRewriteActions(ctx, nsClient)
	if err != nil {
		e.logger.Error("failed to get rewrite actions", "url", e.url, "err", err)
		return
	}

	for _, mapping := range resolveHTTPHostRewrites(bindings, policies, actions) {
		labels := e.buildLabelValues(
			mapping.VirtualServer,
			mapping.Policy,
			mapping.Priority,
			mapping.Rule,
			mapping.Action,
			mapping.Host,
		)
		e.lbVServerHTTPHostRewriteInfo.WithLabelValues(labels...).Set(1)
	}

	e.lbVServerHTTPHostRewriteInfo.Collect(ch)
}

func resolveHTTPHostRewrites(
	bindings []netscaler.LBVServerRewritePolicyBinding,
	policies []netscaler.RewritePolicy,
	actions []netscaler.RewriteAction,
) []httpHostRewriteMapping {
	policyByName := make(map[string]netscaler.RewritePolicy, len(policies))
	for _, policy := range policies {
		policyByName[policy.Name] = policy
	}

	actionByName := make(map[string]netscaler.RewriteAction, len(actions))
	for _, action := range actions {
		actionByName[action.Name] = action
	}

	var mappings []httpHostRewriteMapping
	seen := make(map[string]struct{})
	for _, binding := range bindings {
		if !strings.EqualFold(strings.TrimSpace(binding.BindPoint), "REQUEST") {
			continue
		}

		policy, ok := policyByName[binding.PolicyName]
		if !ok || policy.Action == "" {
			continue
		}

		action, ok := actionByName[policy.Action]
		if !ok || !strings.EqualFold(strings.TrimSpace(action.Type), "replace") || !isHostHeaderTarget(action.Target) {
			continue
		}

		host, ok := normalizeStaticHostExpression(action.StringBuilderExpr)
		if !ok {
			continue
		}

		mapping := httpHostRewriteMapping{
			VirtualServer: binding.Name,
			Policy:        binding.PolicyName,
			Priority:      binding.Priority,
			Rule:          policy.Rule,
			Action:        policy.Action,
			Host:          host,
		}
		key := strings.Join([]string{
			mapping.VirtualServer,
			mapping.Policy,
			mapping.Priority,
			mapping.Rule,
			mapping.Action,
			mapping.Host,
		}, "\x00")
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		mappings = append(mappings, mapping)
	}

	sort.Slice(mappings, func(i, j int) bool {
		left := mappings[i]
		right := mappings[j]
		if left.VirtualServer != right.VirtualServer {
			return left.VirtualServer < right.VirtualServer
		}
		if left.Priority != right.Priority {
			return left.Priority < right.Priority
		}
		if left.Policy != right.Policy {
			return left.Policy < right.Policy
		}
		return left.Host < right.Host
	})

	return mappings
}

func isHostHeaderTarget(target string) bool {
	compact := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, target)
	return strings.EqualFold(compact, hostHeaderRewriteTarget)
}

func normalizeStaticHostExpression(expression string) (string, bool) {
	host, err := strconv.Unquote(strings.TrimSpace(expression))
	if err != nil {
		return "", false
	}

	host = strings.TrimSpace(host)
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		host = parsedHost
	}
	host = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")

	if host == "" || strings.ContainsAny(host, " \t\r\n/\\") || strings.Contains(host, ":") {
		return "", false
	}
	return host, true
}
