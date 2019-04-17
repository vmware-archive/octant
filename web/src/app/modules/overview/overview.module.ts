import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
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
import { FiltersComponent } from './components/filters/filters.component';
import { FlexlayoutComponent } from './components/flexlayout/flexlayout.component';
import { GraphvizComponent } from './components/graphviz/graphviz.component';
import { LabelSelectorComponent } from './components/label-selector/label-selector.component';
import { LabelsComponent } from './components/labels/labels.component';
import { LinkComponent } from './components/link/link.component';
import { ListComponent } from './components/list/list.component';
import { LogsComponent } from './components/logs/logs.component';
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
import { ObjectStatusComponent } from './components/object-status/object-status.component';

export function hljsLanguages() {
  return [{ name: 'yaml', func: yaml }, { name: 'json', func: json }];
}

@NgModule({
  declarations: [
    AnnotationsComponent,
    ContainersComponent,
    DatagridComponent,
    ExpressionSelectorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    GraphvizComponent,
    LabelSelectorComponent,
    LabelsComponent,
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
  ],
  imports: [
    CommonModule,
    ClarityModule,
    HighlightModule.forRoot({
      languages: hljsLanguages,
    }),
    RouterModule,
  ],
  exports: [],
})
export class OverviewModule {}
