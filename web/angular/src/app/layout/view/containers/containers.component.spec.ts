import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContainersComponent } from './containers.component';

describe('ContainersComponent', () => {
  let component: ContainersComponent;
  let fixture: ComponentFixture<ContainersComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ContainersComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContainersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
