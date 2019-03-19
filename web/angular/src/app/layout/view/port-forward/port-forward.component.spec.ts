import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PortForwardComponent } from './port-forward.component';

describe('PortForwardComponent', () => {
  let component: PortForwardComponent;
  let fixture: ComponentFixture<PortForwardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PortForwardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PortForwardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
