import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TextComponent } from './components/presentation/text/text.component';
import { MarkdownModule } from 'ngx-markdown';
import { ClarityModule } from '@clr/angular';
import { TitleComponent } from './components/presentation/title/title.component';
import { AlertComponent } from './components/presentation/alert/alert.component';
import { AnnotationsComponent } from './components/presentation/annotations/annotations.component';
import { CardComponent } from './components/presentation/card/card.component';
import { CardListComponent } from './components/presentation/card-list/card-list.component';
import { CodeComponent } from './components/presentation/code/code.component';
import { LabelsComponent } from './components/presentation/labels/labels.component';
import { LinkComponent } from './components/presentation/link/link.component';
import { ListComponent } from './components/presentation/list/list.component';
import { TabsComponent } from './components/presentation/tabs/tabs.component';
import { ContainersComponent } from './components/presentation/containers/containers.component';
import { DatagridComponent } from './components/presentation/datagrid/datagrid.component';
import { DonutChartComponent } from './components/presentation/donut-chart/donut-chart.component';
import { FlexlayoutComponent } from './components/presentation/flexlayout/flexlayout.component';
import { SingleStatComponent } from './components/presentation/single-stat/single-stat.component';
import { QuadrantComponent } from './components/presentation/quadrant/quadrant.component';
import { IFrameComponent } from './components/presentation/iframe/iframe.component';
import { ErrorComponent } from './components/presentation/error/error.component';
import { ExpressionSelectorComponent } from './components/presentation/expression-selector/expression-selector.component';
import { GraphvizComponent } from './components/presentation/graphviz/graphviz.component';
import { ButtonGroupComponent } from './components/presentation/button-group/button-group.component';
import { YamlComponent } from './components/presentation/yaml/yaml.component';
import { TableComponent } from './components/presentation/table/table.component';
import { TimestampComponent } from './components/presentation/timestamp/timestamp.component';
import { LoadingComponent } from './components/presentation/loading/loading.component';
import { HighlightModule } from 'ngx-highlightjs';
import { LabelSelectorComponent } from './components/presentation/label-selector/label-selector.component';
import { CytoscapeComponent } from './components/presentation/cytoscape/cytoscape.component';
import { SelectorsComponent } from './components/presentation/selectors/selectors.component';
import { ResourceViewerComponent } from './components/presentation/resource-viewer/resource-viewer.component';
import { SummaryComponent } from './components/presentation/summary/summary.component';
import { PortForwardComponent } from './components/presentation/port-forward/port-forward.component';
import { FormComponent } from './components/presentation/form/form.component';
import { ContentFilterComponent } from './components/presentation/content-filter/content-filter.component';
import { HeptagonGridComponent } from './components/presentation/heptagon-grid/heptagon-grid.component';
import { HeptagonGridRowComponent } from './components/presentation/heptagon-grid-row/heptagon-grid-row.component';
import { HeptagonLabelComponent } from './components/presentation/heptagon-label/heptagon-label.component';
import { ContentSwitcherComponent } from './components/presentation/content-switcher/content-switcher.component';
import { ObjectStatusComponent } from './components/presentation/object-status/object-status.component';
import { PodStatusComponent } from './components/presentation/pod-status/pod-status.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TerminalComponent } from './components/smart/terminal/terminal.component';
import { LogsComponent } from './components/smart/logs/logs.component';
import { PortsComponent } from './components/presentation/ports/ports.component';
import { FiltersComponent } from './components/smart/filters/filters.component';
import { HeptagonComponent } from './components/smart/heptagon/heptagon.component';
import { ContextSelectorComponent } from './components/smart/context-selector/context-selector.component';
import { SliderViewComponent } from './components/smart/slider-view/slider-view.component';
import { SafePipe } from './pipes/safe/safe.pipe';
import { AnsiPipe } from './pipes/ansiPipe/ansi.pipe';
import { DefaultPipe } from './pipes/default/default.pipe';
import { RouterModule } from '@angular/router';
import { ResizableModule } from 'angular-resizable-element';
import { hljsLanguages } from './highlight';
import { IndicatorComponent } from './components/presentation/indicator/indicator.component';

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
    IndicatorComponent,
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
    AnsiPipe,
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
