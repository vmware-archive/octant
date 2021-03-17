/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { InjectionToken, Type } from '@angular/core';
import { AccordionComponent } from './components/presentation/accordion/accordion.component';
import { AnnotationsComponent } from './components/presentation/annotations/annotations.component';
import { IFrameComponent } from './components/presentation/iframe/iframe.component';
import { ResourceViewerComponent } from './components/presentation/resource-viewer/resource-viewer.component';
import { LogsComponent } from './components/smart/logs/logs.component';
import { PodStatusComponent } from './components/presentation/pod-status/pod-status.component';
import { ButtonGroupComponent } from './components/presentation/button-group/button-group.component';
import { CardComponent } from './components/presentation/card/card.component';
import { YamlComponent } from './components/presentation/yaml/yaml.component';
import { CodeComponent } from './components/presentation/code/code.component';
import { EditorComponent } from './components/smart/editor/editor.component';
import { PortForwardComponent } from './components/presentation/port-forward/port-forward.component';
import { ExpressionSelectorComponent } from './components/presentation/expression-selector/expression-selector.component';
import { ListComponent } from './components/presentation/list/list.component';
import { FlexlayoutComponent } from './components/presentation/flexlayout/flexlayout.component';
import { LoadingComponent } from './components/presentation/loading/loading.component';
import { TextComponent } from './components/presentation/text/text.component';
import { QuadrantComponent } from './components/presentation/quadrant/quadrant.component';
import { SelectorsComponent } from './components/presentation/selectors/selectors.component';
import { TerminalComponent } from './components/smart/terminal/terminal.component';
import { DatagridComponent } from './components/presentation/datagrid/datagrid.component';
import { DonutChartComponent } from './components/presentation/donut-chart/donut-chart.component';
import { GraphvizComponent } from './components/presentation/graphviz/graphviz.component';
import { ErrorComponent } from './components/presentation/error/error.component';
import { LabelSelectorComponent } from './components/presentation/label-selector/label-selector.component';
import { PortsComponent } from './components/presentation/ports/ports.component';
import { CardListComponent } from './components/presentation/card-list/card-list.component';
import { LabelsComponent } from './components/presentation/labels/labels.component';
import { LinkComponent } from './components/presentation/link/link.component';
import { SingleStatComponent } from './components/presentation/single-stat/single-stat.component';
import { SummaryComponent } from './components/presentation/summary/summary.component';
import { StepperComponent } from './components/presentation/stepper/stepper.component';
import { TimelineComponent } from './components/presentation/timeline/timeline.component';
import { TimestampComponent } from './components/presentation/timestamp/timestamp.component';
import { ContainersComponent } from './components/presentation/containers/containers.component';
import { ModalComponent } from './components/presentation/modal/modal.component';
import { DropdownComponent } from './components/presentation/dropdown/dropdown.component';
import { IconComponent } from './components/presentation/icon/icon.component';

export interface ComponentMapping {
  [key: string]: Type<any>;
}

const DynamicComponentMapping: ComponentMapping = {
  accordion: AccordionComponent,
  annotations: AnnotationsComponent,
  buttonGroup: ButtonGroupComponent,
  card: CardComponent,
  cardList: CardListComponent,
  codeBlock: CodeComponent,
  containers: ContainersComponent,
  donutChart: DonutChartComponent,
  dropdown: DropdownComponent,
  editor: EditorComponent,
  expressionSelector: ExpressionSelectorComponent,
  graphviz: GraphvizComponent,
  flexlayout: FlexlayoutComponent,
  labels: LabelsComponent,
  labelSelector: LabelSelectorComponent,
  loading: LoadingComponent,
  error: ErrorComponent,
  iframe: IFrameComponent,
  link: LinkComponent,
  list: ListComponent,
  logs: LogsComponent,
  modal: ModalComponent,
  podStatus: PodStatusComponent,
  portforward: PortForwardComponent,
  ports: PortsComponent,
  quadrant: QuadrantComponent,
  resourceViewer: ResourceViewerComponent,
  selectors: SelectorsComponent,
  singleStat: SingleStatComponent,
  stepper: StepperComponent,
  summary: SummaryComponent,
  table: DatagridComponent,
  terminal: TerminalComponent,
  text: TextComponent,
  timeline: TimelineComponent,
  timestamp: TimestampComponent,
  yaml: YamlComponent,
  icon: IconComponent,
};

export const DYNAMIC_COMPONENTS_MAPPING = new InjectionToken<ComponentMapping>(
  'dynamicComponentsMapping'
);

export const dynamicComponents = () => {
  return {
    provide: DYNAMIC_COMPONENTS_MAPPING,
    useValue: DynamicComponentMapping,
  };
};
