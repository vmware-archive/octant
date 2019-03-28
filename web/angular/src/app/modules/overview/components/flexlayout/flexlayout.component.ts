import { Component, Input } from '@angular/core';
import { FlexLayoutView } from 'src/app/models/content';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';

@Component({
  selector: 'app-view-flexlayout',
  templateUrl: './flexlayout.component.html',
  styleUrls: ['./flexlayout.component.scss'],
})
export class FlexlayoutComponent {
  @Input() view: FlexLayoutView;
  identifySection = trackByIndex;
}
