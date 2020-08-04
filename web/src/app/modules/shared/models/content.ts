// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

export interface ContentResponse {
  content: Content;
}

export interface PathItem {
  title: string;
  url?: string;
}

export interface Content {
  extensionComponent: ExtensionView;
  viewComponents: View[];
  title: View[];
  buttonGroup?: ButtonGroupView;
}

export interface Metadata {
  type: string;
  title?: View[];
  accessor?: string;
}

export interface View {
  metadata: Metadata;
  totalItems?: number;
}

export interface TitleMetadata {
  type: 'text' | 'link';
  title?: TitleView[];
  accessor?: string;
}

export interface TitleView {
  metadata: TitleMetadata;
}

export interface AnnotationsView extends View {
  config: {
    annotations: { [key: string]: string };
  };
}

export interface Alert {
  type: string;
  message: string;
}

export interface CardView extends View {
  config: {
    body: View;
    actions: Action[];
    alert?: Alert;
  };
}

export interface CardListView extends View {
  config: {
    cards: CardView[];
  };
}

export interface ContainerDef {
  name: string;
  image: string;
}

export interface ContainersView extends View {
  config: {
    containers: ContainerDef[];
  };
}

export interface DonutChartLabels {
  plural: string;
  singular: string;
}

export interface DonutChartView extends View {
  config: {
    segments: DonutSegment[];
    labels: DonutChartLabels;
    size: number;
  };
}

export interface GraphvizView extends View {
  config: {
    dot: string;
  };
}

export interface FlexLayoutItem {
  width: number;
  height: number;
  view: View;
}

export interface Confirmation {
  title: string;
  body: string;
}

export interface Button {
  payload: {};
  name: string;
  confirmation?: Confirmation;
}

export interface ButtonGroupView extends View {
  config: {
    buttons: Button[];
  };
}

export interface FlexLayoutView extends View {
  config: {
    sections: FlexLayoutItem[][];
    buttonGroup: ButtonGroupView;
  };
}

export interface GridAction {
  name: string;
  actionPath: string;
  payload: {};
  confirmation?: Confirmation;
  type: string;
}

export interface GridActionsView extends View {
  config: {
    actions: GridAction[];
  };
}

export interface LabelsView extends View {
  config: {
    labels: { [key: string]: string };
  };
}

export interface LinkView extends View {
  config: {
    ref: string;
    value: string;
    status?: number;
    statusDetail?: View;
  };
}

export interface ListView extends View {
  config: {
    iconName: string;
    items: View[];
  };
}

export interface ExpressionSelectorView extends View {
  config: {
    key: string;
    operator: string;
    values: string[];
  };
}

export interface LabelSelectorView extends View {
  config: {
    key: string;
    value: string;
  };
}

export interface SingleStatView extends View {
  config: {
    title: string;
    value: {
      text: string;
      color: string;
    };
  };
}

export interface PodSummary {
  details: View[];
  status: string;
}

export interface PodStatusView extends View {
  config: {
    pods: { [key: string]: PodSummary };
  };
}

export interface PortForwardPortSpec {
  local: number;
  remote: number;
}

export interface PortForwardTarget {
  apiVersion: string;
  kind: string;
  namespace: string;
  name: string;
}

export interface PortForwardView extends View {
  config: {
    text: string;
    id: string;
    action: string;
    status: string;
    ports: PortForwardPortSpec[];
    target: PortForwardTarget;
  };
}

export interface QuadrantValue {
  value: string;
  label: string;
}

export interface QuadrantView extends View {
  config: {
    nw: QuadrantValue;
    ne: QuadrantValue;
    sw: QuadrantValue;
    se: QuadrantValue;
  };
}

export interface Edge {
  node: string;
  edge: string;
}

export interface Node {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
  details: View;
  path: LinkView;
}

export interface ResourceViewerView extends View {
  config: {
    edges: { [key: string]: Edge[] };
    nodes: Node[];
    selected: string;
  };
}

export interface SelectorsView extends View {
  config: {
    selectors: Array<ExpressionSelectorView | LabelSelectorView>;
  };
}

export interface SummaryItem {
  header: string;
  content: View;
}

export interface ActionField {
  configuration: any;
  label: string;
  name: string;
  type: string;
  value: any;
  placeholder: string;
  error: string;
  validators: string[];
}

export interface ActionForm {
  fields: ActionField[];
}

export interface Action {
  name: string;
  title: string;
  form: ActionForm;
}

export interface SummaryView extends View {
  config: {
    sections: SummaryItem[];
    actions: Action[];
    alert?: Alert;
  };
}

export interface TableView extends View {
  config: {
    columns: TableColumn[];
    rows: TableRow[];
    emptyContent: string;
    loading: boolean;
    filters: TableFilters;
  };
}

export interface TableFilters {
  [key: string]: TableFilter;
}

export interface TableFilter {
  values: string[];
  selected: string[];
}

export interface TableRow {
  [key: string]: View;
}

export interface TableRowWithMetadata {
  data: TableRow;
  actions?: GridAction[];
  isDeleted: boolean;
}

export interface TableColumn {
  name: string;
  accessor: string;
}

export interface TextView extends View {
  config: {
    value: string;
    isMarkdown?: boolean;
    status?: number;
  };
}

export interface TimestampView extends View {
  config: {
    timestamp: number;
  };
}

export interface DonutSegment {
  count: number;
  status: string;
}

export interface Series {
  name: string;
  value: number;
  label: string;
  color: string;
}

export interface BulletBand {
  min: number;
  max: number;
  color: string;
  label: string;
}

export interface Resource {
  bands: BulletBand[];
  measure: number;
  measureLabel: string;
  label: string;
}

export interface WorkloadView extends View {
  config: {
    name: string;
    iconName: string;
    segments: DonutSegment[];
    memory: Resource;
    cpu: Resource;
  };
}

export interface WorkloadListView extends View {
  config: {
    workloads: WorkloadView;
  };
}

export interface YAMLView extends View {
  config: {
    data: string;
  };
}

export interface LogsView extends View {
  config: {
    namespace: string;
    name: string;
    containers: string[];
  };
}

export interface LogEntry {
  timestamp: string;
  message: string;
  container: string;
}

export interface LogResponse {
  entries: LogEntry[];
}

export interface TerminalOutput {
  scrollback: string;
  line: string;
  exitMessage: string;
}

export interface TerminalDetail {
  container: string;
  command: string;
  active: boolean;
}

export interface TerminalView extends View {
  config: {
    namespace: string;
    name: string;
    podName: string;
    terminal: TerminalDetail;
    containers: string[];
  };
}

export interface EditorView extends View {
  config: {
    value: string;
    language: string;
    readOnly: boolean;
    metadata: { [key: string]: string };
    submitAction: string;
    submitLabel: string;
  };
}

export interface Port extends View {
  config: {
    port: number;
    protocol: string;
    state: Partial<{
      id: string;
      isForwarded: boolean;
      isForwardable: boolean;
      port: number;
    }>;
    buttonGroup: ButtonGroupView;
  };
}

export interface PortsView extends View {
  config: {
    ports: Port[];
  };
}

export interface LoadingView extends View {
  config: {
    value: string;
  };
}

export interface ErrorView extends View {
  config: {
    data: string;
  };
}

export interface IFrameView extends View {
  config: {
    url: string;
    title: string;
  };
}

export interface ExtensionTab {
  tab: View;
  payload: { [key: string]: string };
}

export interface ExtensionView extends View {
  config: {
    tabs: ExtensionTab[];
  };
}

export interface CodeView extends View {
  config: {
    value: string;
  };
}

export interface StepItem {
  name: string;
  form: ActionForm;
  title: string;
  description: string;
}

export interface StepperView extends View {
  config: {
    action: string;
    steps: StepItem[];
  };
}
