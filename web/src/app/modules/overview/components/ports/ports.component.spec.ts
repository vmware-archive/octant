import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PortsComponent } from './ports.component';

describe('PortsComponent', () => {
  let component: PortsComponent;
  let fixture: ComponentFixture<PortsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PortsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PortsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
