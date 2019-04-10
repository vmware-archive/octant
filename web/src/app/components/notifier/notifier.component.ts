import { Component, OnDestroy, OnInit } from '@angular/core';
import { NotifierService } from 'src/app/services/notifier/notifier.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-notifier',
  templateUrl: './notifier.component.html',
  styleUrls: ['./notifier.component.scss']
})
export class NotifierComponent implements OnInit, OnDestroy {
  private subscriptions: Subscription[];
  loading = false;
  error: string;
  warning: string;

  constructor(private notifierService: NotifierService) { }

  ngOnInit() {
    this.subscriptions = [
      this.notifierService.loading.subscribe((loading) => this.loading = loading),
      this.notifierService.error.subscribe((error) => this.error = error),
      this.notifierService.warning.subscribe((warning) => this.warning = warning),
    ];
  }

  onWarningClose() {
    this.warning = '';
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(subscription => subscription.unsubscribe());
  }
}
