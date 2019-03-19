import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import json from 'highlight.js/lib/languages/json';
import yaml from 'highlight.js/lib/languages/yaml';
import { HighlightModule } from 'ngx-highlightjs';
import { PageNotFoundComponent } from 'src/app/util/page-not-found/page-not-found.component';

import { ContentComponent } from '../content/content.component';
import { FiltersComponent } from '../filters/filters.component';
import { TabsComponent } from '../tabs/tabs.component';
import { AnnotationsComponent } from './annotations/annotations.component';
import { ContainersComponent } from './containers/containers.component';
import { DatagridComponent } from './datagrid/datagrid.component';
import { ExpressionSelectorComponent } from './expression-selector/expression-selector.component';
import { FlexlayoutComponent } from './flexlayout/flexlayout.component';
import { LabelSelectorComponent } from './label-selector/label-selector.component';
import { LabelsComponent } from './labels/labels.component';
import { LinkComponent } from './link/link.component';
import { ViewListComponent } from './list/list.component';
import { NamespaceComponent } from './namespace/namespace.component';
import { PortForwardComponent } from './port-forward/port-forward.component';
import { QuadrantComponent } from './quadrant/quadrant.component';
import { ResourceViewerComponent } from './resource-viewer/resource-viewer.component';
import { SelectorsComponent } from './selectors/selectors.component';
import { SummaryComponent } from './summary/summary.component';
import { TableComponent } from './table/table.component';
import { TextComponent } from './text/text.component';
import { TimestampComponent } from './timestamp/timestamp.component';
import { ViewComponent } from './view.component';
import { YamlComponent } from './yaml/yaml.component';

const hljsLanguages = () => {
  return [{ name: 'yaml', func: yaml }, { name: 'json', func: json }];
};

@NgModule({
  declarations: [
    AnnotationsComponent,
    ContentComponent,
    ContainersComponent,
    DatagridComponent,
    ExpressionSelectorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    LabelSelectorComponent,
    LabelsComponent,
    LinkComponent,
    NamespaceComponent,
    ViewListComponent,
    PageNotFoundComponent,
    PortForwardComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SelectorsComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TextComponent,
    TimestampComponent,
    YamlComponent,
    ViewComponent,
    PortForwardComponent,
  ],
  imports: [
    CommonModule,
    ClarityModule,
    HighlightModule.forRoot({
      languages: hljsLanguages,
    }),
    RouterModule,
  ],
  exports: [NamespaceComponent],
})
export class ViewModule {}
