// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
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
import { IconService } from './modules/overview/services/icon.service';
import { SliderService } from './services/slider/slider.service';
import { Router, NavigationStart } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  @ViewChild('scrollTarget', { static: false }) scrollTarget: ElementRef;
  navigation: Navigation;
  previousUrl: string;
  style: object = {};

  constructor(
    private websocketService: WebsocketService,
    private iconService: IconService,
    private sliderService: SliderService,
    private router: Router
  ) {
    iconService.load({
      iconName: 'octant-logo',
      // tslint:disable-next-line:max-line-length
      iconSource: `<svg width="106px" height="126px" viewBox="0 0 106 126" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"> <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"> <g id="288026-vmw-os-lgo-octant-final" fill-rule="nonzero"> <g id="Group"> <circle id="Oval" fill="#FFFFFF" cx="53.4" cy="67.1" r="52.5"></circle> <path d="M58.6,112.7 L58.6,112.7 C58.6,115.5 56.4,117.7 53.6,117.7 L53.6,117.7 C50.8,117.7 48.6,115.5 48.6,112.7 L48.6,112.7 C48.6,110 50.8,107.7 53.6,107.7 L53.6,107.7 C56.4,107.7 58.6,110 58.6,112.7 Z" id="Path" fill="#62BB46"></path> <path d="M53.4,14.6 C24.4,14.6 0.9,38.1 0.9,67.1 C0.9,81 6.3,93.7 15.2,103 C16.6,104.5 18.1,105.8 19.6,107.1 L23.2,97.3 C21.5,95.6 19.9,93.7 18.5,91.7 C13.6,84.7 10.6,76.2 10.6,67 C10.6,43.3 29.8,24.2 53.4,24.2 C77,24.2 96.2,43.4 96.2,67 C96.2,76 93.4,84.3 88.7,91.2 C87.2,93.3 85.6,95.3 83.8,97.1 L87.4,106.9 C89,105.5 90.6,104.1 92,102.5 C100.6,93.2 105.8,80.7 105.8,67 C105.9,38.1 82.4,14.6 53.4,14.6 Z" id="Path" fill="#727175"></path> <path d="M83.8,97.2 L60.5,33.5 L57.2,34.7 L56.6,31.9 L55.4,32.1 L70.2,106.4 C65.9,108.3 61.2,109.4 56.2,109.8 C55.2,109.9 54.3,109.9 53.3,109.9 C52.5,109.9 51.7,109.9 50.9,109.8 C45.8,109.5 41,108.3 36.5,106.4 L51.3,32.1 L50.1,31.9 L49.5,34.7 L46.3,33.5 L22.9,97.3 L19.3,107.1 C19.7,107.5 20.2,107.9 20.6,108.2 L20.6,108.1 L20.6,108.1 L20.6,108.2 C20.6,108.2 20.6,108.2 20.6,108.2 L20.6,108.2 C21.3,108.8 22.1,109.3 22.9,109.9 C24.1,110.7 25.3,111.5 26.6,112.3 L26.6,112.3 C34.4,116.9 43.5,119.6 53.2,119.6 C63.1,119.6 72.4,116.8 80.3,112.1 L80.3,112.1 C81.5,111.4 82.6,110.6 83.7,109.8 C84.6,109.2 85.4,108.6 86.2,107.9 C86.6,107.6 86.9,107.3 87.3,107 L83.8,97.2 Z M29.9,102.9 C28.8,102.1 27.7,101.3 26.6,100.5 L48.6,40.6 L35.5,106 C33.5,105.1 31.7,104.1 29.9,102.9 L29.9,102.9 Z M77.3,102.6 L77.3,102.6 C75.4,103.9 73.5,105 71.4,105.9 L58.3,40 L80.4,100.4 C79.4,101.1 78.4,101.9 77.3,102.6 Z" id="Shape" fill="#00438C"></path> <g transform="translate(43.000000, 8.000000)" id="Path"> <polygon fill="#FFFFFF" points="10.6 0.7 2.3 4.9 0.1 13.9 5.9 21.2 15.1 21.2 21 13.8 18.9 4.9"></polygon> <path d="M15.5,11.3 L15.5,11.3 C15.5,14.1 13.3,16.3 10.5,16.3 L10.5,16.3 C7.7,16.3 5.5,14.1 5.5,11.3 L5.5,11.3 C5.5,8.6 7.7,6.3 10.5,6.3 L10.5,6.3 C13.3,6.3 15.5,8.5 15.5,11.3 Z" fill="#62BB46"></path> </g> <path d="M71.6,23.5 L68,8.1 L53.6,0.9 L39.3,8 L35.5,23.5 L45.6,36.2 L49.9,36.2 L49.9,105.4 L56.9,105.4 L56.9,36.2 L61.5,36.2 L71.6,23.5 Z M48.9,29.2 L43.1,21.9 L45.3,12.9 L53.6,8.8 L61.9,12.9 L64,21.9 L58.1,29.2 L48.9,29.2 Z" id="Shape" fill="#009CDC"></path> <g transform="translate(42.000000, 103.000000)"> <path d="M10.6,6.2 C7.9,6.6 6,9.1 6.4,11.9 C6.8,14.6 9.3,16.5 12.1,16.1 C14.8,15.7 16.7,13.2 16.3,10.4 C15.9,7.6 13.3,5.8 10.6,6.2 Z" id="Path" fill="#FFFFFF"></path> <path d="M22.2,9.5 C21.3,3.5 15.7,-0.7 9.7,0.2 L9.7,0.2 C3.7,1.1 -0.5,6.7 0.4,12.7 C1.2,18.2 5.9,22.1 11.3,22.1 C11.8,22.1 12.4,22.1 12.9,22 C18.9,21.1 23.1,15.5 22.2,9.5 Z M12.1,16 C9.4,16.4 6.8,14.5 6.4,11.8 C6,9.1 7.9,6.5 10.6,6.1 C13.3,5.7 15.9,7.6 16.3,10.3 C16.7,13.1 14.8,15.6 12.1,16 Z" id="Shape" fill="#62BB46"></path> </g> </g> </g> </g> </svg>`,
    });
    this.sliderService.setHeight$.subscribe((data: number) => {
      Object.assign(this.style, { marginBottom: `${data}px` });
      setTimeout(() => {
        this.sliderService.resizedSliderEvent.emit(true);
      }, 0);
    });
    router.events.subscribe(data => {
      if (data instanceof NavigationStart) {
        this.style = {};
      }
    });
  }

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
