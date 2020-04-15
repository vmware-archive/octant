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
        iconSource:
          '<?xml version="1.0" encoding="UTF-8" standalone="no"?>\n<!-- Created with Inkscape (http://www.inkscape.org/) -->\n\n<svg\n   xmlns:dc="http://purl.org/dc/elements/1.1/"\n   xmlns:cc="http://creativecommons.org/ns#"\n   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"\n   xmlns:svg="http://www.w3.org/2000/svg"\n   xmlns="http://www.w3.org/2000/svg"\n   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"\n   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"\n   width="18.035334mm"\n   height="17.500378mm"\n   viewBox="0 0 18.035334 17.500378"\n   version="1.1"\n   id="svg13826"\n   inkscape:version="0.91 r13725"\n   sodipodi:docname="crd.svg">\n  <defs\n     id="defs13820" />\n  <sodipodi:namedview\n     id="base"\n     pagecolor="#ffffff"\n     bordercolor="#666666"\n     borderopacity="1.0"\n     inkscape:pageopacity="0.0"\n     inkscape:pageshadow="2"\n     inkscape:zoom="8"\n     inkscape:cx="3.972496"\n     inkscape:cy="33.752239"\n     inkscape:document-units="mm"\n     inkscape:current-layer="layer1"\n     showgrid="false"\n     inkscape:window-width="1440"\n     inkscape:window-height="775"\n     inkscape:window-x="0"\n     inkscape:window-y="1"\n     inkscape:window-maximized="1"\n     fit-margin-top="0"\n     fit-margin-left="0"\n     fit-margin-right="0"\n     fit-margin-bottom="0" />\n  <metadata\n     id="metadata13823">\n    <rdf:RDF>\n      <cc:Work\n         rdf:about="">\n        <dc:format>image/svg+xml</dc:format>\n        <dc:type\n           rdf:resource="http://purl.org/dc/dcmitype/StillImage" />\n        <dc:title />\n      </cc:Work>\n    </rdf:RDF>\n  </metadata>\n  <g\n     inkscape:label="Calque 1"\n     inkscape:groupmode="layer"\n     id="layer1"\n     transform="translate(-0.99262638,-1.174181)">\n    <g\n       id="g70"\n       transform="matrix(1.0148887,0,0,1.0148887,16.902146,-2.698726)">\n      <path\n         inkscape:export-ydpi="250.55"\n         inkscape:export-xdpi="250.55"\n         inkscape:export-filename="new.png"\n         inkscape:connector-curvature="0"\n         id="path3055"\n         d="m -6.8492015,4.2724668 a 1.1191255,1.1099671 0 0 0 -0.4288818,0.1085303 l -5.8524037,2.7963394 a 1.1191255,1.1099671 0 0 0 -0.605524,0.7529759 l -1.443828,6.2812846 a 1.1191255,1.1099671 0 0 0 0.151943,0.851028 1.1191255,1.1099671 0 0 0 0.06362,0.08832 l 4.0508,5.036555 a 1.1191255,1.1099671 0 0 0 0.874979,0.417654 l 6.4961011,-0.0015 a 1.1191255,1.1099671 0 0 0 0.8749788,-0.416906 L 1.3818872,15.149453 A 1.1191255,1.1099671 0 0 0 1.5981986,14.210104 L 0.15212657,7.9288154 A 1.1191255,1.1099671 0 0 0 -0.45339794,7.1758396 L -6.3065496,4.3809971 A 1.1191255,1.1099671 0 0 0 -6.8492015,4.2724668 Z"\n         style="fill:#326ce5;fill-opacity:1;stroke:none;stroke-width:0;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" />\n\n    </g>\n    <path\n       style="fill:#ffffff;fill-opacity:1;stroke-width:0.46185368"\n       inkscape:connector-curvature="0"\n       d="m 14.269566,9.6934441 -0.692781,0 0,-1.847416 c 0,-0.5080383 -0.415668,-0.9237063 -0.923706,-0.9237063 l -1.847416,0 0,-0.69278 c 0,-0.637359 -0.517275,-1.154636 -1.1546334,-1.154636 -0.6373578,0 -1.1546338,0.517277 -1.1546338,1.154636 l 0,0.69278 -1.847416,0 c -0.508038,0 -0.919089,0.415668 -0.919089,0.9237063 l 0,1.755045 0.688164,0 c 0.688161,0 1.247005,0.5588429 1.247005,1.2470049 0,0.688162 -0.558844,1.247005 -1.247005,1.247005 l -0.692781,0 0,1.755045 c 0,0.508038 0.415668,0.923706 0.923706,0.923706 l 1.755045,0 0,-0.692781 c 0,-0.688161 0.558843,-1.247005 1.2470048,-1.247005 0.6881624,0 1.2470044,0.558844 1.2470044,1.247005 l 0,0.692781 1.755045,0 c 0.508038,0 0.923706,-0.415668 0.923706,-0.923706 l 0,-1.847416 0.692781,0 c 0.637358,0 1.154635,-0.517276 1.154635,-1.154634 0,-0.637358 -0.517277,-1.1546339 -1.154635,-1.1546339 z"\n       id="path7415" />\n  </g>\n</svg>\n',
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
