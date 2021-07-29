import { ContentResponse } from '../../models/content';
import { BehaviorSubject } from 'rxjs';
import { ContentRoute } from './history.service';

export class HistoryServiceMock {
  history = new BehaviorSubject<ContentRoute[]>([]);

  pushHistory(history: ContentRoute[]) {
    this.history.next(history);
  }
}
