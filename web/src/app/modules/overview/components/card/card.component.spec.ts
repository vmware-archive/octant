import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CardComponent } from './card.component';
import { OverviewModule } from '../../overview.module';
import { CardView, TextView } from '../../../../models/content';

describe('CardComponent', () => {
  let component: CardComponent;
  let fixture: ComponentFixture<CardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [OverviewModule],
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
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
