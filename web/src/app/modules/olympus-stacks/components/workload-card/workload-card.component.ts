import { Component, OnInit, Input, OnChanges, SimpleChanges } from '@angular/core';
import { Workload } from '../../types';
import moment from 'moment';

@Component({
  selector: 'app-workload-card',
  templateUrl: './workload-card.component.html',
  styleUrls: ['./workload-card.component.scss']
})
export class WorkloadCardComponent implements OnInit, OnChanges {
  @Input() workload: Workload;
  editMode = false;
  humanReadableTimestamp: string;
  newRevisionInput: string;

  constructor() {}
  ngOnInit() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.workload.currentValue) {
      const { lastUpdated } = changes.workload.currentValue as Workload;
      const m = moment(lastUpdated);
      const date = m.format('ddd, MMMM Do h:mma');
      this.humanReadableTimestamp = `${date} (${m.fromNow()})`;
    }
  }

  saveChanges() {
    if (!this.newRevisionInput) {
      return;
    }
    const warningText = `Please confirming setting ${this.workload.name}'s revision to:
      ${this.newRevisionInput}
    `
    if (!confirm(warningText)) {
      return;
    }
    this.editMode = false;
    this.newRevisionInput = '';
  }

  pinToStack() {
    const warningText = `Please confirming setting ${this.workload.name}'s revision to:
      ${this.workload.revision}
    `;
    if (!confirm(warningText)) {
      return;
    }
  }

  cancelChanges() {
    this.newRevisionInput = '';
    this.editMode = false;
  }
}
