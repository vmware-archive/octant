import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TextComponent } from './components/presentation/text/text.component';
import { ClarityModule } from '@clr/angular';
import { CdsModule } from '@cds/angular';
import { TitleComponent } from './components/presentation/title/title.component';
import { AccordionComponent } from './components/presentation/accordion/accordion.component';
import { AlertComponent } from './components/presentation/alert/alert.component';
import { AnnotationsComponent } from './components/presentation/annotations/annotations.component';
import { BreadcrumbComponent } from './components/presentation/breadcrumb/breadcrumb.component';
import { CardComponent } from './components/presentation/card/card.component';
import { CardListComponent } from './components/presentation/card-list/card-list.component';
import { CodeComponent } from './components/presentation/code/code.component';
import { DropdownComponent } from './components/presentation/dropdown/dropdown.component';
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
import { EditorComponent } from './components/smart/editor/editor.component';
import { ErrorComponent } from './components/presentation/error/error.component';
import { ExpressionSelectorComponent } from './components/presentation/expression-selector/expression-selector.component';
import { GraphvizComponent } from './components/presentation/graphviz/graphviz.component';
import { ButtonGroupComponent } from './components/presentation/button-group/button-group.component';
import { YamlComponent } from './components/presentation/yaml/yaml.component';
import { TableComponent } from './components/presentation/table/table.component';
import { TimelineComponent } from './components/presentation/timeline/timeline.component';
import { TimestampComponent } from './components/presentation/timestamp/timestamp.component';
import { LoadingComponent } from './components/presentation/loading/loading.component';
import { HighlightModule } from 'ngx-highlightjs';
import { LabelSelectorComponent } from './components/presentation/label-selector/label-selector.component';
import { ModalComponent } from './components/presentation/modal/modal.component';
import { CytoscapeComponent } from './components/presentation/cytoscape/cytoscape.component';
import { Cytoscape2Component } from './components/presentation/cytoscape2/cytoscape2.component';
import { SelectorsComponent } from './components/presentation/selectors/selectors.component';
import { ResourceViewerComponent } from './components/presentation/resource-viewer/resource-viewer.component';
import { SummaryComponent } from './components/presentation/summary/summary.component';
import { PortForwardComponent } from './components/presentation/port-forward/port-forward.component';
import { FormComponent } from './components/presentation/form/form.component';
import { ContentFilterComponent } from './components/presentation/content-filter/content-filter.component';
import { HeptagonGridComponent } from './components/presentation/heptagon-grid/heptagon-grid.component';
import { HeptagonGridRowComponent } from './components/presentation/heptagon-grid-row/heptagon-grid-row.component';
import { HeptagonLabelComponent } from './components/presentation/heptagon-label/heptagon-label.component';
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
import { highlightProvider } from './highlight';
import { IndicatorComponent } from './components/presentation/indicator/indicator.component';
import { MonacoEditorModule } from '@materia-ui/ngx-monaco-editor';
import { FormatPathPipe } from './pipes/formatpath/formatpath.pipe';
import { RelativePipe } from './pipes/relative/relative.pipe';
import { OverflowLabelsComponent } from './components/presentation/overflow-labels/overflow-labels.component';
import { OverflowSelectorsComponent } from './components/presentation/overflow-selectors/overflow-selectors.component';
import { PreferencesComponent } from './components/presentation/preferences/preferences.component';
import { HelperComponent } from './components/smart/helper/helper.component';
import { FilterDeletedDatagridRowPipe } from './pipes/filterDeletedDatagridRow/filter-deleted-datagrid-row.pipe';
import { ContentTextFilterComponent } from './components/presentation/content-text-filter/content-text-filter.component';
import { TruncatePipe } from './pipes/truncate/truncate.pipe';
import { StepperComponent } from './components/presentation/stepper/stepper.component';
import { ViewHostDirective } from './directives/view-host/view-host.directive';
import { dynamicComponents } from './dynamic-components';
import { ViewContainerComponent } from './components/view/view-container.component';
import { MissingComponentComponent } from './components/missing-component/missing-component.component';
import { OctantTooltipComponent } from './components/presentation/octant-tooltip/octant-tooltip';
import { BottomPanelComponent } from './components/smart/bottom-panel/bottom-panel.component';
import { DataModule } from '../../data/data.module';
import { OverlayscrollbarsModule } from 'overlayscrollbars-ngx';
import { StringEscapePipe } from './pipes/stringEscape/string.escape.pipe';
import { IconComponent } from './components/presentation/icon/icon.component';
import { FormViewContainerComponent } from './components/form-view-container/form-view-container.component';
import { SignpostComponent } from './components/presentation/signpost/signpost.component';
import { JoinPipe } from './pipes/join/join.pipe';
import { ButtonComponent } from './components/presentation/button/button.component';

@NgModule({
  declarations: [
    AccordionComponent,
    AlertComponent,
    AnnotationsComponent,
    BreadcrumbComponent,
    ButtonGroupComponent,
    CardComponent,
    CardListComponent,
    CodeComponent,
    DropdownComponent,
    ContainersComponent,
    ContentFilterComponent,
    ContextSelectorComponent,
    CytoscapeComponent,
    Cytoscape2Component,
    DatagridComponent,
    DefaultPipe,
    DonutChartComponent,
    EditorComponent,
    ErrorComponent,
    ExpressionSelectorComponent,
    FilterDeletedDatagridRowPipe,
    FiltersComponent,
    FlexlayoutComponent,
    FormComponent,
    FormViewContainerComponent,
    GraphvizComponent,
    HeptagonComponent,
    HeptagonGridComponent,
    HeptagonGridRowComponent,
    HeptagonLabelComponent,
    IFrameComponent,
    IndicatorComponent,
    IconComponent,
    LabelsComponent,
    LabelSelectorComponent,
    LinkComponent,
    ListComponent,
    LoadingComponent,
    LogsComponent,
    ModalComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    PortForwardComponent,
    PortsComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SafePipe,
    AnsiPipe,
    StringEscapePipe,
    SignpostComponent,
    FormatPathPipe,
    RelativePipe,
    TruncatePipe,
    SelectorsComponent,
    SingleStatComponent,
    SliderViewComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TerminalComponent,
    TextComponent,
    TimelineComponent,
    TimestampComponent,
    TitleComponent,
    YamlComponent,
    OverflowLabelsComponent,
    OverflowSelectorsComponent,
    PreferencesComponent,
    HelperComponent,
    ContentTextFilterComponent,
    StepperComponent,
    ViewHostDirective,
    ViewContainerComponent,
    MissingComponentComponent,
    OctantTooltipComponent,
    BottomPanelComponent,
    JoinPipe,
    ButtonComponent,
  ],
  entryComponents: [
    AccordionComponent,
    AlertComponent,
    AnnotationsComponent,
    BreadcrumbComponent,
    ButtonGroupComponent,
    CardComponent,
    CardListComponent,
    CodeComponent,
    DropdownComponent,
    ContainersComponent,
    ContentFilterComponent,
    ContextSelectorComponent,
    CytoscapeComponent,
    Cytoscape2Component,
    DatagridComponent,
    DonutChartComponent,
    EditorComponent,
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
    IconComponent,
    IFrameComponent,
    IndicatorComponent,
    IconComponent,
    LabelsComponent,
    LabelSelectorComponent,
    LinkComponent,
    ListComponent,
    LoadingComponent,
    LogsComponent,
    ModalComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    PortForwardComponent,
    PortsComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SignpostComponent,
    SelectorsComponent,
    SingleStatComponent,
    SliderViewComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TerminalComponent,
    TextComponent,
    TimelineComponent,
    TimestampComponent,
    TitleComponent,
    YamlComponent,
    OverflowLabelsComponent,
    OverflowSelectorsComponent,
    PreferencesComponent,
    HelperComponent,
    ContentTextFilterComponent,
    StepperComponent,
    ViewContainerComponent,
    MissingComponentComponent,
    ButtonComponent,
  ],
  imports: [
    ClarityModule,
    CdsModule,
    CommonModule,
    DataModule,
    FormsModule,
    HighlightModule,
    MonacoEditorModule,
    OverlayscrollbarsModule,
    ReactiveFormsModule,
    ResizableModule,
    RouterModule,
    CdsModule,
  ],
  providers: [highlightProvider(), dynamicComponents()],
  exports: [
    FormatPathPipe,
    AccordionComponent,
    AlertComponent,
    AnnotationsComponent,
    BreadcrumbComponent,
    ButtonGroupComponent,
    CardComponent,
    CardListComponent,
    CodeComponent,
    DropdownComponent,
    ContainersComponent,
    ContentFilterComponent,
    ContextSelectorComponent,
    CytoscapeComponent,
    Cytoscape2Component,
    DatagridComponent,
    DefaultPipe,
    DonutChartComponent,
    EditorComponent,
    ErrorComponent,
    ExpressionSelectorComponent,
    FiltersComponent,
    FlexlayoutComponent,
    FormComponent,
    FormViewContainerComponent,
    GraphvizComponent,
    HeptagonComponent,
    HeptagonGridComponent,
    HeptagonGridRowComponent,
    HeptagonLabelComponent,
    IconComponent,
    IFrameComponent,
    LabelsComponent,
    LabelSelectorComponent,
    LinkComponent,
    ListComponent,
    LoadingComponent,
    LogsComponent,
    ModalComponent,
    ObjectStatusComponent,
    PodStatusComponent,
    PortForwardComponent,
    PortsComponent,
    QuadrantComponent,
    ResourceViewerComponent,
    SignpostComponent,
    SelectorsComponent,
    SliderViewComponent,
    SingleStatComponent,
    SummaryComponent,
    TableComponent,
    TabsComponent,
    TerminalComponent,
    TextComponent,
    TimelineComponent,
    TimestampComponent,
    TitleComponent,
    TruncatePipe,
    YamlComponent,
    OverflowLabelsComponent,
    PreferencesComponent,
    HelperComponent,
    StepperComponent,
    ViewHostDirective,
    ViewContainerComponent,
    OctantTooltipComponent,
    BottomPanelComponent,
    StringEscapePipe,
    JoinPipe,
    ButtonComponent,
  ],
})
export class SharedModule {}
