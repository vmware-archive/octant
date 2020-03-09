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
import { v4 as uuidv4 } from 'uuid';
import {
  TerminalOutputService,
  TerminalOutputStreamer,
} from 'src/app/modules/shared/terminals/terminals.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { TerminalView } from 'src/app/modules/shared/models/content';
import { WebsocketService } from '../../../services/websocket/websocket.service';
import { ActionService } from '../../../services/action/action.service';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-terminal',
  styleUrls: ['./terminal.component.scss'],
  templateUrl: './terminal.component.html',
})
export class TerminalComponent implements OnDestroy, AfterViewInit {
  constructor(
    private terminalService: TerminalOutputService,
    private wss: WebsocketService,
    private actionService: ActionService
  ) {}

  private static sessionID = uuidv4();
  private terminalStream: TerminalOutputStreamer;
  private term: Terminal;
  private fitAddon: FitAddon;
  trackByIdentity = trackByIdentity;
  @Input() view: TerminalView;
  @ViewChild('terminal', { static: true }) terminalDiv: ElementRef;
  @HostListener('click') onClick() {
    this.term.focus();
    this.fitAddon.fit();
  }

  ngOnDestroy(): void {
    if (this.terminalStream) {
      this.terminalStream.scrollback.unsubscribe();
      this.terminalStream.line.unsubscribe();
      this.terminalStream = null;
    }
  }

  ngAfterViewInit() {
    if (this.view) {
      const logLevel = 'info';
      const { podName, namespace, terminal } = this.view.config;
      const { active, uuid, command, container } = terminal;
      const disableStdin = !active;
      const terminalCommand = {
        containerName: container,
        containerCommand: command,
        name: podName,
        namespace,
        kind: 'Pod',
        apiVersion: 'v1',
        sessionID: TerminalComponent.sessionID,
        action: 'overview/commandExec',
      };
      this.actionService.perform(terminalCommand);
      this.term = new Terminal({
        logLevel,
        disableStdin,
      });
      this.initStream();
      this.enableResize();
      this.term.onData(data => {
        if (active) {
          this.wss.sendMessage('sendTerminalCommand', {
            terminalID: uuid,
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
        const { active, uuid } = this.view.config.terminal;
        if (active) {
          this.wss.sendMessage('sendTerminalResize', {
            terminalID: uuid,
            rows: e.rows,
            cols: e.cols,
          });
          this.fitAddon.fit();
        }
      };

      if (timeOut != null) {
        clearTimeout(timeOut);
      }
      timeOut = setTimeout(resize, 500);
    };
    this.term.onResize(resizeDebounce);
  }

  initStream() {
    const { namespace, podName, terminal } = this.view.config;
    const { active, uuid, container } = terminal;
    if (namespace && podName && container && uuid) {
      this.terminalStream = this.terminalService.createStream(
        namespace,
        podName,
        container,
        uuid
      );
      this.terminalStream.scrollback.subscribe((scrollback: string) => {
        if (scrollback && scrollback.length !== 0) {
          this.term.write(atob(scrollback).replace(/\n/g, '\n\r'));
        }
      });
      if (active) {
        this.terminalStream.line.subscribe((line: string) => {
          if (line && line.length !== 0) {
            this.term.write(atob(line).replace(/\n/g, '\n\r'));
          }
        });
      }
    }
  }
}
