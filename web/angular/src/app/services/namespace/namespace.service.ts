import { Injectable } from '@angular/core';
import { NavigationEnd, PRIMARY_OUTLET, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';

import { DataService } from '../data/data.service';
import { Namespaces } from '../../models/namespace';

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  current = new BehaviorSubject<string>('default');
  list = new BehaviorSubject<string[]>([]);

  constructor(private router: Router, private dataService: DataService) {
    this.dataService.getNamespaces().subscribe((namespaces: Namespaces) => {
      this.list.next(namespaces.namespaces);
    });

    this.dataService.pollNamespaces().subscribe((namespaces: string[]) => {
      this.list.next(namespaces);
    });

    this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        const namespace = this.getNamespaceFromUrl(this.router.url);
        const currentNS = this.current.getValue();
        if (currentNS !== namespace) {
          this.current.next(namespace);
        }
      }
    });
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
