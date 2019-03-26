import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { ReplaySubject } from 'rxjs';

import { DataService } from '../data/data.service';
import { OverviewModule } from '../../modules/overview/overview.module';
import { NamespaceService } from './namespace.service';

const dataServiceSpy = jasmine.createSpyObj('DataService', ['getNamespaces', 'pollNamespaces']);

const routerSpy = jasmine.createSpyObj('Router', ['navigate']);

describe('NamespaceService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [OverviewModule],
      providers: [{ provide: Router, useValue: routerSpy }, { provide: DataService, useValue: dataServiceSpy }],
    })
  );

  it('should be created', () => {
    const namespaces = new ReplaySubject<string[]>();

    dataServiceSpy.getNamespaces.and.returnValue(namespaces);
    dataServiceSpy.pollNamespaces.and.returnValue(namespaces);

    const service: NamespaceService = TestBed.get(NamespaceService);
    expect(service).toBeTruthy();
  });
});
