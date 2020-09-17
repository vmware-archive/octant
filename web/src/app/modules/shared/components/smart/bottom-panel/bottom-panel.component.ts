import { Component, ElementRef, OnInit } from '@angular/core';
import { ResizeEvent } from 'angular-resizable-element';
import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';

export const minimizedHeight = '2rem';
export const sliderHeightPropKey = '--slider-height';

const transitionExpr = (from: string, to: string): string => `${from} => ${to}`;

export enum PanelState {
  Open = 'open',
  Closed = 'closed',
}

@Component({
  selector: 'app-bottom-panel',
  templateUrl: './bottom-panel.component.html',
  styleUrls: ['./bottom-panel.component.scss'],
  animations: [
    trigger('toggleState', [
      state(PanelState.Closed, style({ transform: 'rotate(0)' })),
      state(PanelState.Open, style({ transform: 'rotate(-180deg)' })),
      transition(
        transitionExpr(PanelState.Closed, PanelState.Open),
        animate('500ms ease-out')
      ),
      transition(
        transitionExpr(PanelState.Open, PanelState.Closed),
        animate('500ms ease-in')
      ),
    ]),
  ],
})
export class BottomPanelComponent implements OnInit {
  open = false;
  toggleState = PanelState.Closed;
  previousOpenHeight = '50vh';
  resizeEdges = { top: true };

  constructor(private elRef: ElementRef) {}

  ngOnInit() {
    this.setHeight(minimizedHeight);
  }

  parentHeight() {
    return this.elRef.nativeElement.parentNode.offsetHeight + 'px';
  }

  resizeCursors() {
    return {
      topLeft: 'nw-resize',
      topRight: 'ne-resize',
      bottomLeft: 'sw-resize',
      bottomRight: 'se-resize',
      leftOrRight: 'col-resize',
      topOrBottom: this.open ? 'ns-resize' : 'default',
    };
  }

  updateSliderPosition(event: ResizeEvent) {
    if (!this.open) {
      return;
    }

    const panelTop = event.rectangle.top;
    const height = `calc(${this.parentHeight()} - ${panelTop}px)`;
    this.setHeight(height);
    this.previousOpenHeight = height;
  }

  setHeight(height: string) {
    this.elRef.nativeElement.style.setProperty(sliderHeightPropKey, height);
  }

  toggle() {
    this.open = !this.open;
    this.toggleState = this.open ? PanelState.Open : PanelState.Closed;
    this.setHeight(this.open ? this.previousOpenHeight : minimizedHeight);
  }

  gutterClass() {
    return {
      open: this.open,
    };
  }
}
