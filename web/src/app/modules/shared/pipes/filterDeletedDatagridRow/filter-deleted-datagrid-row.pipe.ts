import { Pipe, PipeTransform } from '@angular/core';
import { TableRowWithMetadata } from '../../models/content';

@Pipe({ name: 'filterDeletedDatagridRow', pure: true })
export class FilterDeletedDatagridRowPipe implements PipeTransform {
  public transform(row: TableRowWithMetadata) {
    return row.isDeleted ? ['row-deleted'] : [];
  }
}
