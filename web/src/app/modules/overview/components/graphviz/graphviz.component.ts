import { AfterViewChecked, Component, ElementRef, Input, ViewChild, ViewEncapsulation } from '@angular/core';
import * as dot from 'graphlib-dot';
import { GraphvizView } from 'src/app/models/content';
import { DagreService } from '../../services/dagre/dagre.service';

@Component({
  selector: 'app-view-graphviz',
  template: `
    <div class="graphviz" #viewer>
    </div>
  `,
  styleUrls: ['./graphviz.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class GraphvizComponent implements AfterViewChecked {
  @ViewChild('viewer') private viewer: ElementRef;
  @Input() view: GraphvizView;

  constructor(private dagreService: DagreService) { }

  ngAfterViewChecked() {
    if (this.view) {
      const current = this.view.config.dot;
      const g = dot.read(current);
      this.dagreService.render(this.viewer, g);
    }
  }
}
