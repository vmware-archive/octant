import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SafePipe } from '../../pipes/safe.pipe';
import { IFrameComponent } from './iframe.component';

describe('IFrameComponent', () => {
  let component: IFrameComponent;
  let fixture: ComponentFixture<IFrameComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [IFrameComponent, SafePipe],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IFrameComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
