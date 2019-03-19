import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from 'src/app/data.service';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';

import { ViewModule } from '../view/view.module';
import { ContentComponent } from './content.component';

const dataServiceStub: Partial<DataService> = {
  stopPoller: () => {},
};
const activatedRouteStub = new ActivatedRouteStub();

describe('ContentComponent', () => {
  let component: ContentComponent;
  let fixture: ComponentFixture<ContentComponent>;

  beforeEach(async(() => {
    const routerSpy = jasmine.createSpyObj('Router', ['navigateByUrl', 'navigate']);
    TestBed.configureTestingModule({
      imports: [ViewModule],
      providers: [
        { provide: ActivatedRoute, useValue: activatedRouteStub },
        { provide: DataService, useValue: dataServiceStub },
        { provide: Router, useValue: routerSpy },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
