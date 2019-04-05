import { TestBed } from '@angular/core/testing';

import { PortForwardService } from './port-forward.service';

describe('PortForwardService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: PortForwardService = TestBed.get(PortForwardService);
    expect(service).toBeTruthy();
  });
});
