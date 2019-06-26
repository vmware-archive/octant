import { Component, OnInit, Input } from '@angular/core';
import { Workload } from '../../types';

@Component({
  selector: 'app-workload-card',
  templateUrl: './workload-card.component.html',
  styleUrls: ['./workload-card.component.scss']
})
export class WorkloadCardComponent implements OnInit {
  @Input() workload: Workload;

  constructor() { }

  ngOnInit() {
  }

}
