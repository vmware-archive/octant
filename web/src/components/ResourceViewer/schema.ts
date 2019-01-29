export interface IResourceViewer {
  metadata: {
    type: 'resourceViewer';
    title: string;
  };
  config: {
    adjacencyList: {
      [key: string]: Array<{
        node: string;
        edge: string;
      }>;
    };
    objects: {
      [key: string]: ResourceObject;
    };
    selected: string;
  };
}

export interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
  isNetwork?: boolean;
}
