<script setup lang="ts">
import { ref } from 'vue';
import * as rt from "../../wailsjs/runtime";

const onLogs = ref('')

rt.EventsOn('onLog', (event) => {
  onLogs.value = event.logged_at + " " + event.direction + ":\n" + event.message + "\n\n" + onLogs.value
});

</script>

<template>
  <div class="modal-dialog modal-fullscreen">
    <div class="modal-content">
      <div class="modal-header">
        <h1 class="modal-title fs-5" id="logsModalToggleLabel">Logs</h1>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body">
        <pre class="word-break: break-all">{{onLogs}}</pre>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" data-bs-target="#logsModalToggle2" data-bs-toggle="modal">Close</button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>