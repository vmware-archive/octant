import { ChangeDetectorRef, Component, Input, OnInit } from '@angular/core';
import { Subject } from 'rxjs';
import { ClrDatagridFilter, ClrDatagridFilterInterface } from '@clr/angular';
import { TableFilter, TableRow } from '../../../models/content';

@Component({
  selector: 'app-content-text-filter',
  templateUrl: './content-text-filter.component.html',
  styleUrls: ['./content-text-filter.component.scss'],
})
export class ContentTextFilterComponent
  implements ClrDatagridFilterInterface<TableRow>, OnInit {
  @Input() filter: TableFilter;
  @Input() column: string;

  changes = new Subject<any>();
  text = '';

  constructor(
    private filterContainer: ClrDatagridFilter,
    private cd: ChangeDetectorRef
  ) {
    filterContainer.setFilter(this);
  }

  ngOnInit(): void {
    this.cd.detectChanges();
  }

  accepts(row: TableRow): boolean {
    if (this.text.trim().length === 0) {
      return true;
    }

    if (
      !row.data[this.column] ||
      !row.data[this.column].config ||
      !row.data[this.column].config.value
    ) {
      return false;
    }

    return row.data[this.column].config.value
      .toLowerCase()
      .includes(this.text.toLowerCase());
  }

  isActive(): boolean {
    return this.text.length > 0;
  }

  onFilterChange(text: string) {
    this.text = text;
    this.changes.next(true);
  }
}
