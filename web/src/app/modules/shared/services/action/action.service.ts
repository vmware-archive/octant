import { Injectable } from '@angular/core';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';

@Injectable({
  providedIn: 'root',
})
export class ActionService {
  constructor(private websocketService: WebsocketService) {}

  perform(update: any) {
    this.websocketService.sendMessage(
      'action.octant.dev/performAction',
      update
    );
  }
}
