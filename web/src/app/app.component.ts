// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { Navigation } from './models/navigation';
import { WebsocketService } from './modules/overview/services/websocket/websocket.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  @ViewChild('scrollTarget') scrollTarget: ElementRef;
  navigation: Navigation;
  previousUrl: string;

  constructor(private websocketService: WebsocketService) {}

  ngOnInit(): void {
    this.websocketService.open();
  }

  closeSocket() {
    this.websocketService.close();
  }

  ngOnDestroy(): void {
    this.closeSocket();
  }
}
