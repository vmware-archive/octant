import { Component, OnInit, Input } from '@angular/core';
import { WorkloadList, Workload } from '../../types';

@Component({
  selector: 'app-workload-list',
  templateUrl: './workload-list.component.html',
  styleUrls: ['./workload-list.component.scss']
})
export class WorkloadListComponent implements OnInit {
  @Input() workloadList: WorkloadList;

  constructor() { }

  ngOnInit() {
  }

  identifyWorkload({ name, revision}: Workload) {
    return `${name}-${revision}`;
  }
}
