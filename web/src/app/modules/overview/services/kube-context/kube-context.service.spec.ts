import { TestBed } from '@angular/core/testing';

import { KubeContextService } from './kube-context.service';

describe('KubeContextService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: KubeContextService = TestBed.get(KubeContextService);
    expect(service).toBeTruthy();
  });
});
