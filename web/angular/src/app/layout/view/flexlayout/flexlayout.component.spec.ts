import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewModule } from '../view.module';
import { FlexlayoutComponent } from './flexlayout.component';

describe('FlexlayoutComponent', () => {
  let component: FlexlayoutComponent;
  let fixture: ComponentFixture<FlexlayoutComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ViewModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FlexlayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
