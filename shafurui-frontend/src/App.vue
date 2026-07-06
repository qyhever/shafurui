<template>
  <div id="app">
    <router-view />
    <VersionUpdatePopup />
  </div>
</template>

<script setup lang="ts">
import VersionUpdatePopup from '@/components/VersionUpdatePopup.vue'
defineOptions({
  name: 'App',
})
</script>
