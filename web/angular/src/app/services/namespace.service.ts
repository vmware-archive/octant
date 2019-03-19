import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';

import { DataService } from '../data.service';
import { Namespaces } from '../models/namespace';

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
  }

  setNamespace(namespace: string) {
    this.current.next(namespace);

    const currentURL = this.router.url.split('/');
    if (currentURL.length > 3) {
      if (currentURL[3] === 'namespace') {
        currentURL[4] = namespace;

        if (currentURL.length > 6) {
          const newURL = currentURL.slice(0, 6);
          this.router.navigate(newURL);
        } else {
          this.router.navigate(currentURL);
        }
      } else {
        this.router.navigate(currentURL);
      }
    }
  }
}
