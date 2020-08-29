// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ModalService {
  private modalOpened = new BehaviorSubject<boolean>(false);
  isOpened = this.modalOpened.asObservable();

  constructor() {}

  openModal() {
    this.modalOpened.next(true);
  }

  closeModal() {
    this.modalOpened.next(false);
  }

  setState(opened: boolean) {
    this.modalOpened.next(opened);
  }
}
