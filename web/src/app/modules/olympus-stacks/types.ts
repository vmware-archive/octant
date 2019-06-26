export interface WorkloadList {
  currentStack: string;
  stackOptions: string[];
  channelFollowing: string;
  workloads: Workload[];
}

export interface Workload {
  name: string;
  lastUpdated: Date;
  revision: string;
  sourceImage: string;
  isPinned?: boolean;
  isFollowingChannel?: boolean;
  channelFollowing?: string;
  isComparedAgainst?: boolean;
  isMismatch?: boolean;
}
