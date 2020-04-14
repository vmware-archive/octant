import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DocumentationComponent } from './documentation.component';

import { ComponentOverviewComponent } from './component-overview/component-overview.component';

const documentationRoutes: Routes = [
  {
    path: 'documentation',
    component: DocumentationComponent,
    children: [
      {
        path: '',
        component: ComponentOverviewComponent,
      },
      {
        path: 'annotations',
        loadChildren: () =>
          import('./demos/annotations/annotations.demo.module').then(
            m => m.AnnotationsDemoModule
          ),
      },
      {
        path: 'button-group',
        loadChildren: () =>
          import('./demos/button-group/button-group.demo.module').then(
            m => m.ButtonGroupDemoModule
          ),
      },
      {
        path: 'card',
        loadChildren: () =>
          import('./demos/card/card.demo.module').then(m => m.CardDemoModule),
      },
      // {
      //   path: 'card-list',
      //   loadChildren: () =>
      //     import('./demos/card-list/card-list.demo.module').then(
      //       m => m.CardListDemoModule
      //     ),
      // },
      {
        path: 'code',
        loadChildren: () =>
          import('./demos/code/code.demo.module').then(m => m.CodeDemoModule),
      },
      // {
      //   path: 'containers',
      //   loadChildren: () =>
      //     import('./demos/containers/containers.demo.module').then(
      //       m => m.ContainersDemoModule
      //     ),
      // },
      {
        path: 'donut-chart',
        loadChildren: () =>
          import('./demos/donut-chart/donut-chart.demo.module').then(
            m => m.DonutChartDemoModule
          ),
      },
      {
        path: 'editor',
        loadChildren: () =>
          import('./demos/editor/editor.demo.module').then(
            m => m.EditorDemoModule
          ),
      },
      // {
      //   path: 'expression-selector',
      //   loadChildren: () =>
      //     import(
      //       './demos/expression-selector/expression-selector.demo.module'
      //     ).then(m => m.ExpressionSelectorDemoModule),
      // },
      {
        path: 'flexlayout',
        loadChildren: () =>
          import('./demos/flexlayout/flexlayout.demo.module').then(
            m => m.FlexLayoutDemoModule
          ),
      },
      // {
      //   path: 'graphviz',
      //   loadChildren: () =>
      //     import('./demos/graphviz/graphviz.demo.module').then(
      //       m => m.GraphvizDemoModule
      //     ),
      // },
      {
        path: 'iframe',
        loadChildren: () =>
          import('./demos/iframe/iframe.demo.module').then(
            m => m.IFrameDemoModule
          ),
      },
      // {
      //   path: 'label-selector',
      //   loadChildren: () =>
      //     import('./demos/label-selector/label-selector.demo.module').then(
      //       m => m.LabelSelectorDemoModule
      //     ),
      // },
      {
        path: 'labels',
        loadChildren: () =>
          import('./demos/labels/labels.demo.module').then(
            m => m.LabelsDemoModule
          ),
      },
      {
        path: 'link',
        loadChildren: () =>
          import('./demos/link/link.demo.module').then(m => m.LinkDemoModule),
      },
      {
        path: 'list',
        loadChildren: () =>
          import('./demos/list/list.demo.module').then(m => m.ListDemoModule),
      },
      {
        path: 'logs',
        loadChildren: () =>
          import('./demos/logs/logs.demo.module').then(m => m.LogsDemoModule),
      },
      // {
      //   path: 'pod-status',
      //   loadChildren: () =>
      //     import('./demos/pod-status/pod-status.demo.module').then(
      //       m => m.PodStatusDemoModule
      //     ),
      // },
      // {
      //   path: 'portforward',
      //   loadChildren: () =>
      //     import('./demos/portforward/portforward.demo.module').then(
      //       m => m.PortForwardDemoModule
      //     ),
      // },
      {
        path: 'ports',
        loadChildren: () =>
          import('./demos/ports/ports.demo.module').then(
            m => m.PortsDemoModule
          ),
      },
      {
        path: 'quadrant',
        loadChildren: () =>
          import('./demos/quadrant/quadrant.demo.module').then(
            m => m.QuadrantDemoModule
          ),
      },
      // {
      //   path: 'selectors',
      //   loadChildren: () =>
      //     import('./demos/selectors/selectors.demo.module').then(
      //       m => m.SelectorsDemoModule
      //     ),
      // },
      {
        path: 'summary',
        loadChildren: () =>
          import('./demos/summary/summary.demo.module').then(
            m => m.SummaryDemoModule
          ),
      },
      {
        path: 'table',
        loadChildren: () =>
          import('./demos/table/table.demo.module').then(
            m => m.TableDemoModule
          ),
      },
      {
        path: 'terminal',
        loadChildren: () =>
          import('./demos/terminal/terminal.demo.module').then(
            m => m.TerminalDemoModule
          ),
      },
      {
        path: 'text',
        loadChildren: () =>
          import('./demos/text/text.demo.module').then(m => m.TextDemoModule),
      },
      {
        path: 'timestamp',
        loadChildren: () =>
          import('./demos/timestamp/timestamp.demo.module').then(
            m => m.TimestampDemoModule
          ),
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(documentationRoutes)],
  exports: [RouterModule],
})
export class DocumentationRoutingModule {}
