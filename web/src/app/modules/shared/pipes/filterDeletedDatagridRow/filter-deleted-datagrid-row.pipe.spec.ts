import { FilterDeletedDatagridRowPipe } from './filter-deleted-datagrid-row.pipe';
import { TableRowWithMetadata } from '../../models/content';

describe('FilterDeletedDatagridRowPipe', () => {
  it('create an instance', () => {
    const pipe = new FilterDeletedDatagridRowPipe();
    expect(pipe).toBeTruthy();
  });

  it('filters a row', () => {
    const pipe = new FilterDeletedDatagridRowPipe();
    const row: TableRowWithMetadata = {
      data: null,
      isDeleted: true,
    };
    expect(pipe.transform(row)).toEqual(['row-deleted']);
  });

  it('does not filter a row', () => {
    const pipe = new FilterDeletedDatagridRowPipe();
    const row: TableRowWithMetadata = {
      data: {
        ['test']: {
          metadata: null,
        },
      },
      isDeleted: false,
    };
    expect(pipe.transform(row)).toEqual([]);
  });
});
