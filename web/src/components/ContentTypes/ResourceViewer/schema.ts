export interface Schema {
  dag: DAG;
  objects: ResourceObjects;
  selected: string;
}

export type DAG = { [key: string]: Edge[] };

export interface Edge {
  node: string;
  edge: string;
}

export type ResourceObjects = { [key: string]: ResourceObject };

export interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
  isNetwork?: boolean;
}
