declare module "dagre" {
  import * as d3 from "d3";
  import * as dagre from "dagre";
  export * from "dagre";

  export const render: { new(): Render };
  export const intersect: { [shapeName: string]: (node: Node, points: Array<{}>, point: any) => void };

    namespace graphlib {

        interface Graph {
            graph(): Graph;
            height: number;
            predecessors(id: string): string[];
            successors(id: string): string[];
            transition?(selection: d3.Selection<any, any, any, any>): d3.Transition<any, any, any, any>;
            width: number;
        }
    }

    export type Edge = {[key:string]:any}

    export interface Render {
        // see http://cpettitt.github.io/project/dagre-d3/latest/demo/user-defined.html for example usage
        arrows(): { [arrowStyleName: string]: (parent: d3.Selection<any, any, any, any>, id: string, edge: Edge, type: string) => void };
        (selection: d3.Selection<any, any, any, any>, g: graphlib.Graph): void;
        shapes(): { [shapeStyleName: string]: (parent: d3.Selection<any, any, any, any>, bbox: any, node: Node) => void };
    }
}
