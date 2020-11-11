import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CardListComponent } from './card-list.component';
import { CardListView } from '../../../models/content';
import { SharedModule } from '../../../shared.module';
import { IndicatorComponent } from '../indicator/indicator.component';

describe('CardListComponent', () => {
  let component: CardListComponent;
  let fixture: ComponentFixture<CardListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [IndicatorComponent],
      imports: [SharedModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CardListComponent);
    component = fixture.componentInstance;

    const view: CardListView = {
      config: { cards: [] },
      metadata: { type: 'cardsList' },
    };

    component.view = view;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
