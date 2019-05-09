import { EventSourceStub } from './event-source.service';

describe('EventSourceStub', () => {
  it('should act like an EventSource', () => {
    const eventSource = new EventSourceStub();

    eventSource.addEventListener('content', (message: MessageEvent) => {
      const { data } = message;
      expect(data.a).toBe(1);
      expect(data.b).toBe(2);
    });

    eventSource.addEventListener('content', (message: MessageEvent) => {
      const { data } = message;
      expect(data.a).toBe(1);
      expect(data.b).toBe(2);
    });

    const errorSpy = jasmine.createSpy();
    eventSource.addEventListener('error', errorSpy);

    eventSource.queueMessage('content', { a: 1, b: 2 });
    eventSource.queueMessage('error');
    eventSource.flush();

    expect(errorSpy.calls.count()).toBe(1);
  });
});
