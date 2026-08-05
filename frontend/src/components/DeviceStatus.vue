<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = withDefaults(defineProps<{ deviceName?: string; batteryLevel: number }>(), {
  deviceName: 'LEVABOX',
})

const now = ref(new Date())
let timer: ReturnType<typeof setInterval> | undefined

const timeText = computed(() => now.value.toLocaleTimeString('zh-CN', {
  hour: '2-digit',
  minute: '2-digit',
  hour12: false,
}))

const safeBatteryLevel = computed(() => Math.min(100, Math.max(0, Math.round(props.batteryLevel))))

onMounted(() => {
  timer = setInterval(() => { now.value = new Date() }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="device-status" aria-label="设备状态">
    <strong>{{ deviceName }}</strong>
    <span class="device-status__divider"></span>
    <div class="battery-status" :aria-label="`电量 ${safeBatteryLevel}%`">
      <span class="battery-status__shell">
        <i class="battery-status__fill" :style="{ width: `${safeBatteryLevel}%` }"></i>
      </span>
      <span class="battery-status__terminal"></span>
      <b>{{ safeBatteryLevel }}%</b>
    </div>
    <time>{{ timeText }}</time>
  </div>
</template>

<style scoped>
.device-status {
  position: absolute;
  z-index: 7;
  top: 24px;
  right: 30px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, .72);
  font-size: .7rem;
  letter-spacing: .06em;
  text-shadow: 0 2px 12px rgba(0, 0, 0, .7);
}

.device-status > strong { color: rgba(255, 255, 255, .9); font-size: .67rem; letter-spacing: .17em; }
.device-status__divider { width: 1px; height: 13px; background: rgba(255, 255, 255, .22); }
.battery-status { display: flex; align-items: center; gap: 6px; }
.battery-status__shell { display: block; width: 23px; height: 11px; padding: 2px; border: 1px solid rgba(255, 255, 255, .62); border-radius: 3px; }
.battery-status__fill { display: block; height: 100%; min-width: 1px; border-radius: 1px; background: rgba(255, 255, 255, .82); }
.battery-status__terminal { width: 2px; height: 5px; margin-left: -5px; border-radius: 0 2px 2px 0; background: rgba(255, 255, 255, .62); }
.battery-status b { min-width: 28px; color: rgba(255, 255, 255, .8); font-size: .68rem; font-weight: 600; letter-spacing: 0; }
.device-status time { min-width: 34px; color: #fff; font-variant-numeric: tabular-nums; }

@media (max-width: 1100px) {
  .device-status { top: 18px; right: 22px; }
}
</style>
