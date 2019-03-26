import { TestBed } from '@angular/core/testing';

import { LabelFilterService } from './label-filter.service';

describe('LabelFilterService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: LabelFilterService = TestBed.get(LabelFilterService);
    expect(service).toBeTruthy();
  });
});
