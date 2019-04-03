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

  constructor(private notifierService: NotifierService) { }

  ngOnInit() {
    this.subscriptions = [
      this.notifierService.loading.subscribe((loading) => this.loading = loading),
      this.notifierService.error.subscribe((error) => this.error = error),
    ];
  }

  ngOnDestroy(): void {
    this.subscriptions.map(subscription => subscription.unsubscribe());
  }
}
