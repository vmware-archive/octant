// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectorRef,
  Component,
  ElementRef,
  isDevMode,
  OnDestroy,
  OnInit,
  Renderer2,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import {
  Node,
  ResourceViewerView,
} from 'src/app/modules/shared/models/content';
import { ElementsDefinition, Stylesheet } from 'cytoscape';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ELEMENTS_STYLE, ELEMENTS_STYLE_DARK } from './octant.style';
import { Router } from '@angular/router';
import { ThemeService } from '../../../services/theme/theme.service';
import { Subscription } from 'rxjs';
import { ResizeEvent } from 'angular-resizable-element';

const statusColorCodes = {
  ok: '#60b515',
  warning: '#f57600',
  error: '#e12200',
};

const edgeColorCode = '#003d79';

const defaultZoom = {
  min: 0.075,
  max: 4.0,
};

@Component({
  selector: 'app-view-resource-viewer',
  templateUrl: './resource-viewer.component.html',
  styleUrls: ['./resource-viewer.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ResourceViewerComponent
  extends AbstractViewComponent<ResourceViewerView>
  implements OnInit, OnDestroy
{
  selectedNodeId: string;
  private subscriptionTheme: Subscription;
  resizeEdges = { left: true, right: true };
  startPosition: number;

  @ViewChild('resourceViewer')
  resourceViewer: ElementRef;

  @ViewChild('viewContainer')
  viewContainer: ElementRef;

  @ViewChild('statusContainer')
  statusContainer: ElementRef;

  layout = {
    name: 'dagre',
    padding: 0,
    nodeSep: 50,
    rankSep: 150,
    rankDir: 'TB',
    directed: true,
    animate: false,
  };

  zoom = defaultZoom;

  style: Stylesheet[] = ELEMENTS_STYLE;
  graphData: ElementsDefinition;

  constructor(
    private renderer: Renderer2,
    private router: Router,
    private themeService: ThemeService,
    private cdr: ChangeDetectorRef
  ) {
    super();
  }

  ngOnInit(): void {
    this.subscriptionTheme = this.themeService.themeType.subscribe(() => {
      this.style = this.themeService.isLightThemeEnabled()
        ? ELEMENTS_STYLE
        : ELEMENTS_STYLE_DARK;
      this.cdr.detectChanges();
    });
  }

  ngOnDestroy(): void {
    this.subscriptionTheme?.unsubscribe();
  }

  update() {
    const nodes: Node[] = this.v.config.nodes;
    if (nodes && Object.keys(nodes).length > 0) {
      const selection = this.v.config?.selected
        ? this.v.config.selected
        : Object.keys(nodes)[0];

      this.graphData = this.generateGraphData();
      this.selectNode(selection);

      if (isDevMode()) {
        console.log(
          'Resource view data:',
          JSON.stringify((this.view as ResourceViewerView).config)
        );
      }
    }
  }

  generateGraphData() {
    return {
      nodes: this.nodes(),
      edges: this.edges(),
    };
  }

  nodes() {
    if (!this.v.config.nodes) {
      return [];
    }

    const nodes = Object.entries(this.v.config.nodes).map(([name, details]) => {
      const colorCode =
        statusColorCodes[details.status] || statusColorCodes.error;

      return {
        data: {
          id: name,
          label1: this.getLabel(details.name, 20),
          label2: this.getLabel(`${details.apiVersion} ${details.kind}`, 36),
          weight: 100,
          status: details.status,
          colorCode,
        },
      };
    });

    return Array.prototype.concat(...nodes);
  }

  edges() {
    if (!this.v.config.edges) {
      return [];
    }

    const edges = Object.entries(this.v.config.edges).map(([parent, maps]) => {
      return maps.map(edge => {
        return {
          data: {
            source: parent,
            target: edge.node,
            colorCode: edgeColorCode,
            strength: 10,
          },
        };
      });
    });

    return Array.prototype.concat(...edges);
  }

  nodeChange(event) {
    this.selectNode(event.id);
  }

  selectNode(id: string) {
    this.selectedNodeId = id;
  }

  selectedNode(): string {
    return this.v?.config?.nodes[this.selectedNodeId];
  }

  openNode(event) {
    const node = this.v.config.nodes[event.id];
    if (node && node.path) {
      this.router.navigateByUrl(node.path.config.ref);
    }
  }

  getLabel(label: string, length: number): string {
    return label.length > length ? label.substring(0, length) + '...' : label;
  }

  resizeCursors() {
    return {
      topLeft: 'nw-resize',
      topRight: 'ne-resize',
      bottomLeft: 'sw-resize',
      bottomRight: 'se-resize',
      leftOrRight: 'ew-resize',
      topOrBottom: 'ns-resize',
    };
  }

  resizeEnd() {
    this.zoom = Object.assign({}, defaultZoom); // update without layout recalc
  }

  resizeStart() {
    this.startPosition = this.viewContainer.nativeElement.offsetWidth;
  }

  updateSliderPosition(event: ResizeEvent) {
    const parentWidth = this.resourceViewer.nativeElement.offsetWidth - 4;
    const sliderOffset = event.edges.left as number;
    const leftSize = Math.max(
      30,
      Math.min(80, (100 * (this.startPosition + sliderOffset)) / parentWidth)
    );

    this.renderer.setStyle(
      this.viewContainer.nativeElement,
      'width',
      `${leftSize}%`
    );
    this.renderer.setStyle(
      this.statusContainer.nativeElement,
      'width',
      `${100 - leftSize}%`
    );
  }
}
