package printer

import (
	"context"
	"sort"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	rbacv1 "k8s.io/api/rbac/v1"
)

func RoleListHandler(ctx context.Context, roleList *rbacv1.RoleList, opts Options) (component.ViewComponent, error) {
	if roleList == nil {
		return nil, errors.New("role list is nil")
	}

	columns := component.NewTableCols("Name", "Age")
	table := component.NewTable("Roles", columns)

	for _, role := range roleList.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&role, role.Name)
		row["Age"] = component.NewTimestamp(role.CreationTimestamp.Time)
		table.Add(row)
	}

	return table, nil
}

func RoleHandler(ctx context.Context, role *rbacv1.Role, opts Options) (component.ViewComponent, error) {
	o := NewObject(role)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printRoleConfig(role)
	}, 16)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return printRolePolicyRules(role)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, opts)
}

func printRoleConfig(role *rbacv1.Role) (component.ViewComponent, error) {
	if role == nil {
		return nil, errors.New("role is nil")
	}

	var sections component.SummarySections
	sections.AddText("Name", role.Name)
	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func printRolePolicyRules(role *rbacv1.Role) (*component.Table, error) {
	if role == nil {
		return nil, errors.New("role is nil")
	}

	rules := role.Rules
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
