import { TestBed } from '@angular/core/testing';

import { ActionService } from './action.service';
import { WebsocketService } from '../websocket/websocket.service';
import { WebsocketServiceMock } from '../websocket/mock';

describe('ActionService', () => {
  let service: ActionService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ActionService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
      ],
    });

    service = TestBed.inject(ActionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('performAction', () => {
    let websocketService: WebsocketService;

    beforeEach(() => {
      websocketService = TestBed.inject(WebsocketService);
      spyOn(websocketService, 'sendMessage');
    });

    it('sends a performAction message to the server', () => {
      const update = { foo: 'bar' };
      service.perform(update);
      expect(websocketService.sendMessage).toHaveBeenCalledWith(
        'action.octant.dev/performAction',
        update
      );
    });
  });
});
