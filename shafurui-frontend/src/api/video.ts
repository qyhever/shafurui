import { get } from "@/utils/request";

export interface VideoItem {
  id: string;
  filename: string;
  relativePath: string;
  url: string;
  coverUrl?: string;
  shotAt?: string;
  groupDate: string;
  durationSec?: number;
  width?: number;
  height?: number;
  sizeBytes?: number;
  mtime?: string;
}

export interface VideoGroup {
  date: string;
  items: VideoItem[];
}

export interface VideoListResponse {
  groups: VideoGroup[];
}

export function getVideoList() {
  return get<VideoListResponse>("/video");
}
