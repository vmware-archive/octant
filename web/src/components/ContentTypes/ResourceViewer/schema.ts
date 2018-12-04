export interface Schema {
  adjacencyList: AdjacencyList;
  objects: ResourceObjects;
  selected: string;
}

export interface AdjacencyList { [key: string]: Edge[] }
export interface Edge {
  node: string;
  edge: string;
}

export interface ResourceObjects {
  [key: string]: ResourceObject;
}

export interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
  isNetwork?: boolean;
}
