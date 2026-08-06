<script setup lang="ts">
import AppContent from './components/AppContent.vue'
import DevicePreviewShell from './components/DevicePreviewShell.vue'

const previewEnabled = import.meta.env.DEV || import.meta.env.VITE_DEVICE_PREVIEW === 'true'
</script>

<template>
  <DevicePreviewShell v-if="previewEnabled">
    <AppContent />
  </DevicePreviewShell>

  <AppContent v-else />
</template>
