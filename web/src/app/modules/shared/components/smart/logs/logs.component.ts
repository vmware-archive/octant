// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterViewChecked,
  Component,
  ElementRef,
  Input,
  IterableDiffer,
  IterableDiffers,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import {
  LogEntry,
  LogsView,
  View,
} from 'src/app/modules/shared/models/content';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  PodLogsService,
  PodLogsStreamer,
} from 'src/app/modules/shared/pod-logs/pod-logs.service';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
})
export class LogsComponent implements OnInit, OnDestroy, AfterViewChecked {
  v: LogsView;

  @Input() set view(v: View) {
    this.v = v as LogsView;
  }
  get view() {
    return this.v;
  }

  private logStream: PodLogsStreamer;
  scrollToBottom = false;

  private containerLogsDiffer: IterableDiffer<LogEntry>;
  @ViewChild('scrollTarget', { static: true }) scrollTarget: ElementRef;
  containerLogs: LogEntry[] = [];

  selectedContainer = '';
  shouldDisplayTimestamp = true;

  constructor(
    private podLogsService: PodLogsService,
    private iterableDiffers: IterableDiffers
  ) {}

  ngOnInit() {
    this.containerLogsDiffer = this.iterableDiffers
      .find(this.containerLogs)
      .create();
    if (this.v) {
      if (this.v.config.containers && this.v.config.containers.length > 0) {
        this.selectedContainer = this.v.config.containers[0];
      }
      this.startStream();
    }
  }

  onContainerChange(containerSelection: string): void {
    this.selectedContainer = containerSelection;
    if (this.logStream) {
      this.containerLogs = [];
      this.logStream.close();
      this.logStream = null;
    }
    this.startStream();
  }

  toggleTimestampDisplay(): void {
    this.shouldDisplayTimestamp = !this.shouldDisplayTimestamp;
  }

  startStream() {
    const namespace = this.v.config.namespace;
    const pod = this.v.config.name;
    const container = this.selectedContainer;
    if (namespace && pod && container) {
      this.logStream = this.podLogsService.createStream(
        namespace,
        pod,
        container
      );
      this.logStream.logEntries
        .pipe(untilDestroyed(this))
        .subscribe((entries: LogEntry[]) => {
          this.containerLogs = entries;
        });
    }
  }

  identifyLog(index: number, item: LogEntry) {
    return `${item.timestamp}-${item.message}`;
  }

  // Note(marlon): to determine if we should continue tailing
  // the incoming logs
  onScroll(evt: { target: HTMLDivElement }) {
    const { target } = evt;
    const { clientHeight, scrollHeight, scrollTop, offsetHeight } = target;
    this.scrollToBottom = false;
    if (scrollHeight <= clientHeight) {
      // Not scrollable
      return;
    }
    if (scrollTop < scrollHeight - offsetHeight) {
      // Not at the bottom
      return;
    }
    this.scrollToBottom = true;
  }

  ngAfterViewChecked() {
    const change = this.containerLogsDiffer.diff(this.containerLogs);
    if (change && this.scrollToBottom) {
      const { nativeElement } = this.scrollTarget;
      nativeElement.scrollTop = nativeElement.scrollHeight;
    }
  }

  ngOnDestroy(): void {
    if (this.logStream) {
      this.logStream.close();
      this.logStream = null;
    }
  }
}
