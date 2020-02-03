import { Component, Input, OnInit } from '@angular/core';
import { SingleStatView, View } from '../../../models/content';

@Component({
  selector: 'app-single-stat',
  templateUrl: './single-stat.component.html',
  styleUrls: ['./single-stat.component.scss'],
})
export class SingleStatComponent implements OnInit {
  v: SingleStatView;

  @Input() set view(v: View) {
    this.v = v as SingleStatView;
  }
  get view() {
    return this.v;
  }

  constructor() {}

  ngOnInit() {}
}
