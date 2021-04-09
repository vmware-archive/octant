import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormViewContainerComponent } from './form-view-container.component';

describe('FormViewContainerComponent', () => {
  let component: FormViewContainerComponent;
  let fixture: ComponentFixture<FormViewContainerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [FormViewContainerComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FormViewContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
