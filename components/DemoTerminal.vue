<script setup lang="ts">
// A live shell in the demo VM, wrapped in the theme's terminal chrome. It
// lands in the companion code repository via scripts/dev-shell.sh, as the
// ordinary user, so every command in the speaker notes works as typed:
// `sudo` is written where it is needed and the rootless beat needs none.
//
// cmd is a constant, so every terminal in the deck shares one PTY session
// unless it passes its own sid. That is what carries scrollback from one
// slide's terminal into the next.
withDefaults(defineProps<{
  fontSize?: number
  sid?: string
  autoConnect?: boolean
}>(), {
  fontSize: 16,
  autoConnect: true,
})

const cmd = './scripts/dev-shell.sh'
</script>

<template>
  <div class="demo-terminal">
    <TerminalWindow label="crs" command="mycontainer">
      <LiveTerminal :cmd="cmd" :sid="sid" :font-size="fontSize" :auto-connect="autoConnect" />
    </TerminalWindow>
  </div>
</template>

<style scoped>
.demo-terminal {
  width: 100%;
  height: 100%;
}

/* The theme pads .terminal-body for static slot content; the live terminal
   brings its own geometry, so pull the padding back in. */
.demo-terminal :deep(.terminal-body) {
  padding: var(--ac-space-2);
  overflow: hidden;
}
</style>
