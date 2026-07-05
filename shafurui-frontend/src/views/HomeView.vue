<template>
  <div class="album-page">
    <header class="topbar">
      <div class="topbar-inner">
        <div class="brand">
          <div class="mark" aria-hidden="true">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M6 5.5h12a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-9a2 2 0 0 1 2-2Z"
                stroke="currentColor"
                stroke-width="1.8"
              />
              <path d="m10.2 9 4.7 3-4.7 3V9Z" fill="currentColor" />
            </svg>
          </div>
          <div>
            <h1>视频相册</h1>
            <p>按拍摄时间归档的私有视频库</p>
          </div>
        </div>

        <div class="search">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path
              d="m21 21-4.35-4.35M10.8 18a7.2 7.2 0 1 1 0-14.4 7.2 7.2 0 0 1 0 14.4Z"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
          <input
            v-model="keyword"
            type="search"
            placeholder="搜索文件名、日期或路径"
            autocomplete="off"
          />
        </div>

        <div class="toolbar" aria-label="视图工具栏">
          <div class="segmented" aria-label="卡片密度">
            <button
              :class="{ active: !compact }"
              type="button"
              title="标准密度"
              aria-label="标准密度"
              @click="compact = false"
            >
              <svg width="19" height="19" viewBox="0 0 24 24" fill="none">
                <path
                  d="M4 5h7v6H4V5Zm9 0h7v6h-7V5ZM4 13h7v6H4v-6Zm9 0h7v6h-7v-6Z"
                  stroke="currentColor"
                  stroke-width="1.8"
                />
              </svg>
            </button>
            <button
              :class="{ active: compact }"
              type="button"
              title="紧凑密度"
              aria-label="紧凑密度"
              @click="compact = true"
            >
              <svg width="19" height="19" viewBox="0 0 24 24" fill="none">
                <path
                  d="M4 4h5v5H4V4Zm5.5 0h5v5h-5V4ZM15 4h5v5h-5V4ZM4 9.5h5v5H4v-5Zm5.5 0h5v5h-5v-5Zm5.5 0h5v5h-5v-5ZM4 15h5v5H4v-5Zm5.5 0h5v5h-5v-5Zm5.5 0h5v5h-5v-5Z"
                  stroke="currentColor"
                  stroke-width="1.5"
                />
              </svg>
            </button>
          </div>
          <button
            class="icon-button"
            type="button"
            title="清空筛选"
            aria-label="清空筛选"
            @click="clearFilters"
          >
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none">
              <path
                d="M4 7h16M9 11v6m6-6v6M10 4h4l1 3H9l1-3Zm-3 3 1 13h8l1-13"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
          <RouterLink class="icon-button" to="/signin" title="登录" aria-label="登录">
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none">
              <path
                d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M15 12H3"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </RouterLink>
        </div>
      </div>
    </header>

    <div class="shell">
      <aside class="sidebar">
        <div class="stats">
          <div class="stat">
            <b>{{ filteredVideos.length }}</b>
            <span>视频总数</span>
          </div>
          <div class="stat">
            <b>{{ totalHours }}h</b>
            <span>总时长</span>
          </div>
        </div>

        <div class="filter-panel">
          <div class="select-row">
            <label for="monthFilter">月份</label>
            <select id="monthFilter" v-model="monthFilter">
              <option value="all">全部月份</option>
              <option v-for="month in months" :key="month" :value="month">{{ month }}</option>
            </select>
          </div>
          <div class="select-row">
            <label for="sourceFilter">来源</label>
            <select id="sourceFilter" v-model="sourceFilter">
              <option value="all">全部来源</option>
              <option v-for="source in sources" :key="source" :value="source">{{ source }}</option>
            </select>
          </div>
        </div>

        <nav class="date-nav" aria-label="日期导航">
          <button
            v-for="(group, index) in groupedVideos"
            :key="group.date"
            class="date-link"
            :class="{ active: index === 0 }"
            type="button"
            @click="scrollToGroup(group.date)"
          >
            <strong>{{ dateLabel(group.date) }}</strong>
            <span>{{ group.items.length }}</span>
          </button>
        </nav>
      </aside>

      <main class="main">
        <section class="overview">
          <div>
            <h2>最近同步的视频已经按日期整理好</h2>
            <p>
              封面来自同名图片或视频首帧，列表接口只返回结构化数据，播放仍由 nginx 静态文件负责。
            </p>
          </div>
          <div class="health" aria-label="扫描状态">
            <div class="health-row"><span>API</span><b>/api/video</b></div>
            <div class="health-row">
              <span>数据源</span><b>{{ dataSource }}</b>
            </div>
            <div class="health-row">
              <span>最近扫描</span><b>{{ scanTime }}</b>
            </div>
            <div class="health-row">
              <span>封面缺失</span><b>{{ missingCovers }}</b>
            </div>
          </div>
        </section>

        <div v-if="loading" class="empty">正在加载视频列表...</div>
        <template v-else>
          <section
            v-for="group in groupedVideos"
            :id="groupId(group.date)"
            :key="group.date"
            class="group"
          >
            <div class="group-head">
              <div class="group-title">
                <h3>{{ group.date }}</h3>
                <span>{{ dateLabel(group.date) }}</span>
              </div>
              <div class="group-meta">
                {{ group.items.length }} 个视频 / {{ groupMinutes(group.items) }} 分钟
              </div>
            </div>
            <div class="grid" :class="{ compact }">
              <article v-for="video in group.items" :key="video.id" class="video-card">
                <button
                  class="thumb-button"
                  type="button"
                  :aria-label="`播放 ${video.filename}`"
                  @click="openPlayer(video)"
                >
                  <img
                    v-if="video.coverUrl"
                    :src="video.coverUrl"
                    loading="lazy"
                    :alt="`${video.filename} 封面`"
                    @error="markCoverMissing(video.id)"
                  />
                  <div v-else class="poster-fallback" :style="{ '--hue': String(video.hue) }">
                    <span>{{ video.groupDate }}</span>
                  </div>
                  <span class="badge">{{
                    video.coverUrl && !missingCoverIds.has(video.id) ? "COVER" : "FIRST FRAME"
                  }}</span>
                  <span class="play" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                      <path d="M8 5v14l11-7L8 5Z" fill="currentColor" />
                    </svg>
                  </span>
                  <span class="duration">{{ formatDuration(video.durationSec) }}</span>
                </button>
                <div class="card-body">
                  <p class="filename" :title="video.filename">{{ video.filename }}</p>
                  <div class="meta">
                    <span>{{ formatShotTime(video) }}</span>
                    <code>{{ resolution(video) }}</code>
                  </div>
                  <div class="meta">
                    <span>{{ sourceLabel(video) }}</span>
                    <code :title="video.relativePath">{{ video.relativePath }}</code>
                  </div>
                </div>
              </article>
            </div>
          </section>
          <div v-if="groupedVideos.length === 0" class="empty">没有匹配的视频</div>
        </template>
      </main>
    </div>

    <div
      class="modal"
      :class="{ open: Boolean(currentVideo) }"
      role="dialog"
      aria-modal="true"
      aria-labelledby="playerTitle"
      @click.self="closePlayer"
    >
      <div class="player">
        <div class="player-bar">
          <div>
            <h4 id="playerTitle">{{ currentVideo?.filename }}</h4>
            <p>
              {{
                currentVideo
                  ? `${formatShotTime(currentVideo)} / ${resolution(currentVideo)} / ${currentVideo.relativePath}`
                  : ""
              }}
            </p>
          </div>
          <button
            class="close"
            type="button"
            title="关闭"
            aria-label="关闭播放器"
            @click="closePlayer"
          >
            <svg width="19" height="19" viewBox="0 0 24 24" fill="none">
              <path
                d="M6 6l12 12M18 6 6 18"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
              />
            </svg>
          </button>
        </div>
        <div class="video-stage">
          <video
            v-if="currentVideo"
            :key="currentVideo.id"
            :src="currentVideo.url"
            :poster="currentVideo.coverUrl"
            controls
            autoplay
            playsinline
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { RouterLink } from "vue-router";
import { getVideoList, type VideoItem } from "@/api/video";

defineOptions({
  name: "HomeView",
});

type AlbumVideo = VideoItem & {
  hue: number;
};

const sampleVideos: AlbumVideo[] = [
  sample("IMG_20260702_184522.mp4", "2026-07-02T18:45:22+08:00", 84, 3840, 2160, "iPhone", 12),
  sample(
    "VID_20260702_091038.mov",
    "2026-07-02T09:10:38+08:00",
    38,
    1920,
    1080,
    "iPhone",
    178,
    false,
  ),
  sample("park_walk_20260701.mp4", "2026-07-01T20:18:00+08:00", 126, 3840, 2160, "Camera", 142),
  sample("family_table_20260701.mp4", "2026-07-01T12:04:00+08:00", 57, 1920, 1080, "iPhone", 42),
  sample("rain_window_20260630.m4v", "2026-06-30T23:31:00+08:00", 44, 1920, 1080, "Camera", 213),
  sample(
    "train_station_20260630.mp4",
    "2026-06-30T08:22:00+08:00",
    203,
    3840,
    2160,
    "Camera",
    31,
    false,
  ),
  sample("birthday_clip_20260618.mov", "2026-06-18T19:45:00+08:00", 72, 1920, 1080, "iPhone", 342),
  sample("kitchen_test_20260618.mp4", "2026-06-18T17:02:00+08:00", 29, 1280, 720, "Other", 88),
];

const videos = ref<AlbumVideo[]>([]);
const keyword = ref("");
const monthFilter = ref("all");
const sourceFilter = ref("all");
const compact = ref(false);
const loading = ref(true);
const dataSource = ref("sample");
const scanTime = ref("--:--");
const currentVideo = ref<AlbumVideo | null>(null);
const missingCoverIds = ref(new Set<string>());

const filteredVideos = computed(() => {
  const query = keyword.value.trim().toLowerCase();
  return videos.value.filter((video) => {
    const matchedKeyword =
      !query ||
      video.filename.toLowerCase().includes(query) ||
      video.relativePath.toLowerCase().includes(query) ||
      video.groupDate.includes(query);
    const matchedMonth =
      monthFilter.value === "all" || video.groupDate.startsWith(monthFilter.value);
    const matchedSource = sourceFilter.value === "all" || sourceLabel(video) === sourceFilter.value;
    return matchedKeyword && matchedMonth && matchedSource;
  });
});

const groupedVideos = computed(() => {
  const map = new Map<string, AlbumVideo[]>();
  for (const video of filteredVideos.value) {
    if (!map.has(video.groupDate)) {
      map.set(video.groupDate, []);
    }
    map.get(video.groupDate)?.push(video);
  }
  return [...map.entries()]
    .sort(([left], [right]) => right.localeCompare(left))
    .map(([date, items]) => ({
      date,
      items: items.sort((left, right) => (right.shotAt || "").localeCompare(left.shotAt || "")),
    }));
});

const months = computed(() =>
  [...new Set(videos.value.map((video) => video.groupDate.slice(0, 7)))].sort().reverse(),
);
const sources = computed(() => [...new Set(videos.value.map(sourceLabel))].sort());
const totalHours = computed(() =>
  (filteredVideos.value.reduce((sum, video) => sum + (video.durationSec || 0), 0) / 3600).toFixed(
    1,
  ),
);
const missingCovers = computed(
  () =>
    videos.value.filter((video) => !video.coverUrl || missingCoverIds.value.has(video.id)).length,
);

onMounted(async () => {
  await loadVideos();
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});

async function loadVideos() {
  loading.value = true;
  try {
    const response = await getVideoList();
    const items = response.groups.flatMap((group) =>
      group.items.map((item) => normalizeVideo(item, group.date)),
    );
    videos.value = items.length > 0 ? items : sampleVideos;
    dataSource.value = items.length > 0 ? "api" : "sample";
  } catch (error) {
    console.warn("Failed to load video list, using sample data.", error);
    videos.value = sampleVideos;
    dataSource.value = "sample";
  } finally {
    scanTime.value = new Intl.DateTimeFormat("zh-CN", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    }).format(new Date());
    loading.value = false;
  }
}

function sample(
  filename: string,
  shotAt: string,
  durationSec: number,
  width: number,
  height: number,
  source: string,
  hue: number,
  hasCover = true,
): AlbumVideo {
  const groupDate = shotAt.slice(0, 10);
  const relativePath = `${groupDate.slice(0, 4)}/${groupDate.slice(5, 7)}/${filename}`;
  return {
    id: relativePath,
    filename,
    relativePath,
    url: `/videos/${relativePath}`,
    coverUrl: hasCover ? generatedPoster(groupDate, hue) : undefined,
    shotAt,
    groupDate,
    durationSec,
    width,
    height,
    sizeBytes: durationSec * 2800000,
    mtime: shotAt,
    hue,
    source,
  } as AlbumVideo & { source: string };
}

function normalizeVideo(item: VideoItem, fallbackDate: string): AlbumVideo {
  return {
    ...item,
    id: item.id || item.relativePath || item.url,
    groupDate: item.groupDate || fallbackDate,
    relativePath: item.relativePath || item.filename,
    hue: hueFromString(item.relativePath || item.filename),
  };
}

function generatedPoster(date: string, hue: number) {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="960" height="540" viewBox="0 0 960 540">
    <defs><linearGradient id="g" x1="0" y1="0" x2="1" y2="1"><stop stop-color="hsl(${hue},48%,24%)"/><stop offset=".52" stop-color="hsl(${(hue + 34) % 360},42%,34%)"/><stop offset="1" stop-color="hsl(${(hue + 174) % 360},28%,16%)"/></linearGradient></defs>
    <rect width="960" height="540" fill="url(#g)"/>
    <g opacity=".22">${Array.from({ length: 18 }, (_, index) => `<rect x="${index * 66 - 120}" width="16" height="540" fill="${index % 3 === 0 ? "#fff" : "#111"}"/>`).join("")}</g>
    <rect x="48" y="398" width="420" height="64" fill="rgba(255,255,255,.84)"/>
    <text x="72" y="438" fill="#111113" font-size="32" font-weight="700" font-family="Fira Sans, sans-serif">${date}</text>
  </svg>`;
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
}

function hueFromString(value: string) {
  return [...value].reduce((sum, char) => sum + char.charCodeAt(0), 0) % 360;
}

function sourceLabel(video: AlbumVideo) {
  const unknownVideo = video as AlbumVideo & { source?: string };
  if (unknownVideo.source) return unknownVideo.source;
  if (/iphone|img_|vid_/i.test(video.filename)) return "iPhone";
  if (/camera|dji|gopro/i.test(video.filename)) return "Camera";
  return "Other";
}

function formatDuration(seconds = 0) {
  const minutes = Math.floor(seconds / 60);
  const rest = String(Math.floor(seconds % 60)).padStart(2, "0");
  return `${minutes}:${rest}`;
}

function formatShotTime(video: AlbumVideo) {
  if (!video.shotAt) return video.groupDate;
  const date = new Date(video.shotAt);
  if (Number.isNaN(date.getTime())) return video.groupDate;
  return new Intl.DateTimeFormat("zh-CN", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(date);
}

function dateLabel(date: string) {
  const today = new Date();
  const yesterday = new Date();
  yesterday.setDate(today.getDate() - 1);
  if (date === today.toISOString().slice(0, 10)) return "今天";
  if (date === yesterday.toISOString().slice(0, 10)) return "昨天";
  return date;
}

function resolution(video: AlbumVideo) {
  if (!video.width || !video.height) return "UNKNOWN";
  if (video.width >= 3840 || video.height >= 2160) return "4K";
  if (video.width >= 1920 || video.height >= 1080) return "1080p";
  if (video.width >= 1280 || video.height >= 720) return "720p";
  return `${video.width}x${video.height}`;
}

function groupMinutes(items: AlbumVideo[]) {
  return (items.reduce((sum, video) => sum + (video.durationSec || 0), 0) / 60).toFixed(1);
}

function groupId(date: string) {
  return `group-${date}`;
}

function scrollToGroup(date: string) {
  document
    .querySelector(`#${groupId(date)}`)
    ?.scrollIntoView({ behavior: "smooth", block: "start" });
}

function clearFilters() {
  keyword.value = "";
  monthFilter.value = "all";
  sourceFilter.value = "all";
}

function openPlayer(video: AlbumVideo) {
  currentVideo.value = video;
}

function closePlayer() {
  currentVideo.value = null;
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === "Escape") closePlayer();
}

function markCoverMissing(id: string) {
  missingCoverIds.value = new Set([...missingCoverIds.value, id]);
}
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Fira+Code:wght@500;600;700&family=Fira+Sans:wght@400;500;600;700&display=swap");

.album-page {
  --bg: #f6f5f2;
  --panel: #ffffff;
  --ink: #101013;
  --muted: #696966;
  --line: #d9d6cf;
  --line-strong: #b9b5aa;
  --dark: #171719;
  --accent: #c43f2f;
  --good: #18794e;
  --shadow: 0 18px 50px rgba(18, 18, 20, 0.12);
  --radius: 8px;
  min-height: 100vh;
  color: var(--ink);
  background:
    linear-gradient(90deg, rgba(16, 16, 19, 0.045) 1px, transparent 1px) 0 0 / 44px 44px,
    linear-gradient(rgba(16, 16, 19, 0.035) 1px, transparent 1px) 0 0 / 44px 44px,
    var(--bg);
  font-family: "Fira Sans", "PingFang SC", "Microsoft YaHei", sans-serif;
}

button,
input,
select {
  font: inherit;
}

button {
  cursor: pointer;
}

.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  background: rgba(246, 245, 242, 0.92);
  border-bottom: 1px solid var(--line);
  backdrop-filter: blur(18px);
}

.topbar-inner {
  max-width: 1480px;
  margin: 0 auto;
  min-height: 88px;
  padding: 18px 28px;
  display: grid;
  grid-template-columns: minmax(260px, 1fr) minmax(380px, 560px) auto;
  gap: 18px;
  align-items: center;
}

.brand {
  display: flex;
  gap: 14px;
  align-items: center;
  min-width: 0;
}

.mark {
  width: 48px;
  height: 48px;
  border-radius: var(--radius);
  background:
    linear-gradient(135deg, transparent 0 42%, rgba(255, 255, 255, 0.28) 43% 56%, transparent 57%),
    var(--dark);
  display: grid;
  place-items: center;
  color: #fff;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.12);
  flex: 0 0 auto;
}

.brand h1,
.overview h2,
.group-title h3,
.player-bar h4,
.filename {
  margin: 0;
  letter-spacing: 0;
}

.brand h1 {
  font-size: 24px;
  line-height: 1.05;
}

.brand p {
  margin: 5px 0 0;
  color: var(--muted);
  font-size: 13px;
}

.search {
  position: relative;
}

.search svg {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--muted);
}

.search input {
  width: 100%;
  height: 48px;
  border-radius: var(--radius);
  border: 1px solid var(--line-strong);
  background: var(--panel);
  padding: 0 16px 0 46px;
  color: var(--ink);
  outline: none;
}

.search input:focus,
.select-row select:focus {
  border-color: var(--dark);
  box-shadow: 0 0 0 3px rgba(16, 16, 19, 0.12);
}

.toolbar {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  align-items: center;
}

.icon-button,
.segmented button {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border-radius: var(--radius);
  border: 1px solid var(--line-strong);
  background: var(--panel);
  color: var(--dark);
  text-decoration: none;
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    background 180ms ease;
}

.icon-button:hover,
.segmented button:hover {
  transform: translateY(-1px);
  border-color: var(--dark);
}

.segmented {
  display: flex;
  gap: 6px;
  padding: 4px;
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.58);
  border-radius: 10px;
}

.segmented button.active {
  background: var(--dark);
  color: #fff;
  border-color: var(--dark);
}

.shell {
  width: 100%;
  max-width: 1480px;
  margin: 0 auto;
  padding: 24px 28px 44px;
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 24px;
}

.sidebar {
  position: sticky;
  top: 112px;
  align-self: start;
  display: grid;
  gap: 14px;
}

.stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.stat,
.date-nav,
.filter-panel {
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.7);
  border-radius: var(--radius);
}

.stat {
  padding: 14px;
  min-height: 82px;
}

.stat b {
  display: block;
  font-family: "Fira Code", monospace;
  font-size: 23px;
  line-height: 1;
}

.stat span {
  display: block;
  margin-top: 8px;
  color: var(--muted);
  font-size: 12px;
}

.filter-panel {
  padding: 12px;
  display: grid;
  gap: 10px;
}

.select-row {
  display: grid;
  gap: 6px;
}

.select-row label {
  color: var(--muted);
  font-size: 12px;
}

.select-row select {
  width: 100%;
  height: 38px;
  border-radius: 7px;
  border: 1px solid var(--line);
  background: #fff;
  color: var(--ink);
  padding: 0 10px;
  outline: none;
}

.date-nav {
  padding: 8px;
}

.date-link {
  width: 100%;
  border: 0;
  background: transparent;
  border-radius: 7px;
  padding: 10px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
  text-align: left;
  color: var(--ink);
  transition:
    background 180ms ease,
    color 180ms ease;
}

.date-link:hover,
.date-link.active {
  background: var(--dark);
  color: #fff;
}

.date-link strong {
  font-size: 13px;
}

.date-link span {
  font-family: "Fira Code", monospace;
  font-size: 12px;
  opacity: 0.78;
}

.main {
  min-width: 0;
}

.overview {
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: var(--dark);
  color: #fff;
  min-height: 172px;
  padding: 22px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 22px;
  overflow: hidden;
  position: relative;
}

.overview::after {
  content: "";
  position: absolute;
  inset: auto -8% -40% 34%;
  height: 170px;
  background:
    repeating-linear-gradient(90deg, rgba(255, 255, 255, 0.22) 0 1px, transparent 1px 18px),
    linear-gradient(90deg, transparent, rgba(196, 63, 47, 0.36), transparent);
  transform: rotate(-6deg);
}

.overview h2 {
  font-size: 34px;
  line-height: 1.05;
  max-width: 720px;
  position: relative;
  z-index: 1;
}

.overview p {
  margin: 12px 0 0;
  color: rgba(255, 255, 255, 0.68);
  max-width: 660px;
  position: relative;
  z-index: 1;
}

.health {
  min-width: 250px;
  position: relative;
  z-index: 1;
  align-self: stretch;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: var(--radius);
  padding: 14px;
  background: rgba(255, 255, 255, 0.06);
}

.health-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12);
  font-size: 13px;
}

.health-row:last-child {
  border-bottom: 0;
}

.health-row span {
  color: rgba(255, 255, 255, 0.64);
}

.health-row b {
  font-family: "Fira Code", monospace;
  color: #fff;
}

.group {
  scroll-margin-top: 116px;
  margin-top: 24px;
}

.group-head {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 18px;
  padding: 4px 2px 12px;
  border-bottom: 1px solid var(--line);
}

.group-title {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.group-title h3 {
  font-family: "Fira Code", monospace;
  font-size: 20px;
}

.group-title span,
.group-meta {
  color: var(--muted);
  font-size: 13px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  padding-top: 14px;
}

.grid.compact {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.video-card {
  border: 1px solid var(--line);
  background: var(--panel);
  border-radius: var(--radius);
  overflow: hidden;
  min-width: 0;
  transition:
    transform 180ms ease,
    box-shadow 180ms ease,
    border-color 180ms ease;
}

.video-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow);
  border-color: var(--dark);
}

.thumb-button {
  width: 100%;
  border: 0;
  padding: 0;
  background: #111;
  display: block;
  position: relative;
  aspect-ratio: 16 / 9;
  overflow: hidden;
}

.thumb-button img,
.poster-fallback {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  filter: saturate(0.9) contrast(1.03);
  transition: transform 240ms ease;
}

.video-card:hover img,
.video-card:hover .poster-fallback {
  transform: scale(1.035);
}

.poster-fallback {
  display: grid;
  place-items: end start;
  padding: 28px;
  color: #111113;
  background:
    linear-gradient(
      135deg,
      hsl(var(--hue), 48%, 24%),
      hsl(calc(var(--hue) + 34), 42%, 34%) 52%,
      hsl(calc(var(--hue) + 174), 28%, 16%)
    ),
    #111;
}

.poster-fallback span {
  min-width: 210px;
  padding: 12px 18px;
  background: rgba(255, 255, 255, 0.84);
  font-family: "Fira Code", monospace;
  font-weight: 700;
}

.play {
  position: absolute;
  left: 12px;
  bottom: 12px;
  width: 38px;
  height: 38px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  color: #fff;
  background: rgba(16, 16, 19, 0.82);
  border: 1px solid rgba(255, 255, 255, 0.26);
}

.duration {
  position: absolute;
  right: 10px;
  bottom: 12px;
  color: #fff;
  background: rgba(16, 16, 19, 0.82);
  border-radius: 5px;
  padding: 5px 7px;
  font-family: "Fira Code", monospace;
  font-size: 12px;
}

.badge {
  position: absolute;
  left: 10px;
  top: 10px;
  color: #111;
  background: rgba(255, 255, 255, 0.86);
  border-radius: 5px;
  padding: 5px 7px;
  font-size: 12px;
  font-weight: 700;
}

.card-body {
  padding: 12px;
  display: grid;
  gap: 10px;
}

.filename {
  font-size: 15px;
  line-height: 1.25;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
  color: var(--muted);
  font-size: 12px;
  min-width: 0;
}

.meta span,
.meta code {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta code {
  color: var(--dark);
  font-family: "Fira Code", monospace;
  font-size: 11px;
}

.empty {
  margin-top: 24px;
  border: 1px dashed var(--line-strong);
  border-radius: var(--radius);
  padding: 42px 18px;
  text-align: center;
  background: rgba(255, 255, 255, 0.6);
  color: var(--muted);
}

.modal {
  position: fixed;
  inset: 0;
  z-index: 50;
  display: none;
  place-items: center;
  padding: 28px;
  background: rgba(10, 10, 12, 0.78);
}

.modal.open {
  display: grid;
}

.player {
  width: min(1080px, 100%);
  background: #0f0f12;
  color: #fff;
  border-radius: var(--radius);
  box-shadow: 0 30px 90px rgba(0, 0, 0, 0.48);
  overflow: hidden;
}

.player-bar {
  min-height: 58px;
  padding: 12px 14px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12);
}

.player-bar h4 {
  font-size: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-bar p {
  margin: 3px 0 0;
  color: rgba(255, 255, 255, 0.58);
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.close {
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  border: 1px solid rgba(255, 255, 255, 0.18);
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  display: grid;
  place-items: center;
}

.video-stage {
  aspect-ratio: 16 / 9;
  display: grid;
  place-items: center;
  background: #050506;
}

.video-stage video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

@media (max-width: 1120px) {
  .topbar-inner {
    grid-template-columns: 1fr;
  }

  .toolbar {
    justify-content: flex-start;
  }

  .shell {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
    grid-template-columns: 1fr 1fr;
  }

  .date-nav {
    max-height: 220px;
    overflow: auto;
  }

  .grid,
  .grid.compact {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .topbar-inner,
  .shell {
    padding-left: 16px;
    padding-right: 16px;
  }

  .brand h1 {
    font-size: 21px;
  }

  .sidebar,
  .overview {
    grid-template-columns: 1fr;
  }

  .health {
    min-width: 0;
  }

  .grid,
  .grid.compact {
    grid-template-columns: 1fr;
  }

  .group-head {
    align-items: start;
    flex-direction: column;
  }

  .modal {
    padding: 14px;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    scroll-behavior: auto !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
