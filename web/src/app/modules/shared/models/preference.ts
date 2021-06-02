export type PreferenceElement =
  | RadioElement
  | InputElement
  | LabelDropDown
  | TextElement;

export enum Operation {
  Equal = 'Equal',
}

export interface Condition {
  lhs: string;
  rhs: any;
  op: Operation;
}

interface Element {
  name: string;
  value: string;
  /**
   * If any of the disable conditions evaluate to true, this
   * element is disabled.
   */
  disableConditions?: Condition[];
}

export interface RadioElement extends Element {
  type: 'radio';
  config: {
    values: { label: string; value: string }[];
  };
}

export interface LabelDropDown extends Element {
  type: 'dropdown';
  label: string;
  metadata: {
    type: string;
    title: [
      {
        metadata: {
          type: string;
        };
        config: {
          value: string;
        };
      }
    ];
  };
  config: {
    position?: string;
    type: string;
    selection: string;
    useSelection: boolean;
    items: { name: string; type: string; label: string }[];
  };
}

export interface InputElement extends Element {
  type: 'input';
  config: {
    label: string;
    placeholder: string;
  };
}

export interface TextElement extends Element {
  type: 'text';
  textConfig: {
    config: {
      value: string;
      clipboardValue: string;
    };
  };
}

export interface PreferenceSection {
  name: string;
  elements: PreferenceElement[];
}

export interface PreferencePanel {
  name: string;
  sections: PreferenceSection[];
}

export interface Preferences {
  panels: PreferencePanel[];
  updateName: string;
}
