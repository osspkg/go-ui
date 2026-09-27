export const uiRPCMethods = {
  getView: "ui.get",
  dataCall: "data.call",
  action: "ui.action",
} as const;

export interface UIGetViewRequest {
  plugin?: string;
  view: string;
  context?: unknown;
}

export interface UIGetViewResponse {
  schema: import("./schema.js").ViewSchema;
}

export interface UIDataCallRequest {
  plugin?: string;
  view: string;
  source: string;
  operation: string;
  input?: unknown;
}

export interface UIActionRequest {
  plugin?: string;
  view: string;
  action: string;
  input?: unknown;
}
