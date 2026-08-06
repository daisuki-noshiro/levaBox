<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{ batteryLevel: number }>()

const now = ref(new Date())
let clockTimer: ReturnType<typeof setInterval> | undefined

const safeBatteryLevel = computed(() => Math.min(100, Math.max(0, Math.round(props.batteryLevel))))
const timeText = computed(() => now.value.toLocaleTimeString('zh-CN', {
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
}))

onMounted(() => {
  clockTimer = setInterval(() => { now.value = new Date() }, 1000)
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
})
</script>

<template>
  <div class="device-status" aria-label="设备状态">
    <time>{{ timeText }}</time>
    <div class="battery-status" :aria-label="`电量 ${safeBatteryLevel}%`">
      <span class="battery-status__shell">
        <i class="battery-status__fill" :style="{ width: `${safeBatteryLevel}%` }"></i>
      </span>
      <span class="battery-status__terminal"></span>
      <b>{{ safeBatteryLevel }}%</b>
    </div>
  </div>
</template>

<style scoped>
.device-status {
  position: absolute;
  z-index: 7;
  top: clamp(20px, 2.2cqh, 35px);
  left: clamp(22px, 2.2cqw, 46px);
  display: flex;
  align-items: center;
  gap: clamp(13px, 1cqw, 20px);
  color: rgba(255, 255, 255, .9);
  text-shadow: 0 2px 12px rgba(0, 0, 0, .7);
}

.device-status time { color: #fff; font-size: clamp(.8rem, .76cqw, .96rem); font-weight: 700; font-variant-numeric: tabular-nums; letter-spacing: .02em; }
.battery-status { display: flex; align-items: center; gap: clamp(7px, .5cqw, 10px); }
.battery-status__shell { display: block; width: clamp(24px, 1.55cqw, 31px); height: clamp(12px, .78cqw, 16px); padding: 2px; border: 1px solid rgba(255, 255, 255, .72); border-radius: 3px; }
.battery-status__fill { display: block; height: 100%; min-width: 1px; border-radius: 1px; background: rgba(255, 255, 255, .82); }
.battery-status__terminal { width: 3px; height: clamp(6px, .4cqw, 8px); margin-left: calc(clamp(7px, .5cqw, 10px) * -1 - 1px); border-radius: 0 2px 2px 0; background: rgba(255, 255, 255, .7); }
.battery-status b { min-width: 36px; color: #fff; font-size: clamp(.78rem, .72cqw, .9rem); font-weight: 700; font-variant-numeric: tabular-nums; letter-spacing: .01em; }

@container levabox-app (max-width: 1100px) {
  .device-status { top: 17px; left: 18px; }
}
</style>
