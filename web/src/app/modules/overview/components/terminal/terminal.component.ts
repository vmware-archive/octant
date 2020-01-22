// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ViewChild,
  OnInit,
  OnDestroy,
  AfterViewInit,
  Input,
  ElementRef,
  ViewEncapsulation,
  HostBinding,
  HostListener,
} from '@angular/core';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import {
  TerminalOutputStreamer,
  TerminalOutputService,
} from 'src/app/services/terminals/terminals.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { TerminalView, TerminalDetail } from 'src/app/models/content';
import { WebsocketService } from '../../services/websocket/websocket.service';
import { SliderService } from 'src/app/services/slider/slider.service';
import { Subscription } from 'rxjs';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-terminal',
  styleUrls: ['./terminal.component.scss'],
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements OnDestroy, AfterViewInit {
  private terminalStream: TerminalOutputStreamer;
  private term: Terminal;
  private fitAddon: FitAddon;
  trackByIdentity = trackByIdentity;
  sliderServiceSubscription: Subscription;

  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) terminalDiv: ElementRef;
  @HostListener('click') onClick() {
    this.term.focus();
    this.fitAddon.fit();
  }

  constructor(
    private terminalService: TerminalOutputService,
    private wss: WebsocketService,
    private sliderService: SliderService
  ) {}

  compareFn(c1: TerminalDetail, c2: TerminalDetail): boolean {
    return c1 && c2 ? c1.uuid === c2.uuid : c1 === c2;
  }

  ngOnDestroy(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }

    if (this.sliderServiceSubscription) {
      this.sliderServiceSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    if (this.view) {
      const disableStdin = this.view.config.terminal.active ? false : true;
      const logLevel = 'info';

      this.term = new Terminal({
        logLevel,
        disableStdin,
      });

      this.initStream();
      this.enableResize();
      this.term.onData(data => {
        if (this.view.config.terminal.active === true) {
          this.wss.sendMessage('sendTerminalCommand', {
            terminalID: this.view.config.terminal.uuid,
            key: data,
          });
        }
      });
      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.term.open(this.terminalDiv.nativeElement);
      this.term.focus();
      this.fitAddon.fit();

      this.sliderServiceSubscription = this.sliderService.resizedSliderEvent.subscribe(
        () => {
          this.term.focus();
          setTimeout(() => {
            this.fitAddon.fit();
          }, 0);
        }
      );
    }
  }

  enableResize() {
    let timeOut = null;
    const resizeDebounce = (e: { cols: number; rows: number }) => {
      const resize = () => {
        if (this.view.config.terminal.active) {
          this.wss.sendMessage('sendTerminalResize', {
            terminalID: this.view.config.terminal.uuid,
            rows: e.rows,
            cols: e.cols,
          });
        }
      };

      if (timeOut != null) {
        clearTimeout(timeOut);
      }
      timeOut = setTimeout(resize, 1000);
      this.fitAddon.fit();
    };
    this.term.onResize(resizeDebounce);
  }

  initStream() {
    const namespace = this.view.config.namespace;
    const name = this.view.config.name;
    const terminal = this.view.config.terminal;

    if (namespace && name && terminal.container && terminal.uuid) {
      this.terminalStream = this.terminalService.createStream(
        namespace,
        name,
        terminal.container,
        terminal.uuid
      );
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        if (scrollback && scrollback.length !== 0) {
          this.term.write(atob(scrollback).replace(/\n/g, '\n\r'));
        }
      });
      if (terminal.active) {
        this.terminalStream.line.subscribe((line: string) => {
          if (line && line.length !== 0) {
            this.term.write(atob(line).replace(/\n/g, '\n\r'));
          }
        });
      }
    }
  }
}
