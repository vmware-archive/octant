// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import json from 'highlight.js/lib/languages/json';
import yaml from 'highlight.js/lib/languages/yaml';
import { HighlightModule } from 'ngx-highlightjs';

import { AnnotationsComponent } from './components/annotations/annotations.component';
import { ContainersComponent } from './components/containers/containers.component';
import { ContentSwitcherComponent } from './components/content-switcher/content-switcher.component';
import { DatagridComponent } from './components/datagrid/datagrid.component';
import { ExpressionSelectorComponent } from './components/expression-selector/expression-selector.component';
import { ErrorComponent } from './components/error/error.component';
import { FiltersComponent } from './components/filters/filters.component';
import { FlexlayoutComponent } from './components/flexlayout/flexlayout.component';
import { GraphvizComponent } from './components/graphviz/graphviz.component';
import { LabelSelectorComponent } from './components/label-selector/label-selector.component';
import { LabelsComponent } from './components/labels/labels.component';
import { LinkComponent } from './components/link/link.component';
import { ListComponent } from './components/list/list.component';
import { LoadingComponent } from './components/loading/loading.component';
import { LogsComponent } from './components/logs/logs.component';
import { ObjectStatusComponent } from './components/object-status/object-status.component';
import { PodStatusComponent } from './components/pod-status/pod-status.component';
import { PortForwardComponent } from './components/port-forward/port-forward.component';
import { PortsComponent } from './components/ports/ports.component';
import { QuadrantComponent } from './components/quadrant/quadrant.component';
import { ResourceViewerComponent } from './components/resource-viewer/resource-viewer.component';
import { SelectorsComponent } from './components/selectors/selectors.component';
import { SummaryComponent } from './components/summary/summary.component';
import { TableComponent } from './components/table/table.component';
import { TabsComponent } from './components/tabs/tabs.component';
import { TextComponent } from './components/text/text.component';
import { TimestampComponent } from './components/timestamp/timestamp.component';
import { YamlComponent } from './components/yaml/yaml.component';
import { OverviewComponent } from './overview.component';
import { HeptagonGridComponent } from './components/heptagon-grid/heptagon-grid.component';
import { HeptagonGridRowComponent } from './components/heptagon-grid-row/heptagon-grid-row.component';
import { HeptagonComponent } from './components/heptagon/heptagon.component';
import { HeptagonLabelComponent } from './components/heptagon-label/heptagon-label.component';
import { ContextSelectorComponent } from './components/context-selector/context-selector.component';
import { DefaultPipe } from './pipes/default.pipe';

export function hljsLanguages() {
  return [{ name: 'yaml', func: yaml }, { name: 'json', func: json }];
}

@NgModule({
  declarations: [
    AnnotationsComponent,
    ContainersComponent,
    DatagridComponent,
    ExpressionSelectorComponent,
    ErrorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    GraphvizComponent,
    LabelSelectorComponent,
    LabelsComponent,
    LoadingComponent,
    LinkComponent,
    ListComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SelectorsComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TextComponent,
    TimestampComponent,
    YamlComponent,
    OverviewComponent,
    PortForwardComponent,
    ContentSwitcherComponent,
    LogsComponent,
    PortsComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    HeptagonGridComponent,
    HeptagonGridRowComponent,
    HeptagonComponent,
    HeptagonLabelComponent,
    ContextSelectorComponent,
    DefaultPipe,
  ],
  imports: [
    CommonModule,
    ClarityModule,
    FormsModule,
    HighlightModule.forRoot({
      languages: hljsLanguages,
    }),
    RouterModule,
  ],
  exports: [ContextSelectorComponent, DefaultPipe],
})
export class OverviewModule {}
