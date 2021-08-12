import { NamespacedTitle } from '../../models/content';
import { BehaviorSubject } from 'rxjs';

const emptyNsTitle: NamespacedTitle = {
  namespace: '',
  title: '',
  path: '',
};

export class ContentServiceMock {
  title = new BehaviorSubject<NamespacedTitle>(emptyNsTitle);

  pushTitle(title: string, path: string, namespace: string) {
    this.title.next({ namespace, title, path });
  }
}
