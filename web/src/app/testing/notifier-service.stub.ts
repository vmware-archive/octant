export function notifierServiceStubFactory() {
  return {
    notifierSessionStub: jasmine.createSpyObj(['removeAllSignals', 'pushSignal']),
    createSession() {
      return this.notifierSessionStub;
    }
  };
}

export const notifierServiceStub = notifierServiceStubFactory();
