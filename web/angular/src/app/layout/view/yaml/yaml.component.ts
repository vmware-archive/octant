import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { YAMLView } from 'src/app/models/content';

@Component({
  selector: 'app-view-yaml',
  templateUrl: './yaml.component.html',
  styleUrls: ['./yaml.component.scss'],
})
export class YamlComponent implements OnChanges {
  @Input() view: YAMLView;

  source: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as YAMLView;
      this.source = view.config.data;
    }
  }
}
