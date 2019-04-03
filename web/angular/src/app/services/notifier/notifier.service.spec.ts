import { TestBed } from '@angular/core/testing';

import { NotifierService } from './notifier.service';

describe('NotifierService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: NotifierService = TestBed.get(NotifierService);
    expect(service).toBeTruthy();
  });
});
