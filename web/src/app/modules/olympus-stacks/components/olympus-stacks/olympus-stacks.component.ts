import { Component, OnInit } from '@angular/core';
import { WorkloadList } from '../../types';
import example from './examplepayload';

@Component({
  selector: 'app-olympus-stacks',
  templateUrl: './olympus-stacks.component.html',
  styleUrls: ['./olympus-stacks.component.scss']
})
export class OlympusStacksComponent implements OnInit {
  workloadsLists: WorkloadList[];

  constructor() {
    this.workloadsLists = example;
  }

  ngOnInit() {
  }
}
