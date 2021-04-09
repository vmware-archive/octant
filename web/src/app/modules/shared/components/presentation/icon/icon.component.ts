import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import '@cds/core/icon/register';
import {
  loadCoreIconSet,
  loadEssentialIconSet,
  loadCommerceIconSet,
  loadMediaIconSet,
  loadSocialIconSet,
  loadTravelIconSet,
  loadTextEditIconSet,
  loadTechnologyIconSet,
  loadChartIconSet,
} from '@cds/core/icon';
import { IconView } from '../../../models/content';

@Component({
  selector: 'app-view-icon',
  templateUrl: './icon.component.html',
  styleUrls: ['./icon.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class IconComponent extends AbstractViewComponent<IconView> {
  direction: string;
  shape: string;
  flip: string;
  size: string;
  isInverse: boolean;
  isSolid: boolean;
  badge: string;
  status: string;
  iconStyle: string;
  label: string;

  constructor(private cdr: ChangeDetectorRef) {
    super();
    loadCoreIconSet();
    loadEssentialIconSet();
    loadCommerceIconSet();
    loadMediaIconSet();
    loadSocialIconSet();
    loadTravelIconSet();
    loadTextEditIconSet();
    loadTechnologyIconSet();
    loadChartIconSet();
  }

  protected update(): void {
    const view = this.v;

    this.shape = view.config.shape;
    this.flip = view.config.flip;
    this.size = view.config.size;
    this.direction = view.config.direction;
    this.status = view.config.status;
    this.badge = view.config.badge;

    this.isInverse = view.config.inverse;
    this.isSolid = view.config.solid;

    this.iconStyle = '';
    if (view.config.color !== '') {
      this.iconStyle += `--color: ${view.config.color};`;
    }
    if (view.config.badgeColor !== '') {
      this.iconStyle += `--badge-color: ${view.config.badgeColor};`;
    }
    this.label = view.config.label;
    this.cdr.markForCheck();
  }
}
