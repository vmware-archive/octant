import {
  Component,
  Input,
  OnChanges,
  OnInit,
  SimpleChanges,
  ViewChild,
  ElementRef,
  ViewEncapsulation,
  DoCheck,
  AfterViewChecked,
} from '@angular/core';
import * as dagreD3 from 'dagre-d3';
import { ResourceViewerView } from 'src/app/models/content';
import { graphlib } from 'dagre-d3';
import * as d3 from 'd3';
import { zoom } from 'd3';

interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
}

class ResourceNode {
  constructor(
    private readonly id: string,
    private readonly object: ResourceObject,
    private readonly isSelected: boolean
  ) {}

  toDescriptor(): dagreD3.Label {
    let nodeClass = `node-${this.object.status}`;
    if (this.isSelected) {
      nodeClass += ` selected`;
    }

    return {
      id: this.id,
      label: `${this.title()}${this.subTitle()}`,
      labelType: 'html',
      class: `${nodeClass}`,
    };
  }

  title(): string {
    return `<div class="resource-name">${this.object.name}</div>`;
  }

  subTitle(): string {
    return `<div class="resource-type">${this.object.apiVersion} ${this.object.kind}</div>`;
  }
}


@Component({
  selector: 'app-view-resource-viewer',
  template: `
    <div class="resourceViewer" #viewer>
      <svg:svg class="dagre-d3" #parent>
        <g #container></g>
      </svg:svg>
    </div>
  `,
  styleUrls: ['./resource-viewer.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ResourceViewerComponent implements OnInit, OnChanges, AfterViewChecked {
  constructor() {}
  @ViewChild('viewer') private viewer: ElementRef;
  @ViewChild('parent') private parent: ElementRef;
  @ViewChild('container') private container: ElementRef;

  @Input() view: ResourceViewerView;

  private changed = false;

  ngOnInit() {}

  ngOnChanges(changes: SimpleChanges): void {
    const a = JSON.stringify(changes.view.previousValue);
    const b = JSON.stringify(changes.view.currentValue);
    if (a !== b) {
      this.changed = true;
      const view = changes.view.currentValue as ResourceViewerView;

      const objects = view.config.nodes;
      const adjacencyList = view.config.edges;

      const nodes: { [key: string]: dagreD3.Label } = {};

      for (const [id, object] of Object.entries(objects)) {
        nodes[id] = new ResourceNode(id, object, false).toDescriptor();
      }
      const edges = [];
      if (adjacencyList) {
        for (const [node, nodeEdges] of Object.entries(adjacencyList)) {
          edges.push(
            ...nodeEdges.map((e) => [
              node,
              e.node,
              {
                arrowhead: 'undirected',
                arrowheadStyle: 'fill: rgba(173, 187, 196, 0.3)',
              },
            ])
          );
        }
      }

      const g = new dagreD3.graphlib.Graph().setGraph({
        align: 'DR',
      });

      for (const [id, label] of Object.entries(nodes)) {
        g.setNode(id, label);
      }

      g.nodes().forEach((v) => {
        const node = g.node(v);
        node.rx = node.ry = 4;
      });

      edges.forEach((edge) => g.setEdge(edge[0], edge[1], edge[2]));

      const containerElement = this.container.nativeElement;
      const inner = d3.select(containerElement);

      const render = new dagreD3.render();
      // @ts-ignore
      render(inner, g);
    }
  }

  ngAfterViewChecked() {
    // this translates/scales after the view has been shown. it is not optimal, and should
    // be performed sooner.
    const viewerElement = this.viewer.nativeElement;
    const viewerHeight = viewerElement.offsetHeight;
    const viewerWidth = viewerElement.offsetWidth;
    if (viewerHeight < 1 || viewerWidth < 1 || !this.changed) {
      return;
    }

    this.resize();
  }

  resize() {
    const viewerElement = this.viewer.nativeElement;

    const parentElement = this.parent.nativeElement;
    const svg = d3.select(parentElement);

    const containerElement = this.container.nativeElement;
    const inner = d3.select(containerElement);

    const viewerHeight = viewerElement.offsetHeight;
    const viewerWidth = viewerElement.offsetWidth;
    const { height, width } = parentElement.getBBox();

    svg.attr('height', viewerHeight);
    svg.attr('width', viewerWidth);

    const bounds = inner.node().getBBox();
    const parent = inner.node().parentElement;
    const fullWidth = parent.clientWidth;
    const fullHeight = parent.clientHeight;
    const innerWidth = bounds.width;
    const innerHeight = bounds.height;
    const midX = bounds.x + innerWidth / 2;
    const midY = bounds.y + innerHeight / 2;
    const scale = 0.4 / Math.max(width / fullWidth, height / fullHeight);
    const translate = [fullWidth / 2 - scale * midX, fullHeight / 2 - scale * midY];

    // @ts-ignore
    inner.attr('transform', `translate(${translate[0]}, ${translate[1]}) scale(${scale})`);

    this.changed = false;
  }
}
