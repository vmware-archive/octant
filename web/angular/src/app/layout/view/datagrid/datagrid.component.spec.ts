import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewModule } from '../view.module';
import { DatagridComponent } from './datagrid.component';

describe('DatagridComponent', () => {
  let component: DatagridComponent;
  let fixture: ComponentFixture<DatagridComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ViewModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DatagridComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
