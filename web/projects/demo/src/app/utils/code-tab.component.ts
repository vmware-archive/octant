import { Component, Input } from '@angular/core';
import { View } from 'src/app/modules/shared/models/content';

@Component({
  selector: 'app-code-tabs',
  templateUrl: './code-tab.component.html',
  styleUrls: ['./code-tab.component.scss'],
})
export class CodeTabsComponent {
  @Input() preview: View;
  @Input() code: string;
  @Input() json: string;
}
