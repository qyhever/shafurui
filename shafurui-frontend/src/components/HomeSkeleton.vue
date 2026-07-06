<template>
  <template v-if="section === 'dates'">
    <div v-for="index in 5" :key="index" class="date-link-skeleton" aria-hidden="true">
      <span class="skeleton-block skeleton-date-title"></span>
      <span class="skeleton-block skeleton-date-count"></span>
    </div>
  </template>

  <div v-else class="loading-groups" aria-label="正在加载视频列表" aria-busy="true">
    <section v-for="groupIndex in 2" :key="groupIndex" class="group skeleton-group">
      <div class="group-head">
        <div class="group-title skeleton-group-title">
          <span class="skeleton-block skeleton-group-date"></span>
          <span class="skeleton-block skeleton-group-label"></span>
        </div>
        <span class="skeleton-block skeleton-group-meta"></span>
      </div>
      <div class="video-list" :class="viewMode">
        <article v-for="cardIndex in 3" :key="cardIndex" class="video-card skeleton-card">
          <div class="skeleton-thumb">
            <span class="skeleton-badge"></span>
            <span class="skeleton-play"></span>
            <span class="skeleton-duration"></span>
          </div>
          <div v-if="viewMode === 'list'" class="list-body skeleton-list-body">
            <div class="list-primary">
              <span class="skeleton-block skeleton-file"></span>
              <span class="skeleton-block skeleton-path"></span>
            </div>
            <div class="list-meta">
              <div v-for="metaIndex in 4" :key="metaIndex">
                <span class="skeleton-block skeleton-meta-label"></span>
                <span class="skeleton-block skeleton-meta-value"></span>
              </div>
            </div>
          </div>
          <div v-else class="card-body skeleton-card-body">
            <span class="skeleton-block skeleton-file"></span>
            <div class="meta">
              <span class="skeleton-block skeleton-meta-wide"></span>
              <span class="skeleton-block skeleton-meta-code"></span>
            </div>
            <div class="meta">
              <span class="skeleton-block skeleton-meta-mid"></span>
              <span class="skeleton-block skeleton-meta-code"></span>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "HomeSkeleton",
});

defineProps<{
  section: "dates" | "groups";
  viewMode?: "standard" | "compact" | "list";
}>();
</script>

<style scoped>
.date-link-skeleton {
  width: 100%;
  border-radius: 7px;
  padding: 10px;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
  align-items: center;
}

.loading-groups {
  padding-top: 1px;
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

.video-list {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  padding-top: 14px;
}

.video-list.compact {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.video-list.list {
  grid-template-columns: 1fr;
  gap: 10px;
}

.video-card {
  border: 1px solid var(--line);
  background: var(--panel);
  border-radius: var(--radius);
  overflow: hidden;
  min-width: 0;
}

.video-list.list .video-card {
  display: grid;
  grid-template-columns: 128px minmax(0, 1fr);
  min-height: 128px;
}

.skeleton-block,
.skeleton-thumb,
.skeleton-badge,
.skeleton-play,
.skeleton-duration {
  position: relative;
  overflow: hidden;
  background: #e6e2d8;
}

.skeleton-block::after,
.skeleton-thumb::after,
.skeleton-badge::after,
.skeleton-play::after,
.skeleton-duration::after {
  content: "";
  position: absolute;
  inset: 0;
  transform: translateX(-100%);
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.62), transparent);
  animation: skeleton-shimmer 1.3s ease-in-out infinite;
}

.skeleton-block {
  display: block;
  border-radius: 5px;
}

.skeleton-date-title {
  width: min(118px, 72%);
  height: 15px;
}

.skeleton-date-count {
  width: 28px;
  height: 15px;
}

.skeleton-card:hover {
  transform: none;
  box-shadow: none;
  border-color: var(--line);
}

.skeleton-thumb {
  width: 100%;
  aspect-ratio: 16 / 9;
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.24) 1px, transparent 1px) 0 0 / 34px 34px,
    linear-gradient(rgba(255, 255, 255, 0.18) 1px, transparent 1px) 0 0 / 34px 34px,
    #e6e2d8;
}

.video-list.list .skeleton-thumb {
  width: 128px;
  height: 128px;
  aspect-ratio: auto;
}

.skeleton-badge,
.skeleton-play,
.skeleton-duration {
  position: absolute;
  display: block;
  background: rgba(255, 255, 255, 0.72);
}

.skeleton-badge {
  left: 10px;
  top: 10px;
  width: 64px;
  height: 23px;
  border-radius: 5px;
}

.skeleton-play {
  left: 12px;
  bottom: 12px;
  width: 38px;
  height: 38px;
  border-radius: 50%;
}

.skeleton-duration {
  right: 10px;
  bottom: 12px;
  width: 46px;
  height: 24px;
  border-radius: 5px;
}

.card-body {
  padding: 12px;
  display: grid;
  gap: 10px;
}

.skeleton-card-body {
  min-height: 104px;
}

.list-body {
  min-width: 0;
  padding: 12px 14px;
  display: grid;
  grid-template-columns: minmax(220px, 1.2fr) minmax(360px, 1fr);
  gap: 16px;
  align-items: center;
}

.list-primary {
  min-width: 0;
  display: grid;
  gap: 7px;
}

.list-meta {
  min-width: 0;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.list-meta div {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.skeleton-group-title {
  width: min(100%, 360px);
}

.skeleton-group-date {
  width: 148px;
  height: 24px;
}

.skeleton-group-label {
  width: 58px;
  height: 15px;
}

.skeleton-group-meta {
  width: 154px;
  height: 16px;
}

.skeleton-file {
  width: min(100%, 76%);
  height: 18px;
}

.skeleton-path {
  width: min(100%, 92%);
  height: 13px;
}

.skeleton-meta-label {
  width: 42px;
  height: 12px;
}

.skeleton-meta-value {
  width: min(100%, 72px);
  height: 14px;
}

.skeleton-meta-wide,
.skeleton-meta-mid,
.skeleton-meta-code {
  height: 14px;
}

.skeleton-meta-wide {
  width: min(100%, 132px);
}

.skeleton-meta-mid {
  width: min(100%, 86px);
}

.skeleton-meta-code {
  width: 62px;
}

.meta {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
  min-width: 0;
}

@keyframes skeleton-shimmer {
  100% {
    transform: translateX(100%);
  }
}

@media (max-width: 1120px) {
  .video-list,
  .video-list.compact {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .video-list.list {
    grid-template-columns: 1fr;
  }

  .list-body {
    grid-template-columns: minmax(0, 1fr);
    gap: 10px;
  }
}

@media (max-width: 720px) {
  .video-list,
  .video-list.compact {
    grid-template-columns: 1fr;
  }

  .video-list.list .video-card {
    grid-template-columns: 128px minmax(0, 1fr);
    min-height: 128px;
  }

  .list-body {
    padding: 10px;
    gap: 8px;
  }

  .list-meta {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 7px 10px;
  }

  .list-meta div:nth-child(3) {
    display: none;
  }

  .group-head {
    align-items: start;
    flex-direction: column;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
  }
}
</style>
