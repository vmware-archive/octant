import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnInit,
} from '@angular/core';

import '@cds/core/alert/register.js';
import { Alert } from '../../../models/content';

const alertLookup = {
  error: 'danger',
  warning: 'warning',
  info: 'info',
  success: 'success',
};

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AlertComponent implements OnInit {
  @Input() alert: Alert;
  message = '';
  status = '';
  type = '';
  closable = false;
  buttonGroup = null;
  showAlert = false;

  constructor() {}

  ngOnInit(): void {
    if (this.alert) {
      this.type = this.alert.type;
      this.message = this.alert.message;
      this.status = alertLookup[this.alert.status] || alertLookup.error;
      this.closable = this.alert.closable;
      this.buttonGroup = this.alert.buttonGroup;
      this.showAlert = true;
    }
  }

  close(): void {
    this.showAlert = false;
  }
}
