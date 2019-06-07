import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from 'src/app/models/content';

@Component({
  selector: 'app-content-switcher',
  templateUrl: './content-switcher.component.html',
  styleUrls: ['./content-switcher.component.scss'],
})
export class ContentSwitcherComponent {
  @Input() view: View;
  constructor() {}
}
