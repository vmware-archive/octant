import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnInit,
} from '@angular/core';
import { Alert } from '../../../models/content';

const alertLookup = {
  error: 'danger',
  warning: 'warning',
  info: 'info',
  success: 'success',
};

const alertShapes = {
  danger: 'exclamation-circle',
  warning: 'exclamation-triangle',
  info: 'info-circle',
  success: 'check-circle',
};

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AlertComponent implements OnInit {
  @Input() alert: Alert;
  shape = '';
  alertClass = '';
  message = '';
  alertType = '';

  constructor() {}

  ngOnInit(): void {
    if (this.alert) {
      const alertClass = alertLookup[this.alert.type] || alertLookup.error;
      this.shape = alertShapes[alertClass];
      this.alertClass = `alert-${alertClass}`;
      this.message = this.alert.message;
      this.alertType = this.alert.type;
    }
  }
}
