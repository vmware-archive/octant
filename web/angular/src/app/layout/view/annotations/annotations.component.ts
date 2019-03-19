import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { AnnotationsView } from 'src/app/models/content';

@Component({
  selector: 'app-view-annotations',
  templateUrl: './annotations.component.html',
  styleUrls: ['./annotations.component.scss'],
})
export class AnnotationsComponent implements OnChanges {
  @Input() view: AnnotationsView;

  annotations: { [key: string]: string };
  annotationKeys: string[];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as AnnotationsView;
      this.annotations = view.config.annotations;
      this.annotationKeys = Object.keys(this.annotations);
    }
  }
}
