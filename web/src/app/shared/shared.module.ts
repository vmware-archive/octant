import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TextComponent } from './shared/components/presentation/text/text.component';
import { MarkdownModule } from 'ngx-markdown';
import { ClarityModule } from '@clr/angular';
import { TitleComponent } from './shared/components/presentation/title/title.component';
import { AlertComponent } from './shared/components/presentation/alert/alert.component';
import { AnnotationsComponent } from './shared/components/presentation/annotations/annotations.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CardComponent } from './shared/components/presentation/card/card.component';
import { CardListComponent } from './shared/components/presentation/card-list/card-list.component';
import { CodeComponent } from './shared/components/presentation/code/code.component';
import { LabelsComponent } from './shared/components/presentation/labels/labels.component';
import { LinkComponent } from './shared/components/presentation/link/link.component';
import { ListComponent } from './shared/components/presentation/list/list.component';
import { TabsComponent } from './shared/components/presentation/tabs/tabs.component';
import { ContainersComponent } from './shared/components/presentation/containers/containers.component';
import { DatagridComponent } from './shared/components/presentation/datagrid/datagrid.component';
import { DonutChartComponent } from './shared/components/presentation/donut-chart/donut-chart.component';
import { FlexlayoutComponent } from './shared/components/presentation/flexlayout/flexlayout.component';
import { SingleStatComponent } from './shared/components/presentation/single-stat/single-stat.component';
import { QuadrantComponent } from './shared/components/presentation/quadrant/quadrant.component';
import { IFrameComponent } from './shared/components/presentation/iframe/iframe.component';
import { ErrorComponent } from './shared/components/presentation/error/error.component';
import { ExpressionSelectorComponent } from './shared/components/presentation/expression-selector/expression-selector.component';
import { GraphvizComponent } from './shared/components/presentation/graphviz/graphviz.component';
import { ButtonGroupComponent } from './shared/components/presentation/button-group/button-group.component';
import { YamlComponent } from './shared/components/presentation/yaml/yaml.component';
import { TableComponent } from './shared/components/presentation/table/table.component';
import { TimestampComponent } from './shared/components/presentation/timestamp/timestamp.component';
import { LoadingComponent } from './shared/components/presentation/loading/loading.component';
import { HighlightModule } from 'ngx-highlightjs';
import { LabelSelectorComponent } from './shared/components/presentation/label-selector/label-selector.component';
import { CytoscapeComponent } from './shared/components/presentation/cytoscape/cytoscape.component';
import { SelectorsComponent } from './shared/components/presentation/selectors/selectors.component';
import { ResourceViewerComponent } from './shared/components/presentation/resource-viewer/resource-viewer.component';
import { SummaryComponent } from './shared/components/presentation/summary/summary.component';
import { PortForwardComponent } from './shared/components/presentation/port-forward/port-forward.component';
import { FormComponent } from './shared/components/presentation/form/form.component';
import { ContentFilterComponent } from './shared/components/presentation/content-filter/content-filter.component';
import { HeptagonGridComponent } from './shared/components/presentation/heptagon-grid/heptagon-grid.component';
import { HeptagonGridRowComponent } from './shared/components/presentation/heptagon-grid-row/heptagon-grid-row.component';
import { HeptagonLabelComponent } from './shared/components/presentation/heptagon-label/heptagon-label.component';
import { ContentSwitcherComponent } from './shared/components/presentation/content-switcher/content-switcher.component';
import { ObjectStatusComponent } from './shared/components/presentation/object-status/object-status.component';
import { PodStatusComponent } from './shared/components/presentation/pod-status/pod-status.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TerminalComponent } from './shared/components/smart/terminal/terminal.component';
import { LogsComponent } from './shared/components/smart/logs/logs.component';
import { PortsComponent } from './shared/components/presentation/ports/ports.component';
import { FiltersComponent } from './shared/components/smart/filters/filters.component';
import { HeptagonComponent } from './shared/components/smart/heptagon/heptagon.component';
import { ContextSelectorComponent } from './shared/components/smart/context-selector/context-selector.component';
import { SliderViewComponent } from './shared/components/smart/slider-view/slider-view.component';
import { SafePipe } from './shared/pipes/safe/safe.pipe';
import { DefaultPipe } from './shared/pipes/default/default.pipe';
import { RouterModule } from '@angular/router';
import { ResizableModule } from 'angular-resizable-element';
import { hljsLanguages } from './highlight';

@NgModule({
  declarations: [
    AlertComponent,
    AnnotationsComponent,
    ButtonGroupComponent,
    CardComponent,
    CardListComponent,
    CodeComponent,
    ContainersComponent,
    ContentFilterComponent,
    ContentSwitcherComponent,
    ContextSelectorComponent,
    CytoscapeComponent,
    DatagridComponent,
    DefaultPipe,
    DonutChartComponent,
    ErrorComponent,
    ExpressionSelectorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    FormComponent,
    GraphvizComponent,
    HeptagonComponent,
    HeptagonGridComponent,
    HeptagonGridRowComponent,
    HeptagonLabelComponent,
    IFrameComponent,
    LabelsComponent,
    LabelSelectorComponent,
    LinkComponent,
    ListComponent,
    LoadingComponent,
    LogsComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    PortForwardComponent,
    PortsComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SafePipe,
    SelectorsComponent,
    SingleStatComponent,
    SliderViewComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TerminalComponent,
    TextComponent,
    TimestampComponent,
    TitleComponent,
    YamlComponent,
  ],
  imports: [
    ClarityModule,
    CommonModule,
    FormsModule,
    HighlightModule.forRoot({
      languages: hljsLanguages,
    }),
    MarkdownModule.forChild(),
    ReactiveFormsModule,
    ResizableModule,
    RouterModule,
  ],
  exports: [
    AlertComponent,
    AnnotationsComponent,
    ButtonGroupComponent,
    CardComponent,
    CardListComponent,
    CodeComponent,
    ContainersComponent,
    ContentFilterComponent,
    ContentSwitcherComponent,
    ContextSelectorComponent,
    CytoscapeComponent,
    DatagridComponent,
    DefaultPipe,
    DonutChartComponent,
    ErrorComponent,
    ExpressionSelectorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    FormComponent,
    GraphvizComponent,
    HeptagonComponent,
    HeptagonGridComponent,
    HeptagonGridRowComponent,
    HeptagonLabelComponent,
    IFrameComponent,
    LabelsComponent,
    LabelSelectorComponent,
    LinkComponent,
    ListComponent,
    LoadingComponent,
    LogsComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    PortForwardComponent,
    PortsComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SelectorsComponent,
    SliderViewComponent,
    SingleStatComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TerminalComponent,
    TextComponent,
    TimestampComponent,
    TitleComponent,
    YamlComponent,
  ],
})
export class SharedModule {}
