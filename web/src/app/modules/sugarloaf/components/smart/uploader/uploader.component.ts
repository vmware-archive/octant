// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { WebsocketService } from '../../../../shared/services/websocket/websocket.service';
import { Subscription } from 'rxjs';
import { KubeContextService } from '../../../../shared/services/kube-context/kube-context.service';
import { has } from 'lodash';

@Component({
  selector: 'app-uploader',
  templateUrl: './uploader.component.html',
  styleUrls: ['./uploader.component.scss'],
})
export class UploaderComponent implements OnInit, OnDestroy {
  private kubeContextSubscription: Subscription;
  inputValue: string;
  showModal: boolean;

  constructor(private websocketService: WebsocketService) {}

  ngOnInit(): void {
    this.websocketService.registerHandler('event.octant.dev/loading', () => {
      this.showModal = true;
    });
    this.websocketService.registerHandler('event.octant.dev/refresh', () => {
      setTimeout(window.location.reload.bind(window.location), 1000);
    });
  }

  ngOnDestroy(): void {}

  upload() {
    this.websocketService.sendMessage('action.octant.dev/uploadKubeConfig', {
      kubeConfig: window.btoa(this.inputValue),
    });
  }
}
