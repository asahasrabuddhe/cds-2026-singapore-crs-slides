<script setup lang="ts">
// One slide per workshop step: the idea on the left, hints on the right,
// and the commands to check it across the full width at the bottom, so a
// long command line never wraps or clips. Based on the theme's
// two-cols-header.
</script>

<template>
  <div class="slidev-layout layout-step">
    <div class="header-strip">
      <h2><slot name="title" /></h2>
    </div>
    <div class="cols-area">
      <div class="col-left">
        <span class="ac-label">THE IDEA</span>
        <slot />
      </div>
      <div class="col-right">
        <span class="ac-label">HINTS</span>
        <slot name="right" />
      </div>
    </div>
    <div class="run-strip">
      <span class="ac-label">CHECK IT</span>
      <slot name="run" />
    </div>
  </div>
</template>

<style scoped>
.layout-step {
  display: flex;
  flex-direction: column;
  padding: 0;
}

.header-strip {
  background: var(--ac-surface);
  padding: var(--ac-space-5) var(--ac-space-8) var(--ac-space-3);
  flex-shrink: 0;
}

.cols-area {
  background: var(--ac-surface-low);
  flex-grow: 1;
  display: flex;
  flex-direction: row;
  min-height: 0;
}

.col-left {
  width: 55%;
  padding: var(--ac-space-4) var(--ac-space-6);
}

.col-right {
  width: 45%;
  background: var(--ac-surface-mid);
  padding: var(--ac-space-4) var(--ac-space-6);
}

.run-strip {
  background: var(--ac-surface);
  padding: var(--ac-space-3) var(--ac-space-6) var(--ac-space-4);
  flex-shrink: 0;
}

.run-strip :deep(pre) {
  margin: var(--ac-space-1) 0 0;
}

.layout-step::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--ac-primary);
  opacity: 0.7;
  z-index: 10;
}
</style>
