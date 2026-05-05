<template>
  <div class="cloud-logo" :style="{ width: size + 'px', height: size + 'px' }">
    <svg viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient :id="id" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="#4facfe" />
          <stop offset="100%" stop-color="#00f2fe" />
        </linearGradient>
      </defs>
      <!-- 云朵主体 -->
      <path
        d="M75 65C66.7 65 60 58.3 60 50C60 41.7 66.7 35 75 35C83.3 35 90 41.7 90 50C90 58.3 83.3 65 75 65ZM42 65C32.1 65 24 56.9 24 47C24 37.1 32.1 29 42 29C44.2 29 46.3 29.4 48.2 30.1C51.9 21.6 60.3 15.8 70.1 15.8C82.8 15.8 93.1 26.1 93.1 38.8C93.1 40.3 93 41.7 92.7 43.1C91.8 42.5 90.7 42.1 89.5 42.1C87.8 42.1 86.3 42.8 85.2 43.9C82.5 38.8 77.1 35.5 70.8 35.5C63.8 35.5 57.9 39.7 55.4 45.7C53.7 43.5 51 42.1 48 42.1C44.5 42.1 41.7 44.9 41.7 48.3V65H42Z"
        :fill="'url(#' + id + ')'"
        class="main-cloud"
      />
      <!-- 核心/状态圆点 -->
      <circle cx="70" cy="55" r="8" fill="white" fill-opacity="0.7" class="core-node" />
    </svg>
  </div>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({
  size: { type: Number, default: 100 },
  id: { type: String, default: 'cloud-grad' }
});
</script>

<style scoped>
.cloud-logo {
  display: inline-block;
  filter: drop-shadow(0 8px 16px rgba(79, 172, 254, 0.2));
}

.core-node {
  animation: breathe 3s ease-in-out infinite;
}

@keyframes breathe {
  0%, 100% { fill-opacity: 0.4; transform: scale(0.95); transform-origin: 70px 55px; }
  50% { fill-opacity: 1; transform: scale(1.05); transform-origin: 70px 55px; }
}

.main-cloud {
  transition: filter 0.3s ease;
}

.cloud-logo:hover .main-cloud {
  filter: brightness(1.1);
}
</style>
