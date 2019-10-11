// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterViewChecked,
  Component,
  ElementRef,
  Input,
  OnDestroy,
  OnInit,
  ViewChild,
  IterableDiffers,
  IterableDiffer,
} from '@angular/core';
import { LogsView, LogEntry } from 'src/app/models/content';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  PodLogsService,
  PodLogsStreamer,
} from 'src/app/services/pod-logs/pod-logs.service';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
})
export class LogsComponent implements OnInit, OnDestroy, AfterViewChecked {
  private logStream: PodLogsStreamer;
  scrollToBottom = false;

  private containerLogsDiffer: IterableDiffer<LogEntry>;
  @Input() view: LogsView;
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
    if (this.view) {
      if (
        this.view.config.containers &&
        this.view.config.containers.length > 0
      ) {
        this.selectedContainer = this.view.config.containers[0];
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
    const namespace = this.view.config.namespace;
    const pod = this.view.config.name;
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
