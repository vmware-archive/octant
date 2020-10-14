import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import {
  BottomPanelComponent,
  minimizedHeight,
  PanelState,
  sliderHeightPropKey,
} from './bottom-panel.component';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { ResizableModule, ResizeEvent } from 'angular-resizable-element';
import { CssStyleDeclaration } from 'cytoscape';

describe('BottomPanelComponent', () => {
  let component: BottomPanelComponent;
  let fixture: ComponentFixture<BottomPanelComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [BottomPanelComponent],
        imports: [NoopAnimationsModule, ResizableModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(BottomPanelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('resizeCursors', () => {
    describe('panel is open', () => {
      beforeEach(() => {
        component.open = true;
      });

      it('sets topOrBottom cursor to ns-resize', () => {
        expect(component.resizeCursors().topOrBottom).toEqual('ns-resize');
      });
    });

    describe('panel is closed', () => {
      beforeEach(() => {
        component.open = false;
      });

      it('sets topOrBottom cursor to default', () => {
        expect(component.resizeCursors().topOrBottom).toEqual('default');
      });
    });
  });

  describe('updateSliderPosition', () => {
    const event: ResizeEvent = {
      rectangle: { bottom: 0, left: 0, right: 0, top: 777 },
      edges: {},
    };

    beforeEach(() => {
      spyOn(component, 'setHeight');
      spyOn(component, 'parentHeight').and.returnValue('100vh');
    });

    describe('panel is closed', () => {
      beforeEach(() => {
        component.open = false;
      });

      it('does not update the slider position', () => {
        component.updateSliderPosition(event);
        expect(component.setHeight).not.toHaveBeenCalled();
      });
    });

    describe('panel is open', () => {
      const expectedHeight = 'calc(100vh - 777px)';

      beforeEach(() => {
        component.open = true;
        component.updateSliderPosition(event);
      });

      it('sets the height of the panel', () => {
        expect(component.setHeight).toHaveBeenCalledWith(expectedHeight);
      });

      it('sets the previous open height for the panel', () => {
        expect(component.previousOpenHeight).toEqual(expectedHeight);
      });
    });
  });

  describe('setHeight', () => {
    it('sets the height of the panel', () => {
      component.setHeight('777px');
      const style = fixture.nativeElement.style as CssStyleDeclaration;
      const value = style.getPropertyValue(sliderHeightPropKey);
      expect(value).toEqual('777px');
    });
  });

  describe('toggle', () => {
    beforeEach(() => {
      spyOn(component, 'setHeight');
    });

    describe('when panel is open', () => {
      beforeEach(() => {
        component.open = true;
        component.toggle();
      });

      it('closes the panel', () => {
        expect(component.open).toBeFalse();
      });

      it('sets the toggle state to closed', () => {
        expect(component.toggleState).toEqual(PanelState.Closed);
      });

      it('sets the height to the minimized height', () => {
        expect(component.setHeight).toHaveBeenCalledWith(minimizedHeight);
      });
    });

    describe('when panel is closed', () => {
      beforeEach(() => {
        component.open = false;
        component.toggle();
      });

      it('opens the panel', () => {
        expect(component.open).toBeTrue();
      });

      it('sets the toggle state to open', () => {
        expect(component.toggleState).toEqual(PanelState.Open);
      });

      it('sets the height to the minimized height', () => {
        expect(component.setHeight).toHaveBeenCalledWith(
          component.previousOpenHeight
        );
      });
    });
  });

  describe('gutterClass', () => {
    it('returns if the gutter class should be activated', () => {
      component.open = false;
      expect(component.gutterClass()).toEqual({ open: false });
    });
  });
});
