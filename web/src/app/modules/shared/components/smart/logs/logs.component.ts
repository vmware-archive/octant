// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterContentChecked,
  AfterViewChecked,
  Component,
  ElementRef,
  Input,
  IterableDiffer,
  IterableDiffers,
  OnDestroy,
  OnInit,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import {
  LogEntry,
  LogsView,
  View,
} from 'src/app/modules/shared/models/content';
import {
  PodLogsService,
  PodLogsStreamer,
} from 'src/app/modules/shared/pod-logs/pod-logs.service';
import { formatDate } from '@angular/common';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class LogsComponent
  implements OnInit, OnDestroy, AfterContentChecked, AfterViewChecked {
  v: LogsView;

  @Input() set view(v: View) {
    this.v = v as LogsView;
  }
  get view() {
    return this.v;
  }

  private logStream: PodLogsStreamer;
  scrollToBottom = true;

  private containerLogsDiffer: IterableDiffer<LogEntry>;
  @ViewChild('scrollTarget', { static: true }) scrollTarget: ElementRef;
  containerLogs: LogEntry[] = [];

  selectedContainer = '';
  shouldDisplayTimestamp = true;
  showOnlyFiltered = false;
  filterText = '';
  oldFilterText = '';
  currentSelection = 0;
  totalSelections = 0;
  timeFormat = 'MMM d, y h:mm:ss a z';

  private logSubscription: Subscription;

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
    this.updateSelectedCount();
    this.scrollToHighlight(0, 0);
  }

  toggleShowOnlyFiltered(): void {
    this.showOnlyFiltered = !this.showOnlyFiltered;
    this.scrollToHighlight(0, 0);
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
      this.logSubscription = this.logStream.logEntries.subscribe(
        (entries: LogEntry[]) => {
          this.containerLogs = entries;
          this.updateSelectedCount();
        }
      );
    }
  }

  identifyLog(index: number, item: LogEntry) {
    return `${item.timestamp}-${item.message}`;
  }

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

  ngAfterContentChecked() {
    if (this.filterText !== this.oldFilterText) {
      this.updateSelectedCount();
    }
  }

  ngAfterViewChecked() {
    const change = this.containerLogsDiffer.diff(this.containerLogs);
    if (change && this.scrollToBottom) {
      const { nativeElement } = this.scrollTarget;
      nativeElement.scrollTop = nativeElement.scrollHeight;
    }
    if (this.filterText !== this.oldFilterText) {
      this.oldFilterText = this.filterText;
      this.scrollToHighlight(0, 0);
    }
  }

  ngOnDestroy(): void {
    if (this.logStream) {
      this.logStream.close();
      this.logStream = null;
    }

    if (this.logSubscription) {
      this.logSubscription.unsubscribe();
    }
  }

  public highlightText(text: string) {
    let highlighted = text;

    if (this.filterText) {
      highlighted = text.replace(new RegExp(this.filterText, 'g'), match => {
        return '<span class="highlight">' + match + '</span>';
      });
    }
    return `${highlighted}`;
  }

  public filterFunction(logs: LogEntry[]): LogEntry[] {
    if (this.showOnlyFiltered) {
      if (this.shouldDisplayTimestamp) {
        return logs.filter(
          log =>
            log.message.match(new RegExp(this.filterText, 'g')) ||
            formatDate(log.timestamp, this.timeFormat, 'en-US').match(
              new RegExp(this.filterText, 'g')
            )
        );
      }
      return logs.filter(log =>
        log.message.match(new RegExp(this.filterText, 'g'))
      );
    }

    return logs;
  }

  onPreviousHighlight(): void {
    if (this.currentSelection > 0) {
      this.scrollToHighlight(-1);
    } else {
      this.scrollToHighlight(0, this.totalSelections - 1);
    }
  }

  onNextHighlight(): void {
    if (this.getHighlightedElement(this.currentSelection + 1)) {
      this.scrollToHighlight(1);
    } else {
      this.scrollToHighlight(0, 0);
    }
  }

  scrollToHighlight(scrollBy: number, newSelection?: number) {
    this.removeHighlightSelection();
    if (newSelection !== undefined) {
      this.currentSelection = newSelection;
    }

    if (this.getHighlightedElement(this.currentSelection + scrollBy)) {
      this.currentSelection += scrollBy;
      const nextSelection: HTMLElement = this.getHighlightedElement(
        this.currentSelection
      );
      const {
        clientHeight,
        offsetTop,
        scrollTop,
      } = this.scrollTarget.nativeElement;
      const top = nextSelection.offsetTop - offsetTop;

      if (top > clientHeight + scrollTop || top < scrollTop) {
        nextSelection.scrollIntoView(true);
      }
      nextSelection.className = 'highlight highlight-selected';
    }
  }

  removeHighlightSelection(): HTMLElement {
    const element: HTMLElement = this.getHighlightedElement(
      this.currentSelection
    );
    if (element) {
      element.className = 'highlight';
    }
    return element;
  }

  getHighlightedElement(index: number): HTMLElement {
    return document.getElementsByClassName('highlight')[index] as HTMLElement;
  }

  updateSelectedCount() {
    let count = 0;
    if (this.filterText.length > 0) {
      this.containerLogs.map(log => {
        count += (log.message.match(new RegExp(this.filterText, 'g')) || [])
          .length;
        if (this.shouldDisplayTimestamp) {
          count += (
            formatDate(log.timestamp, this.timeFormat, 'en-US').match(
              new RegExp(this.filterText, 'g')
            ) || []
          ).length;
        }
      });
    }
    this.totalSelections = count;
  }
}
