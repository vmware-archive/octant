import { Component } from '@angular/core';
import { PortsView } from '../../../../../../../src/app/modules/shared/models/content';

const portsView1: PortsView = {
  config: {
    ports: [
      {
        metadata: {
          type: 'port',
        },
        config: {
          port: 80,
          protocol: 'TCP',
          state: {
            isForwardable: true,
          },
          buttonGroup: {
            metadata: {
              type: 'buttonGroup',
            },
            config: {
              buttons: [
                {
                  name: 'Start port forward',
                  payload: {
                    action: 'overview/startPortForward',
                    apiVersion: 'v1',
                    kind: 'Pod',
                    name: 'httpbin-db6d74d85-nltjq',
                    namespace: 'default',
                    port: 80,
                  },
                },
              ],
            },
          },
        },
      },
    ],
  },
  metadata: {
    type: 'ports',
  },
};

const portsView2: PortsView = {
  config: {
    ports: [
      {
        metadata: {
          type: 'port',
        },
        config: {
          port: 80,
          protocol: 'TCP',
          state: {
            isForwardable: true,
            isForwarded: true,
            port: 64247,
            id: '6267f3b1-41ed-45c0-a4b0-36c9fea314ca',
          },
          buttonGroup: {
            metadata: {
              type: 'buttonGroup',
            },
            config: {
              buttons: [
                {
                  name: 'Stop port forward',
                  payload: {
                    action: 'overview/stopPortForward',
                    id: '6267f3b1-41ed-45c0-a4b0-36c9fea314ca',
                  },
                },
              ],
            },
          },
        },
      },
    ],
  },
  metadata: {
    type: 'ports',
  },
};

const code = `apiVersion, kind := gvk.ToAPIVersionAndKind()
pfs := component.PortForwardState{}
port = component.NewPort(namespace, apiVersion, kind, portName, portNumber, portProtocol, pfs)
`;

const json1 = JSON.stringify(portsView1, null, 4);
const json2 = JSON.stringify(portsView2, null, 4);

@Component({
  selector: 'app-angular-ports-demo',
  templateUrl: './angular-ports.demo.html',
})
export class AngularPortsDemoComponent {
  view1 = portsView1;
  view2 = portsView2;
  json1 = json1;
  json2 = json2;
  code = code;
}
