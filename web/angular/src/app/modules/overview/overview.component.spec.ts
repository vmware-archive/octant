import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from 'src/app/services/data.service';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';

import { OverviewModule } from './overview.module';
import { OverviewComponent } from './overview.component';

const dataServiceStub: Partial<DataService> = {
  stopPoller: () => {},
};
const activatedRouteStub = new ActivatedRouteStub();

describe('OverviewComponent', () => {
  let component: OverviewComponent;
  let fixture: ComponentFixture<OverviewComponent>;

  beforeEach(async(() => {
    const routerSpy = jasmine.createSpyObj('Router', ['navigateByUrl', 'navigate']);
    TestBed.configureTestingModule({
      imports: [OverviewModule],
      providers: [
        { provide: ActivatedRoute, useValue: activatedRouteStub },
        { provide: DataService, useValue: dataServiceStub },
        { provide: Router, useValue: routerSpy },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
