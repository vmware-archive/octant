import {
  ChangeDetectorRef,
  Component,
  Input,
  OnInit,
  ElementRef,
  ViewChild,
  NgZone,
} from '@angular/core';
import { Subject } from 'rxjs';
import { ClrDatagridFilter, ClrDatagridFilterInterface } from '@clr/angular';
import { TableRow } from '../../../models/content';
import { RelativePipe } from '../../../pipes/relative/relative.pipe';

@Component({
  selector: 'app-content-text-filter',
  templateUrl: './content-text-filter.component.html',
  styleUrls: ['./content-text-filter.component.scss'],
})
export class ContentTextFilterComponent
  implements ClrDatagridFilterInterface<TableRow>, OnInit {
  @Input() column: string;
  @ViewChild('search') search: ElementRef;

  relativeTime: RelativePipe;
  changes = new Subject<any>();
  text = '';

  constructor(
    filterContainer: ClrDatagridFilter,
    private cd: ChangeDetectorRef,
    ngZone: NgZone
  ) {
    this.relativeTime = new RelativePipe(cd, ngZone);

    filterContainer.setFilter(this);
    filterContainer.openChange.subscribe(() => {
      setTimeout(() => {
        this.search.nativeElement.focus();
      });
    });
  }

  ngOnInit(): void {
    this.cd.detectChanges();
  }

  accepts(row: TableRow): boolean {
    if (this.text.length === 0) {
      return true;
    }
    const lowerCaseText = this.text.toLowerCase();
    return this.getStringDataForColumn(row, this.column).some(data =>
      data.toLowerCase().includes(lowerCaseText)
    );
  }

  getStringDataForColumn(row: TableRow, column: string): string[] {
    switch (row.data[column].metadata.type) {
      case 'link':
      case 'text':
        return [row.data[column].config.value];
      case 'timestamp':
        return [this.relativeTime.transform(row.data[column].config.timestamp)];
      case 'containers':
        return row.data[column].config.containers.map(
          container => container.name
        );
      case 'labels':
        return Object.entries(
          row.data[column].config.labels
        ).map((labels: any[]) => labels.join(':'));
      case 'selectors':
        return row.data[column].config.selectors.map(
          selector => selector.config.key + ':' + selector.config.value
        );
    }
    return [];
  }

  isActive(): boolean {
    return this.text.length !== 0;
  }

  onFilterChange(text: string) {
    this.text = text;
    this.changes.next(true);
  }

  reset() {
    this.onFilterChange('');
    this.search.nativeElement.focus();
  }
}
