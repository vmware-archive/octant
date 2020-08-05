import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { StepperComponent } from './stepper.component';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { StepperView } from '../../../models/content';
import {
  BrowserAnimationsModule,
  NoopAnimationsModule,
} from '@angular/platform-browser/animations';

describe('StepperComponent', () => {
  let component: StepperComponent;
  let fixture: ComponentFixture<StepperComponent>;
  const formBuilder: FormBuilder = new FormBuilder();

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        ReactiveFormsModule,
        BrowserAnimationsModule,
        NoopAnimationsModule,
      ],
      providers: [{ provide: FormBuilder, useValue: formBuilder }],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StepperComponent);
    component = fixture.componentInstance;

    const view: StepperView = {
      metadata: {
        type: 'stepper',
      },
      config: {
        action: 'action.octant.dev/test',
        steps: [
          {
            name: 'step name',
            form: { fields: [] },
            title: 'step title',
            description: 'step description',
          },
        ],
      },
    };
    component.view = view;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
