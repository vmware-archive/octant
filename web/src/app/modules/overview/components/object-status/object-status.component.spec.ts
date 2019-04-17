import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { OverviewModule } from '../../overview.module';
import { ObjectStatusComponent } from './object-status.component';
import { ContentSwitcherComponent } from '../content-switcher/content-switcher.component';

describe('ObjectStatusComponent', () => {
  let component: ObjectStatusComponent;
  let fixture: ComponentFixture<ObjectStatusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        OverviewModule,
      ],
      declarations: [
        // ObjectStatusComponent,
        // ContentSwitcherComponent,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ObjectStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
