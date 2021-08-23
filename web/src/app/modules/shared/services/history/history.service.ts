import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { ContentService } from '../content/content.service';

export interface ContentRoute {
  title: string;
  path: string;
}

@Injectable({
  providedIn: 'root',
})
export class HistoryService {
  localStorage: Storage;
  history: BehaviorSubject<ContentRoute[]>;
  private previousRoutes: ContentRoute[];
  private HISTORY_SIZE = 10;
  private STORAGE_KEY = '__octant_history';

  constructor(private contentService: ContentService) {
    this.localStorage = window.localStorage;
    this.history = new BehaviorSubject<ContentRoute[]>([]);
    this.previousRoutes = this.getHistory();
    this.history.next(this.previousRoutes);

    this.contentService.title.subscribe(nsTitle => {
      if (
        nsTitle.title?.length > 0 &&
        nsTitle.path !== '' &&
        nsTitle.path !== this.previousRoutes[0]?.path
      ) {
        const historyTitle = nsTitle.namespace
          ? `${nsTitle.title} | ${nsTitle.namespace}`
          : nsTitle.title;
        const currentEntry = { title: historyTitle, path: nsTitle.path };
        let newHistory: ContentRoute[];
        let pr = this.getHistory();

        if (pr.find(r => r.path === currentEntry.path)) {
          const rest = pr.filter(r => r.path !== currentEntry.path);
          newHistory = [currentEntry, ...rest];
        } else {
          newHistory = [currentEntry, ...pr].slice(0, this.HISTORY_SIZE);
        }

        this.saveHistory(newHistory);
        this.history.next(newHistory);
      }
    });
  }

  getHistory() {
    if (this.isLocalStorageSupported) {
      return JSON.parse(this.localStorage.getItem(this.STORAGE_KEY)) || [];
    }

    return this.previousRoutes || [];
  }

  saveHistory(history: ContentRoute[]) {
    if (this.isLocalStorageSupported) {
      this.localStorage.setItem(this.STORAGE_KEY, JSON.stringify(history));
    } else {
      this.previousRoutes = history;
    }
  }

  get isLocalStorageSupported(): boolean {
    return !!this.localStorage;
  }
}
