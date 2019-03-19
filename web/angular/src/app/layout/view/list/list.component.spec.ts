import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewModule } from '../view.module';
import { ViewListComponent } from './list.component';

describe('ViewListComponent', () => {
  let component: ViewListComponent;
  let fixture: ComponentFixture<ViewListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ViewModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ViewListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
