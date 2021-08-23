/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"path"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/module"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/util/path_util"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/vmware-tanzu/octant/pkg/navigation"
)

type LinkDefinition struct {
	Title string
	Url   string
}

// Generate breadcrumb for specified path
func GenerateBreadcrumb(cm *ContentManager, contentPath string, state octant.State, m module.Module, options module.ContentOptions) []component.TitleComponent {
	var title []component.TitleComponent
	crPath := "custom-resources"

	navs, err := cm.moduleManager.Navigation(cm.ctx, state.GetNamespace(), m.Name())
	if err != nil {
		return title
	}

	parent, title := CreateNavigationBreadcrumb(navs, contentPath)
	if title == nil {
		return title
	} else if parent.Title == "" {
		title = append(title, component.NewText(path.Base(contentPath)))
		return title
	}

	if strings.Contains(contentPath, crPath) {
		if path.Base(contentPath) == parent.Title || path.Base(parent.Url) == crPath {
			title = append(title, component.NewText(parent.Title))
		} else {
			title = append(title, component.NewLink("", parent.Title, path_util.PrefixedPath(parent.Url)), component.NewText(path.Base(contentPath)))
		}
	} else {
		gvk, err := cm.moduleManager.GvkFromPath(path.Dir(contentPath), state.GetNamespace())
		if err != nil {
			title = append(title, component.NewText(parent.Title))
			return title
		}
		key := store.KeyFromGroupVersionKind(gvk)
		key.Selector = options.LabelSet

		if !isClusterScoped(m) {
			key.Namespace = state.GetNamespace()
		}
		list, _, err := cm.dashConfig.ObjectStore().List(cm.ctx, key)
		if err == nil && list.Items != nil && len(list.Items) > 0 {
			second := dropdownFromList(parent, path.Base(contentPath), list)
			title = append(title, second, component.NewText(path.Base(contentPath)))
		} else {
			title = append(title, component.NewText(parent.Title))
		}
	}
	return title
}

// Create first part of breadcrumb from module navigation entries.
// Performs reverse path traversal and creates all related breadcrumbs.
func CreateNavigationBreadcrumb(navs []navigation.Navigation, contentPath string) (LinkDefinition, []component.TitleComponent) {
	var last LinkDefinition
	var title []component.TitleComponent

	// When there is a single non-nested navigation entry, then show it as a link.
	if len(navs) == 1 && contentPath != navs[0].Path && strings.HasPrefix(contentPath, navs[0].Path) {
		title = append(title, component.NewLink("", path.Base(navs[0].Title), path_util.PrefixedPath(navs[0].Path)))
		return LinkDefinition{}, title
	}

	thisPath := contentPath
	for {
		if thisPath == "." { // done
			break
		}
		navItems, parent, selection := NavigationFromPath(navs, thisPath)
		if len(navItems) > 0 {
			dropdown := dropdownFromNavigation(parent, selection.Title, navItems)
			title = append(title, dropdown)
			if len(last.Title) == 0 {
				last = LinkDefinition{Title: selection.Title, Url: selection.Url}
			}
		}
		thisPath = path.Dir(thisPath)
	}
	return last, reverseTitle(title)
}

// Returns all navigation elements for specified path
func NavigationFromPath(navs []navigation.Navigation, navPath string) ([]navigation.Navigation, LinkDefinition, LinkDefinition) {
	for _, nav := range navs {
		for _, child := range nav.Children {
			if child.Path == navPath {
				return nav.Children, LinkDefinition{Title: nav.Title, Url: nav.Path}, LinkDefinition{Title: child.Title, Url: child.Path}
			}
		}
		if nav.Path == navPath && nav.Title != navs[0].Title {
			return navs, LinkDefinition{Title: navs[0].Title, Url: navs[0].Path}, LinkDefinition{Title: nav.Title, Url: nav.Path}
		}
	}
	return []navigation.Navigation{}, LinkDefinition{}, LinkDefinition{}
}

func dropdownFromNavigation(title LinkDefinition, selection string, items []navigation.Navigation) *component.Dropdown {
	dropItems := make([]component.DropdownItemConfig, 0)
	for _, item := range items {
		item := component.NewDropdownItem(item.Title, component.Url, item.Title, item.Path, "")
		dropItems = append(dropItems, item)
	}

	return createLinkDropdown(title, selection, dropItems, false)
}

func dropdownFromList(title LinkDefinition, selection string, items *unstructured.UnstructuredList) *component.Dropdown {
	dropItems := make([]component.DropdownItemConfig, 0)
	for _, item := range items.Items {
		item := component.NewDropdownItem(item.GetName(), component.Url, item.GetName(),
			path.Join(title.Url, item.GetName()), "")
		dropItems = append(dropItems, item)
	}

	return createLinkDropdown(title, selection, dropItems, true)
}

func createLinkDropdown(title LinkDefinition, selection string, items []component.DropdownItemConfig, sortItems bool) *component.Dropdown {
	if sortItems {
		sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	}

	dropdown := component.NewDropdown(title.Title, component.DropdownLink, "", items...)
	dropdown.SetTitle(append([]component.TitleComponent{}, component.NewLink("", title.Title, path_util.PrefixedPath(title.Url))))
	dropdown.SetSelection(selection)
	return dropdown
}

func reverseTitle(title []component.TitleComponent) []component.TitleComponent {

	for i, j := 0, len(title)-1; i < j; i, j = i+1, j-1 {
		title[i], title[j] = title[j], title[i]
	}
	return title
}

func isClusterScoped(m module.Module) bool {
	return m.Name() == "cluster-overview"
}
