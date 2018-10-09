package overview

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

type pathFilter struct {
	path      string
	describer Describer

	re *regexp.Regexp
}

func newPathFilter(path string, describer Describer) *pathFilter {
	re := regexp.MustCompile(fmt.Sprintf("^%s/?$", path))

	return &pathFilter{
		re:        re,
		path:      path,
		describer: describer,
	}
}

func (pf *pathFilter) Match(path string) bool {
	return pf.re.MatchString(path)
}

func (pf *pathFilter) Fields(path string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(path)
	for i, name := range pf.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}

var (
	workloadsDescriber = NewSectionDescriber(
		NewCronJobDescriber(),
		NewDeploymentsDescriber(),
	)

	discoveryAndLoadBalancingDescriber = NewSectionDescriber()

	configAndStorageDescriber = NewSectionDescriber()

	customResourcesDescriber = NewSectionDescriber()

	rbacDescriber = NewSectionDescriber()

	rootDescriber = NewSectionDescriber(
		workloadsDescriber,
		discoveryAndLoadBalancingDescriber,
		configAndStorageDescriber,
		customResourcesDescriber,
		rbacDescriber,
	)
)

// defaultPathFilters are a default set of path filters. These are
// currently hand crafted, but they could be generated from the API
// types as well.
var defaultPathFilters = []pathFilter{
	*newPathFilter(
		"/",
		rootDescriber,
	),
	*newPathFilter(
		"/workloads",
		workloadsDescriber,
	),
	*newPathFilter(
		"/workloads/cron-jobs",
		NewCronJobsDescriber(),
	),
	*newPathFilter(
		"/workloads/cron-jobs/(?P<name>.*?)",
		NewCronJobDescriber(),
	),
	*newPathFilter(
		"/workloads/deployments",
		NewDeploymentsDescriber(),
	),
	*newPathFilter(
		"/workloads/deployments/(?P<name>.*?)",
		NewDeploymentDescriber(),
	),
	*newPathFilter(
		"/discovery-and-load-balancing",
		discoveryAndLoadBalancingDescriber,
	),
	*newPathFilter(
		"/config-and-storage",
		configAndStorageDescriber,
	),
	*newPathFilter(
		"/custom-resources",
		customResourcesDescriber,
	),
	*newPathFilter(
		"/rbac",
		rbacDescriber,
	),
	*newPathFilter(
		"/events",
		NewEventsDescriber(),
	),
}

var navPaths = []string{
	"/workloads/daemon-sets",
	"/workloads/jobs",
	"/workloads/pods",
	"/workloads/replica-sets",
	"/workloads/replication-controllers",
	"/workloads/stateful-sets",

	"/discovery-and-load-balancing/ingresses",
	"/discovery-and-load-balancing/services",

	"/config-and-storage/config-maps",
	"/config-and-storage/persistent-volume-claims",
	"/config-and-storage/secrets",
	"/config-and-storage",

	"/rbac/roles",
	"/rbac/role-bindings",
}

var contentNotFound = errors.Errorf("content not found")

type generator interface {
	Generate(path, prefix, namespace string) ([]Content, error)
}

type realGenerator struct {
	cache       Cache
	pathFilters []pathFilter

	mu sync.Mutex
}

func newGenerator(cache Cache, pathFilters []pathFilter) *realGenerator {
	return &realGenerator{
		cache:       cache,
		pathFilters: pathFilters,
	}
}

func (g *realGenerator) Generate(path, prefix, namespace string) ([]Content, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if stringInSlice(path, navPaths) {
		return stubContent(path), nil
	}

	for _, pf := range g.pathFilters {
		if !pf.Match(path) {
			continue
		}

		fields := pf.Fields(path)

		return pf.describer.Describe(prefix, namespace, g.cache, fields)
	}

	return nil, contentNotFound
}

func stringInSlice(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

func stubContent(name string) []Content {
	t := newTable(name)
	t.Columns = []tableColumn{
		{Name: "foo", Accessor: "foo"},
		{Name: "bar", Accessor: "bar"},
		{Name: "baz", Accessor: "baz"},
	}

	t.Rows = []tableRow{
		{
			"foo": newStringText("r1c1"),
			"bar": newStringText("r1c2"),
			"baz": newStringText("r1c3"),
		},
		{
			"foo": newStringText("r2c1"),
			"bar": newStringText("r2c2"),
			"baz": newStringText("r2c3"),
		},
		{
			"foo": newStringText("r3c1"),
			"bar": newStringText("r3c2"),
			"baz": newStringText("r3c3"),
		},
	}

	return []Content{t}
}
