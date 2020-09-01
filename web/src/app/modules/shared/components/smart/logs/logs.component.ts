// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterContentChecked,
  AfterViewChecked,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  IterableDiffer,
  IterableDiffers,
  OnDestroy,
  OnInit,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { LogEntry, LogsView } from 'src/app/modules/shared/models/content';
import {
  PodLogsService,
  PodLogsStreamer,
} from 'src/app/modules/shared/pod-logs/pod-logs.service';
import { formatDate } from '@angular/common';
import { Subscription } from 'rxjs';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.scss'],
  encapsulation: ViewEncapsulation.None,
  changeDetection: ChangeDetectionStrategy.Default,
})
export class LogsComponent
  extends AbstractViewComponent<LogsView>
  implements OnInit, OnDestroy, AfterContentChecked, AfterViewChecked {
  private logStream: PodLogsStreamer;
  scrollToBottom = true;

  private containerLogsDiffer: IterableDiffer<LogEntry>;
  @ViewChild('scrollTarget', { static: true }) scrollTarget: ElementRef;
  containerLogs: LogEntry[] = [];

  selectedContainer = '';
  shouldDisplayTimestamp = false;
  shouldDisplayName = true;
  showOnlyFiltered = false;
  filterText = '';
  oldFilterText = '';
  currentSelection = 0;
  totalSelections = 0;
  timeFormat = 'MMM d, y h:mm:ss a z';
  regexFlags = 'gi';

  private logSubscription: Subscription;

  constructor(
    private podLogsService: PodLogsService,
    private iterableDiffers: IterableDiffers,
    private cdr: ChangeDetectorRef
  ) {
    super();
  }

  ngOnInit() {
    this.containerLogsDiffer = this.iterableDiffers
      .find(this.containerLogs)
      .create();
    this.startStream();
  }

  protected update() {
    if (this.v.config.containers && this.v.config.containers.length > 0) {
      this.selectedContainer = this.v.config.containers[0];
    }
  }

  onContainerChange(containerSelection: string): void {
    this.selectedContainer = containerSelection;
    if (this.logStream) {
      this.containerLogs = [];
      this.logStream.close();
      this.logStream = null;
    }
    if (this.selectedContainer === '') {
      this.shouldDisplayName = true;
    } else {
      this.shouldDisplayName = false;
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
    if (namespace && pod) {
      this.logStream = this.podLogsService.createStream(
        namespace,
        pod,
        container
      );
      this.logSubscription = this.logStream.logEntry.subscribe(
        (entry: LogEntry) => {
          if (entry.message == null) {
            return;
          }
          this.containerLogs.push(entry);
          this.updateSelectedCount();
          this.cdr.markForCheck();
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
    if (!this.filterText) {
      return text;
    }

    const matched = new RegExp(this.filterText, this.regexFlags).exec(text);
    if (matched === null) {
      return text;
    }

    const filter =
      matched[0] && matched[0].length > 0
        ? this.filterText
        : this.filterText + '.*$';

    return text.replace(new RegExp(filter, this.regexFlags), match => {
      return '<span class="highlight">' + match + '</span>';
    });
  }

  public filterFunction(logs: LogEntry[]): LogEntry[] {
    if (this.showOnlyFiltered) {
      return logs.filter(log => {
        const hasFiltered = this.matchRegex(log);
        return hasFiltered && hasFiltered.length > 0;
      });
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

  matchRegex(input: LogEntry) {
    let match = input.message.match(
      new RegExp(this.filterText, this.regexFlags)
    );
    if (match) {
      return match;
    }

    if (this.shouldDisplayTimestamp && input.timestamp) {
      const timestamp = formatDate(input.timestamp, this.timeFormat, 'en-US');
      if (timestamp && timestamp.length > 0) {
        match = timestamp.match(new RegExp(this.filterText, this.regexFlags));
        if (match) {
          return match;
        }
      }
    }

    if (this.shouldDisplayName) {
      match = input.container.match(
        new RegExp(this.filterText, this.regexFlags)
      );
      return match || [];
    }
    return [];
  }

  updateSelectedCount() {
    let count = 0;
    if (this.filterText.length > 0) {
      this.containerLogs.map(log => {
        count += (this.matchRegex(log) || []).length;
      });
    }
    this.totalSelections = count;
  }
}
