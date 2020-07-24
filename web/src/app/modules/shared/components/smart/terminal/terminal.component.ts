// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  AfterViewInit,
  Component,
  ElementRef,
  HostListener,
  Input,
  OnDestroy,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import {
  TerminalOutputService,
  TerminalOutputStreamer,
} from 'src/app/modules/shared/terminals/terminals.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { TerminalView } from 'src/app/modules/shared/models/content';
import { WebsocketService } from '../../../services/websocket/websocket.service';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-terminal',
  styleUrls: ['./terminal.component.scss'],
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements OnDestroy, AfterViewInit {
  constructor(
    private terminalService: TerminalOutputService,
    private wss: WebsocketService
  ) {}

  selectedContainer = '';
  private terminalStream: TerminalOutputStreamer;
  private term: Terminal;
  private fitAddon: FitAddon;
  trackByIdentity = trackByIdentity;
  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) terminalDiv: ElementRef;

  ngOnDestroy(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream.exitMessage.unsubscribe();
      this.terminalStream = null;
    }
  }

  ngAfterViewInit() {
    if (this.view) {
      const logLevel = 'info';
      const { terminal } = this.view.config;
      const { active } = terminal;
      const disableStdin = !active;
      this.term = new Terminal({
        logLevel,
        disableStdin,
      });
      setTimeout(() => {
        this.initStream();
      });
      this.enableResize();
      this.term.onData(data => {
        if (active) {
          this.wss.sendMessage('action.octant.dev/sendTerminalCommand', {
            key: data,
          });
        }
      });
      this.fitAddon = new FitAddon();
      this.term.loadAddon(this.fitAddon);
      this.term.open(this.terminalDiv.nativeElement);
      this.term.focus();
      this.fitAddon.fit();
    }
  }

  enableResize() {
    let timeOut = null;
    const resizeDebounce = (e: { cols: number; rows: number }) => {
      const resize = () => {
        this.wss.sendMessage('action.octant.dev/sendTerminalResize', {
          rows: e.rows,
          cols: e.cols,
        });
        this.fitAddon.fit();
      };

      if (timeOut != null) {
        clearTimeout(timeOut);
      }
      timeOut = setTimeout(resize, 500);
    };
    this.term.onResize(resizeDebounce);
  }

  onContainerChange(containerSelection: string): void {
    this.terminalService.selectedContainer = containerSelection;
    this.selectedContainer = containerSelection;
    this.term.reset();
    this.initStream();
    this.term.focus();
    this.fitAddon.fit();
  }

  initStream() {
    const { namespace, podName, terminal } = this.view.config;
    const { container } = terminal;
    if (this.terminalService.selectedContainer) {
      this.selectedContainer = this.terminalService.selectedContainer;
    }
    if (namespace && podName && container) {
      if (
        this.terminalService.namespace === namespace &&
        this.terminalService.podName === podName &&
        this.selectedContainer
      ) {
        this.selectedContainer = this.terminalService.selectedContainer;
        this.terminalStream = this.terminalService.createStream(
          namespace,
          podName,
          this.selectedContainer
        );
      } else {
        this.terminalStream = this.terminalService.createStream(
          namespace,
          podName,
          container
        );
      }
      this.terminalStream.exitMessage.subscribe((exitMessage: string) => {
        if (exitMessage && exitMessage.length !== 0) {
          this.selectedContainer = undefined;
          this.terminalService.selectedContainer = this.selectedContainer;
        }
      });
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        if (scrollback && scrollback.length !== 0) {
          this.term.write(atob(scrollback).replace(/\n/g, '\n\r'));
        }
      });
      this.terminalStream.line.subscribe((line: string) => {
        if (line && line.length !== 0) {
          this.term.write(atob(line).replace(/\n/g, '\n\r'));
        }
      });
      this.terminalService.namespace = namespace;
      this.terminalService.podName = podName;
    }
  }

  onResize() {
    this.fitAddon.fit();
  }
}
