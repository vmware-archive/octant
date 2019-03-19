import { Component, Input, OnInit } from '@angular/core';
import { PortForwardView } from 'src/app/models/content';

@Component({
  selector: 'app-view-port-forward',
  templateUrl: './port-forward.component.html',
  styleUrls: ['./port-forward.component.scss'],
})
export class PortForwardComponent implements OnInit {
  @Input() view: PortForwardView;

  constructor() {}

  ngOnInit() {}
}
