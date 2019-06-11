import {
  AfterViewChecked,
  Component,
  ElementRef,
  HostListener,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import * as d3 from 'd3';
import {Edge} from 'dagre';
import * as dagreD3 from 'dagre-d3';
import _ from 'lodash';
import {Node, ResourceViewerView} from 'src/app/models/content';
import {ResourceNode} from 'src/app/models/resource-node';

import {DagreService} from '../../services/dagre/dagre.service';

@Component({
  selector: 'app-view-resource-viewer',
  templateUrl: './resource-viewer.component.html',
  styleUrls: ['./resource-viewer.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ResourceViewerComponent implements OnChanges, AfterViewChecked {

  @Input() view: ResourceViewerView;
  currentView: ResourceViewerView;
  selected: string;
  selectedNode: Node;
  @ViewChild('viewer') private viewer: ElementRef;
  private runUpdate = false;
  private hasDrawn = false;

  constructor(private dagreService: DagreService) {
  }

  ngOnChanges(changes: SimpleChanges): void {

    this.currentView = changes.view.currentValue as ResourceViewerView;

    const isEqual = _.isEqual(changes.view.currentValue, changes.view.previousValue);

    if (changes.view.isFirstChange()) {
      this.select(this.currentView.config.selected);
      this.runUpdate = true;
    } else if (!isEqual) {
      this.select(this.selected);
      this.runUpdate = true;
    } else if (!this.hasDrawn) {
      this.runUpdate = true;
    }
  }

  ngAfterViewChecked(): void {
    this.updateGraph();
  }

  @HostListener('window:resize') onResize() {
    if (this.viewer) {
      this.runUpdate = true;
      this.updateGraph();
    }
  }

  edges(): Array<Edge> {
    const adjacencyList = this.currentView.config.edges;
    const edges: Array<Edge> = [];

    if (adjacencyList) {
      for (const [node, nodeEdges] of Object.entries(adjacencyList)) {
        edges.push(
          ...nodeEdges.map((e) => [
            node,
            e.node,
            {
              arrowhead: 'undirected',
              lineCurve: d3.curveBasis,
            },
          ])
        );
      }
    }

    return edges;
  }

  nodes() {
    const objects = this.currentView.config.nodes;

    const nodes: { [key: string]: dagreD3.Label } = {};

    if (objects) {
      for (const [id, object] of Object.entries(objects)) {
        const isSelected = id === this.selected;
        nodes[id] = new ResourceNode(id, object, isSelected).toDescriptor();
      }
    }

    return nodes;
  }

  private updateGraph() {
    if (!this.runUpdate || !this.currentView.config.nodes) {
      return;
    }

    try {
      const g = new dagreD3.graphlib.Graph().setGraph({});

      for (const [id, label] of Object.entries(this.nodes())) {
        g.setNode(id, label);
      }

      g.nodes().forEach((v: any) => {
        const node = g.node(v);
        node.rx = node.ry = 4;
      });

      this.edges().forEach((edge) => g.setEdge(edge[0], edge[1], edge[2]));

      this.dagreService.render(this.viewer, g);

      const svg = d3.select('.viewer svg');
      const nodes = svg.selectAll('g.node');

      if (nodes.nodes().length > 0) {
        this.runUpdate = false;
        this.hasDrawn = true;

        nodes.on('click', (id: string) => {
          this.onClick(id);
        });
      }
    } catch (error) {
      console.log(`render resource viewer failed ${error}`, this.currentView.config);
    }

  }

  private onClick(id: string) {
    this.select(id);
  }

  private select(id: string) {
    this.runUpdate = true;
    this.selected = id;

    const nodes = this.currentView.config.nodes;

    if (nodes && nodes[id]) {
      this.selectedNode = nodes[id];
    }
  }
}
