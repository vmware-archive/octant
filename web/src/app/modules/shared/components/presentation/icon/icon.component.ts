import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { IconView } from '../../../models/content';

@Component({
  selector: 'app-view-icon',
  templateUrl: './icon.component.html',
  styleUrls: ['./icon.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class IconComponent extends AbstractViewComponent<IconView> {
  shape: string;
  flip: string;
  size: string;

  styleAttr = '';
  classAttr = '';

  sizeHash = { sm: '16', md: '24', lg: '36"', xl: '48', xxl: '64' };
  directionHash = {
    up: 'transform: rotate(0deg);',
    right: 'transform: rotate(90deg);',
    down: 'transform: rotate(180deg);',
    left: 'transform: rotate(270deg);',
  };
  statusHash = {
    info: 'is-info ',
    success: 'is-success ',
    warning: 'is-warning ',
    danger: 'is-error ',
  };
  badge = {
    info: 'has-badge--info ',
    success: 'has-badge--success ',
    danger: 'has-badge ',
    'warning-triangle': 'has-alert ',
  };

  constructor(private cdr: ChangeDetectorRef) {
    super();
  }

  // This will work for version 5 of clarity but we are using
  // version 4 so function are been added to map the differences
  // dir needs to be change to direction on the template
  protected update(): void {
    const view = this.v;

    // reset values to be re-calculated
    this.styleAttr = '';
    this.classAttr = '';

    this.shape = view.config.shape;
    this.flip = view.config.flip;
    this.mapSize(view.config.size);
    this.mapDirection(view.config.direction);
    this.mapSolid(view.config.solid);
    this.mapStatus(view.config.status);
    this.mapInverse(view.config.inverse);
    this.mapBadge(view.config.badge);
    if (view.config.status === '' || !view.config.status) {
      this.setColor(view.config.color);
    }

    this.cdr.markForCheck();
  }

  mapSize(size: string): string {
    if (!size || size === '') {
      return;
    }
    this.size = this.sizeHash[size] || size;
  }

  mapDirection(direction: string) {
    if (!direction || direction === '') {
      return;
    }
    this.styleAttr += this.directionHash[direction] || '';
  }

  mapSolid(isSolid: boolean) {
    this.classAttr += isSolid ? 'is-solid ' : '';
  }

  mapStatus(status: string) {
    if (!status || status === '') {
      return;
    }
    this.classAttr += this.statusHash[status] || '';
  }

  mapInverse(inverse: boolean) {
    this.classAttr += inverse ? 'is-inverse ' : '';
  }

  mapBadge(badge: string) {
    if (!badge || badge === '') {
      return;
    }
    this.classAttr += this.badge[badge] || '';
  }

  setColor(color: string) {
    console.log('color', color);
    if (!color || color === '') {
      return;
    }
    this.styleAttr += `fill: ${color};`;
  }
}
