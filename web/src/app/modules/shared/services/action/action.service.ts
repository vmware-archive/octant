import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';

@Injectable({
  providedIn: 'root',
})
export class ActionService {
  constructor(private websocketService: WebsocketService) {}

  perform(update: any) {
    this.websocketService.sendMessage('performAction', update);
  }
}
