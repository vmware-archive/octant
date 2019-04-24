import { Injectable } from '@angular/core';
import { NavigationEnd, PRIMARY_OUTLET, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import _ from 'lodash';
import { DataService } from '../data/data.service';
import { Namespaces } from '../../models/namespace';
import { NotifierService, NotifierServiceSession, NotifierSignalType } from '../notifier/notifier.service';

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  notifierServiceSession: NotifierServiceSession;
  current = new BehaviorSubject<string>('default');
  list = new BehaviorSubject<string[]>([]);

  constructor(private router: Router, private dataService: DataService, notifierService: NotifierService) {
    this.notifierServiceSession = notifierService.createSession();

    this.dataService.getNamespaces().subscribe((namespaces: Namespaces) => {
      this.list.next(namespaces.namespaces);
    });

    this.dataService.pollNamespaces().subscribe((namespaces: string[]) => {
      this.list.next(namespaces);
    });

    this.router.events.subscribe((event) => {
      if (!(event instanceof NavigationEnd)) {
        return;
      }
      this.notifierServiceSession.removeAllSignals();
      const namespace = this.getNamespaceFromUrl(this.router.url);

      if (namespace) {
        const listOfNamespaces = this.list.getValue();
        if (listOfNamespaces.length > 0) {
          if (!this.isNamespaceValid(namespace)) {
            this.notifierServiceSession.pushSignal(NotifierSignalType.ERROR, 'The current set namespace is not valid');
            return;
          }
        }
      }

      const currentNS = this.current.getValue();
      if (currentNS !== namespace) {
        this.current.next(namespace);
      }
    });
  }

  isNamespaceValid(namespaceToCheck: string) {
    const listOfNamespaces = this.list.getValue();
    return _.includes(listOfNamespaces, namespaceToCheck);
  }

  getNamespaceFromUrl(url: string): string {
    if (!url) {
      throw new Error('No url');
    }
    const urlTree = this.router.parseUrl(url);
    const urlSegments = urlTree.root.children[PRIMARY_OUTLET].segments;
    if (urlSegments.length > 3 && urlSegments[2].path === 'namespace') {
      return urlSegments[3].path;
    }
  }

  setNamespace(namespace: string) {
    this.current.next(namespace);
    this.router.navigate(['/content', 'overview', 'namespace', namespace]);
  }
}
