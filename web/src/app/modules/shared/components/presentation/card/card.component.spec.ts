import {
  ComponentFixture,
  fakeAsync,
  TestBed,
  waitForAsync,
} from '@angular/core/testing';

import { CardComponent } from './card.component';
import { Action, CardView, TextView } from '../../../models/content';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ViewService } from '../../../services/view/view.service';
import { viewServiceStub } from 'src/app/testing/view-service.stub';
import { SharedModule } from '../../../shared.module';
import { FormComponent } from '../form/form.component';
import { windowProvider, WindowToken } from '../../../../../window';
import { EditorComponent } from '../../smart/editor/editor.component';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';

describe('CardComponent', () => {
  let component: CardComponent;
  let fixture: ComponentFixture<CardComponent>;
  let formComponent: FormComponent;
  let formFixture: ComponentFixture<FormComponent>;
  const formBuilder: FormBuilder = new FormBuilder();

  const action: Action = {
    name: 'actionName',
    title: 'actionTitle',
    form: {
      fields: [],
    },
  };

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [EditorComponent, OverlayScrollbarsComponent],
        imports: [SharedModule, OverlayscrollbarsModule],
        providers: [
          { provide: FormBuilder, useValue: formBuilder },
          { provide: ViewService, useValue: viewServiceStub },
          { provide: WindowToken, useFactory: windowProvider },
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(CardComponent);
    component = fixture.componentInstance;

    const textView: TextView = {
      metadata: {
        type: 'text',
        accessor: '',
        title: [],
      },
      config: {
        value: 'text',
      },
    };

    component.view = {
      config: {
        actions: [],
        body: textView,
        alert: {
          message: 'message',
          type: 'default',
          status: 'error',
          closable: true,
        },
      },
      metadata: {
        type: 'card',
        title: [textView],
      },
    } as CardView;
    component.currentAction = action;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should cancel action', () => {
    component.onActionCancel();
    expect(component.currentAction).toBeUndefined();
  });

  it('should set action', () => {
    component.setAction(action);
    expect(component.currentAction).toBe(action);
  });

  it('should submit action', () => {
    const formGroup: FormGroup = formBuilder.group({
      formGroupExample: 'justForTest',
    });

    formFixture = TestBed.createComponent(FormComponent);
    formComponent = formFixture.componentInstance;
    formComponent.formGroup = formGroup;
    component.appForm = formComponent;

    component.onActionSubmit();
    expect(component.currentAction).toBeUndefined();
  });

  it('should not submit action', () => {
    component.onActionSubmit();
    expect(component.currentAction).toBeDefined();
  });

  it('should save title & body correctly', () => {
    const view: CardView = {
      config: {
        actions: [],
        body: {
          metadata: {
            type: 'textChanged',
            accessor: '',
            title: [],
          },
        },
      },
      metadata: {
        type: 'card',
      },
    };

    component.view = view;
    fixture.detectChanges();

    expect(component.body).toBe(view.config.body);
  });

  it('should call "onActionCancel" method when cancelling the form', fakeAsync(() => {
    spyOn(component, 'onActionCancel');
    component.onActionCancel();
    expect(component.onActionCancel).toHaveBeenCalled();
  }));
});
