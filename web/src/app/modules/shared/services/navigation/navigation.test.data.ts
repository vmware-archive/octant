/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
import { NavigationChild } from '../../../sugarloaf/models/navigation';

export const NAVIGATION_MOCK_DATA: NavigationChild[] = [
  {
    title: 'Applications',
    path: 'workloads/namespace/default',
    iconName: 'application',
    isLoading: false,
  },
  {
    title: 'Namespace Overview',
    path: 'overview/namespace/default',
    iconName: 'dashboard',
    isLoading: false,
  },
  {
    title: 'Workloads',
    path: 'overview/namespace/default/workloads',
    children: [
      {
        title: 'Overview',
        path: 'overview/namespace/default/workloads',
        isLoading: false,
      },
      {
        title: 'Cron Jobs',
        path: 'overview/namespace/default/workloads/cron-jobs',
        isLoading: false,
      },
      {
        title: 'Daemon Sets',
        path: 'overview/namespace/default/workloads/daemon-sets',
        isLoading: false,
      },
      {
        title: 'Deployments',
        path: 'overview/namespace/default/workloads/deployments',
        isLoading: false,
      },
      {
        title: 'Jobs',
        path: 'overview/namespace/default/workloads/jobs',
        isLoading: false,
      },
      {
        title: 'Pods',
        path: 'overview/namespace/default/workloads/pods',
        isLoading: false,
      },
      {
        title: 'Replica Sets',
        path: 'overview/namespace/default/workloads/replica-sets',
        isLoading: false,
      },
      {
        title: 'Replication Controllers',
        path: 'overview/namespace/default/workloads/replication-controllers',
        isLoading: false,
      },
      {
        title: 'Stateful Sets',
        path: 'overview/namespace/default/workloads/stateful-sets',
        isLoading: false,
      },
    ],
    iconName: 'applications',
    isLoading: false,
  },
  {
    title: 'Discovery and Load Balancing',
    path: 'overview/namespace/default/discovery-and-load-balancing',
    children: [
      {
        title: 'Overview',
        path: 'overview/namespace/default/discovery-and-load-balancing',
        isLoading: false,
      },
      {
        title: 'Horizontal Pod Autoscalers',
        path:
          'overview/namespace/default/discovery-and-load-balancing/horizontal-pod-autoscalers',
        isLoading: false,
      },
      {
        title: 'Ingresses',
        path:
          'overview/namespace/default/discovery-and-load-balancing/ingresses',
        isLoading: false,
      },
      {
        title: 'Services',
        path:
          'overview/namespace/default/discovery-and-load-balancing/services',
        isLoading: false,
      },
    ],
    iconName: 'network-globe',
    isLoading: false,
  },
  {
    title: 'Config and Storage',
    path: 'overview/namespace/default/config-and-storage',
    children: [
      {
        title: 'Overview',
        path: 'overview/namespace/default/config-and-storage',
        isLoading: false,
      },
      {
        title: 'Config Maps',
        path: 'overview/namespace/default/config-and-storage/config-maps',
        isLoading: false,
      },
      {
        title: 'Persistent Volume Claims',
        path:
          'overview/namespace/default/config-and-storage/persistent-volume-claims',
        isLoading: false,
      },
      {
        title: 'Secrets',
        path: 'overview/namespace/default/config-and-storage/secrets',
        isLoading: false,
      },
      {
        title: 'Service Accounts',
        path: 'overview/namespace/default/config-and-storage/service-accounts',
        isLoading: false,
      },
    ],
    iconName: 'storage',
    isLoading: false,
  },
  {
    title: 'Custom Resources',
    path: 'overview/namespace/default/custom-resources',
    children: [
      {
        title: 'Overview',
        path: 'overview/namespace/default/custom-resources',
        isLoading: false,
      },
    ],
    iconName: 'file-group',
    isLoading: false,
  },
  {
    title: 'RBAC',
    path: 'overview/namespace/default/rbac',
    children: [
      {
        title: 'Overview',
        path: 'overview/namespace/default/rbac',
        isLoading: false,
      },
      {
        title: 'Roles',
        path: 'overview/namespace/default/rbac/roles',
        isLoading: false,
      },
      {
        title: 'Role Bindings',
        path: 'overview/namespace/default/rbac/role-bindings',
        isLoading: false,
      },
    ],
    iconName: 'assign-user',
    isLoading: false,
  },
  {
    title: 'Events',
    path: 'overview/namespace/default/events',
    iconName: 'event',
    isLoading: false,
  },
  {
    module: 'Cluster',
    title: 'Cluster Overview',
    path: 'cluster-overview',
    iconName: 'dashboard',
    isLoading: false,
  },
  {
    title: 'Namespaces',
    path: 'cluster-overview/namespaces',
    iconName: 'namespace',
    isLoading: false,
  },
  {
    title: 'Custom Resources',
    path: 'cluster-overview/custom-resources',
    children: [
      {
        title: 'Overview',
        path: 'cluster-overview/custom-resources',
        isLoading: false,
      },
      {
        title:
          'antreacontrollerinfos.clusterinformation.antrea.tanzu.vmware.com',
        path:
          'cluster-overview/custom-resources/antreacontrollerinfos.clusterinformation.antrea.tanzu.vmware.com',
        iconName: 'internal:crd',
        isLoading: false,
      },
    ],
    iconName: 'file-group',
    isLoading: false,
  },
  {
    title: 'RBAC',
    path: 'cluster-overview/rbac',
    children: [
      {
        title: 'Overview',
        path: 'cluster-overview/rbac',
        isLoading: false,
      },
      {
        title: 'Cluster Roles',
        path: 'cluster-overview/rbac/cluster-roles',
        isLoading: false,
      },
      {
        title: 'Cluster Role Bindings',
        path: 'cluster-overview/rbac/cluster-role-bindings',
        isLoading: false,
      },
    ],
    iconName: 'assign-user',
    isLoading: false,
  },
  {
    title: 'Nodes',
    path: 'cluster-overview/nodes',
    iconName: 'nodes',
    isLoading: false,
  },
  {
    title: 'Storage',
    path: 'cluster-overview/storage',
    children: [
      {
        title: 'Overview',
        path: 'cluster-overview/storage',
        isLoading: false,
      },
      {
        title: 'Persistent Volumes',
        path: 'cluster-overview/storage/persistent-volumes',
        isLoading: false,
      },
    ],
    iconName: 'storage',
    isLoading: false,
  },
  {
    title: 'Port Forwards',
    path: 'cluster-overview/port-forward',
    iconName: 'router',
    isLoading: false,
  },
  {
    module: 'Configuration',
    title: 'Plugin',
    path: 'configuration/plugins',
    iconName: 'plugin',
    isLoading: false,
  },
  {
    title: 'Argo UI',
    path: 'argo-ui',
    iconName: 'cloud',
    isLoading: false,
  },
  {
    title: 'Sample Plugin',
    path: 'plugin-name',
    children: [
      {
        title: 'Nested Once',
        path: 'plugin-name/nested-once',
        children: [
          {
            title: 'Nested Twice',
            path: 'plugin-name/nested-once/nested-twice',
            iconName: 'folder',
            isLoading: false,
          },
        ],
        iconName: 'folder',
        isLoading: false,
      },
    ],
    iconName: 'cloud',
    isLoading: false,
  },
  {
    title: 'OpenStack',
    path: 'openstack',
    iconName: 'cloud',
    isLoading: false,
  },
];
