import type { ReactElement } from 'react';
import React from 'react';
import { Globe, FileEdit, CheckCircle2, XCircle, Clock } from 'lucide-react';

export interface StatusConfig {
  icon: ReactElement;
  textColor: string;
  label: string;
}

export const getPublishDraftStatusConfig = (status: string): StatusConfig => {
  switch (status) {
    case "publish":
      return {
        icon: React.createElement(Globe, { className: "h-3.5 w-3.5 text-chart-2" }),
        textColor: "text-chart-2",
        label: "เผยแพร่"
      };
    case "draft":
      return {
        icon: React.createElement(FileEdit, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: "ฉบับร่าง"
      };
    default:
      return {
        icon: React.createElement(FileEdit, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: "ฉบับร่าง"
      };
  }
};

export const getPublishDraftStatusConfigForBanner = (status: string): StatusConfig => {
  switch (status) {
    case "publish":
      return {
        icon: React.createElement(Globe, { className: "h-3.5 w-3.5 text-chart-2" }),
        textColor: "text-chart-2",
        label: "เผยแพร่"
      };
    case "draft":
      return {
        icon: React.createElement(FileEdit, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: "ร่าง"
      };
    default:
      return {
        icon: React.createElement(FileEdit, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: "ร่าง"
      };
  }
};

export const getActiveInactiveStatusConfig = (active: boolean): StatusConfig => {
  return active
    ? {
        icon: React.createElement(CheckCircle2, { className: "h-3.5 w-3.5 text-chart-2" }),
        textColor: "text-chart-2",
        label: "ใช้งาน"
      }
    : {
        icon: React.createElement(XCircle, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: "ไม่ใช้งาน"
      };
};

export const getUserStatusConfig = (status: string): StatusConfig => {
  switch (status) {
    case "active":
      return {
        icon: React.createElement(CheckCircle2, { className: "h-3.5 w-3.5 text-chart-2" }),
        textColor: "text-chart-2",
        label: "ยืนยันตัวตน"
      };
    case "pending":
      return {
        icon: React.createElement(Clock, { className: "h-3.5 w-3.5 text-chart-3" }),
        textColor: "text-chart-3",
        label: "รออนุมัติ"
      };
    case "inactive":
      return {
        icon: React.createElement(XCircle, { className: "h-3.5 w-3.5 text-destructive" }),
        textColor: "text-destructive",
        label: "ถูกระงับ"
      };
    default:
      return {
        icon: React.createElement(Clock, { className: "h-3.5 w-3.5 text-muted-foreground" }),
        textColor: "text-muted-foreground",
        label: status
      };
  }
};

export const getPeriodDialogStatusIcon = (status: string): ReactElement => {
  switch (status) {
    case "publish":
      return React.createElement(Globe, { className: "h-4 w-4 text-green-600" });
    case "draft":
      return React.createElement(FileEdit, { className: "h-4 w-4 text-muted-foreground" });
    default:
      return React.createElement(FileEdit, { className: "h-4 w-4 text-muted-foreground" });
  }
};

export const getPeriodDialogStatusLabel = (status: string): string => {
  return status === "publish" ? "เผยแพร่" : "ฉบับร่าง";
};