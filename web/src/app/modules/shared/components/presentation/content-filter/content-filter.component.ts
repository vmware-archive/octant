import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { ClrDatagridFilter, ClrDatagridFilterInterface } from '@clr/angular';
import { Subject } from 'rxjs';
import { TableFilter, TableRow, TextView } from '../../../models/content';
import trackByIdentity from '../../../../../util/trackBy/trackByIdentity';

@Component({
  selector: 'app-content-filter',
  templateUrl: './content-filter.component.html',
  styleUrls: ['./content-filter.component.scss'],
})
export class ContentFilterComponent
  implements ClrDatagridFilterInterface<TableRow>, OnInit {
  @Input() filter: TableFilter;
  @Input() column: string;

  changes = new Subject<any>();
  checkboxes: { [key: string]: boolean } = {};
  trackByIdentity = trackByIdentity;

  constructor(
    private filterContainer: ClrDatagridFilter,
    private cd: ChangeDetectorRef
  ) {
    filterContainer.setFilter(this);
  }

  ngOnInit(): void {
    this.filter.selected.forEach(value => (this.checkboxes[value] = true));
    this.cd.detectChanges();
  }

  accepts(row: TableRow): boolean {
    const selected = Object.entries(this.checkboxes)
      .filter(([_, value]) => value)
      .map(([key, _]) => key);

    if (!row.data[this.column]) {
      return false;
    }

    if (row.data[this.column].metadata.type !== 'text') {
      return false;
    }

    const view = row.data[this.column] as TextView;
    return selected.includes(view.config.value);
  }

  isActive(): boolean {
    return Object.keys(this.checkboxes).length > 0;
  }

  onFilterChange(name: string, e: boolean) {
    this.checkboxes[name] = e;
    this.changes.next(true);
  }
}
