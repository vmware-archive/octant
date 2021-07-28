import {
  async,
  ComponentFixture,
  TestBed,
  waitForAsync,
} from '@angular/core/testing';
import { StepperComponent } from './stepper.component';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { StepperView } from '../../../models/content';
import {
  BrowserAnimationsModule,
  NoopAnimationsModule,
} from '@angular/platform-browser/animations';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { anything, deepEqual, instance, mock, verify } from 'ts-mockito';
import { ActionService } from '../../../services/action/action.service';

describe('StepperComponent', () => {
  let component: StepperComponent;
  let fixture: ComponentFixture<StepperComponent>;
  const formBuilder: FormBuilder = new FormBuilder();

  const mockActionService: ActionService = mock(ActionService);

  const action = 'action.octant.dev/test';
  const view: StepperView = {
    metadata: {
      type: 'stepper',
    },
    config: {
      action,
      steps: [
        {
          name: 'step 1',
          form: { fields: [] },
          title: 'step title',
          description: 'step description',
        },
        {
          name: 'confirmation step',
          form: { fields: [] },
          title: 'step title',
          description: 'confirmation description',
        },
      ],
    },
  };

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        imports: [
          ReactiveFormsModule,
          BrowserAnimationsModule,
          NoopAnimationsModule,
        ],
        providers: [
          { provide: FormBuilder, useValue: formBuilder },
          {
            provide: ActionService,
            useValue: instance(mockActionService),
          },
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(StepperComponent);
    component = fixture.componentInstance;

    component.view = view;

    fixture.detectChanges();
  });

  it(
    'should submit form after completing each step',
    waitForAsync(() => {
      fixture.whenStable().then(() => {
        let nextButton =
          fixture.debugElement.nativeElement.querySelector('.next');
        nextButton.click();
        fixture.detectChanges();

        nextButton =
          fixture.debugElement.nativeElement.querySelector('.submit');
        nextButton.click();
        fixture.detectChanges();

        verify(
          mockActionService.perform(
            deepEqual({ action, 'step 1': {}, 'confirmation step': {} })
          )
        ).once();
      });
    })
  );
});
