package printer

import (
	"sort"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"

	rbacv1 "k8s.io/api/rbac/v1"
)

// ClusterRoleListHandler is a printFunc that prints cluster roles
func ClusterRoleListHandler(list *rbacv1.ClusterRoleList, options Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("cluster role list is nil")
	}

	cols := component.NewTableCols("Name", "Age")
	tbl := component.NewTable("Cluster Roles", cols)

	for _, clusterRole := range list.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&clusterRole, clusterRole.Name)
		ts := clusterRole.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// ClusterRoleHandler is a printFunc that prints a cluster role
func ClusterRoleHandler(clusterRole *rbacv1.ClusterRole, options Options) (component.ViewComponent, error) {
	o := NewObject(clusterRole)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printClusterRoleConfig(clusterRole)
	}, 12)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return printClusterRolePolicyRules(clusterRole)
		},
		Width: 24,
	})

	return o.ToComponent(options)
}

func printClusterRoleConfig(clusterRole *rbacv1.ClusterRole) (component.ViewComponent, error) {
	if clusterRole == nil {
		return nil, errors.New("cluster role is nil")
	}

	var sections component.SummarySections

	if clusterRoleAggregation := clusterRole.AggregationRule; clusterRoleAggregation != nil {
		if clusterRoleSelectors := clusterRoleAggregation.ClusterRoleSelectors; clusterRoleSelectors != nil {
			var selectors []component.Selector

			for _, selector := range clusterRoleSelectors {
				for k, v := range selector.MatchLabels {
					ls := component.NewLabelSelector(k, v)
					selectors = append(selectors, ls)
				}

				sections = append(sections, component.SummarySection{
					Header:  "Selectors",
					Content: component.NewSelectors(selectors),
				})
			}
		}
	}

	sections.AddText("Name", clusterRole.Name)
	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func printClusterRolePolicyRules(clusterRole *rbacv1.ClusterRole) (*component.Table, error) {
	if clusterRole == nil {
		return nil, errors.New("cluster role is nil")
	}

	rules := clusterRole.Rules
	breakdownRules := []rbacv1.PolicyRule{}
	for _, rule := range rules {
		breakdownRules = append(breakdownRules, BreakdownRule(rule)...)
	}

	compactRules, err := CompactRules(breakdownRules)
	if err != nil {
		return nil, errors.New("cannot compact rules")
	}
	sort.Stable(SortableRuleSlice(compactRules))

	cols := component.NewTableCols("Resources", "Non-Resource URLs", "Resource Names", "Verbs")
	tbl := component.NewTable("PolicyRules", cols)

	for _, r := range compactRules {
		row := component.TableRow{}
		row["Resources"] = component.NewText(CombineResourceGroup(r.Resources, r.APIGroups))
		row["Non-Resource URLs"] = component.NewText(printSlice(r.NonResourceURLs))
		row["Resource Names"] = component.NewText(printSlice(r.ResourceNames))
		row["Verbs"] = component.NewText(printSlice(r.Verbs))

		tbl.Add(row)
	}

	return tbl, nil
}
