import { Component, OnInit } from '@angular/core';
import { ClarityIcons, bellIcon } from '@cds/core/icon';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';

interface InternalError {
  name: string;
  error: string;
}

@Component({
  selector: 'app-notifications',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.scss'],
})
export class NotificationsComponent implements OnInit {
  hidden = true;
  notifications: InternalError[] = [];

  trackByFn = trackByIndex;

  constructor(private websocketService: WebsocketService) {
    ClarityIcons.addIcons(bellIcon);
  }

  ngOnInit(): void {
    this.websocketService.registerHandler(
      'event.octant.dev/notification',
      (data: InternalError) => {
        this.notifications.push(data);
      }
    );
    this.websocketService.sendMessage('event.octant.dev/notification', {});
  }

  toggleModal(): void {
    this.hidden = !this.hidden;
  }
}
