import { Component } from '@angular/core';
import {
  TableView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const addresses: TextView = {
  metadata: {
    type: 'text',
  },
  config: {
    value: 'Addresses',
  },
};

const address1: TextView = {
  config: {
    value: '192.168.64.7',
  },
  metadata: {
    type: 'text',
  },
};

const address2: TextView = {
  config: {
    value: 'minikube',
  },
  metadata: {
    type: 'text',
  },
};

const type1: TextView = {
  config: {
    value: 'InternalIP',
  },
  metadata: {
    type: 'text',
  },
};

const type2: TextView = {
  config: {
    value: 'Hostname',
  },
  metadata: {
    type: 'text',
  },
};

const view: TableView = {
  metadata: {
    type: 'table',
    title: [addresses],
  },
  config: {
    columns: [
      {
        name: 'Type',
        accessor: 'Type',
      },
      {
        name: 'Address',
        accessor: 'Address',
      },
    ],
    rows: [
      {
        Address: address1,
        Type: type1,
      },
      {
        Address: address2,
        Type: type2,
      },
    ],
    emptyContent: 'There are no addresses!',
    loading: false,
    filters: {},
  },
};

const code = `nodeAddressesColumns := component.NewTableCols("Type", "Address")
table := component.NewTable("Addresses", "There are no addresses!", nodeAddressesColumns)

for _, address := range node.Status.Addresses {
  row := component.TableRow{}
  row["Type"] = component.NewText(string(address.Type))
  row["Address"] = component.NewText(address.Address)

  table.Add(row)
}`;
const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-table-demo',
  templateUrl: './angular-table.demo.html',
})
export class AngularTableDemoComponent {
  view = view;
  code = code;
  json = json;
}
