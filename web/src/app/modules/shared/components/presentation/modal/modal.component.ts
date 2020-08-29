import { Component, ChangeDetectionStrategy, ChangeDetectorRef, OnInit, OnDestroy } from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { TitleView, ModalView, View } from '../../../models/content';
import { ModalService } from '../../../services/modal/modal.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-view-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ModalComponent extends AbstractViewComponent<ModalView> implements OnInit, OnDestroy{
  title: TitleView[];
  body: View;
  opened = false;
  size: string;

  private modalSubscription: Subscription;

  constructor(private modalService: ModalService, private cd: ChangeDetectorRef) {
    super();
  }

  ngOnInit() {
    this.modalSubscription = this.modalService.isOpened.subscribe(isOpened => {
      if (this.opened !== isOpened) {
        this.opened = isOpened;
        this.cd.markForCheck();
      }
    });
  }

  ngOnDestroy(): void {
    if (this.modalSubscription) {
      this.modalSubscription.unsubscribe();
    }
  }

  update() {
    this.title = this.v.metadata.title as TitleView[];
    this.body = this.v.config.body;
    this.size = this.v.config.size;
  }
}
