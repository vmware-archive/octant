import {
  async,
  ComponentFixture,
  TestBed,
  fakeAsync,
  tick,
} from '@angular/core/testing';

import { CardComponent } from './card.component';
import { OverviewModule } from '../../overview.module';
import { CardView, TextView, Action } from '../../../../models/content';
import { ReactiveFormsModule, FormBuilder, FormGroup } from '@angular/forms';
import { SimpleChange } from '@angular/core';
import { ViewService } from '../../services/view/view.service';
import { viewServiceStub } from 'src/app/testing/view-service.stub';

describe('CardComponent', () => {
  let component: CardComponent;
  let fixture: ComponentFixture<CardComponent>;
  const formBuilder: FormBuilder = new FormBuilder();

  const action: Action = {
    name: 'actionName',
    title: 'actionTitle',
    form: {
      fields: [],
    },
  };

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [OverviewModule, ReactiveFormsModule],
      providers: [
        { provide: FormBuilder, useValue: formBuilder },
        { provide: ViewService, useValue: viewServiceStub },
      ],
    }).compileComponents();
  }));

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

    const cardView: CardView = {
      config: {
        actions: [],
        body: textView,
      },
      metadata: undefined,
    };

    component.view = cardView;
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

    component.onActionSubmit(formGroup);
    expect(component.currentAction).toBeUndefined();
  });

  it('should not submit action', () => {
    component.onActionSubmit({} as FormGroup);
    expect(component.currentAction).toBeDefined();
  });

  it('should save title & body correctly', () => {
    component.view = {
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
      metadata: undefined,
    };

    component.ngOnChanges({
      view: new SimpleChange(null, component.view, false),
    });
    fixture.detectChanges();

    expect(component.title).toBe('Just a title');
    expect(component.body).toBe(component.view.config.body);
  });

  it('should call "onActionCancel" method when cancelling the form', fakeAsync(() => {
    spyOn(component, 'onActionCancel');
    fixture.detectChanges();
    component.appForm.onFormCancel();
    tick();
    fixture.detectChanges();
    expect(component.onActionCancel).toHaveBeenCalled();
  }));
});
