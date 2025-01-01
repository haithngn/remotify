<script lang="ts" setup>
import {reactive, onMounted, ref} from 'vue'
import APNS from "./APNS.vue";
import FCM from "./FCM.vue";
import Settings from "./Settings.vue";
import * as rt from "../../wailsjs/runtime/runtime.js";
import { Modal } from 'bootstrap';
import {
  GetSettings,
  ToggleDarkTheme,
  ToggleLightTheme,
  UpdateAPNSMessageNote,
  UpdateFCMMessageNote
} from "../../wailsjs/go/main/App"
import Logging from "../components/Logging.vue";
import { values } from '../../wailsjs/go/models';
import { notify } from '../components/notification';
import {Platform, RemoteNotificationType} from "../common/enums";
import {toast} from "vue3-toastify";

const settingComponentKey = ref(0);
const apnsComponentKey = ref(0);
const fcmComponentKey = ref(0);
const apnsRef = ref()
const fcmRef = ref()

const currentEditMessageId = ref(0)
const currentEditMessageType = ref(RemoteNotificationType.APNS)
const currentEditNote = ref('')
const isDarkMode = ref(false)

onMounted(() => {
  load()
})

function load() {
  GetSettings().then(data => {
    let settings = JSON.parse(data)
    isDarkMode.value = settings.data.theme_mode === 'dark'
    setTheme(isDarkMode.value ? 'dark' : 'light')
  })
}

const setTheme = (theme: string | null) => {
  if (theme === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-bs-theme', 'dark')
  } else {
    isDarkMode.value = false
    document.documentElement.setAttribute('data-bs-theme', theme || 'light')
  }
}

const toggleTheme = (theme: string) => {
  if (theme === 'dark') {
    setTheme('dark')
  } else {
    setTheme('light')
  }
}

async function toggleLightTheme() {
  await ToggleLightTheme();
}

async function toggleDarkTheme() {
  await ToggleDarkTheme();
}

const showLogs = () => {
  let m = document.getElementById('logsModalToggle')
  let modal = new Modal(m!, {});
  modal.show();
}

// Events
rt.EventsOn('onOpenSettingWindow', (event) => {
  const btn = document.getElementById("nav-setting-tab");
  btn?.click()

  // force re-render Settings
  settingComponentKey.value++
});

rt.EventsOn('onOpenAboutWindow', (event) => {
  let m = document.getElementById('modal-about')
  let myModal = new Modal(m!, {});
  myModal.show();
});

rt.EventsOn('onOpenAPNS', (event) => {
  const btn = document.getElementById("nav-apns-tab");
  btn?.click()
});

rt.EventsOn('onOpenFCM', (event) => {
  const btn = document.getElementById("nav-fcm-tab");
  btn?.click()
});

rt.EventsOn('onChangeDarkMode', (event) => {
  toggleTheme('dark');
});

rt.EventsOn('onChangeLightMode', (event) => {
  toggleTheme('light');
});

rt.EventsOn('onOpenLogs', (event) => {
  showLogs();
});

rt.EventsOn('onResetPayload', (event) => {
  // Find Active Tab
  const activeTab = document.querySelector('.nav-tabs .nav-link.active')
  if (activeTab) {
    if (activeTab.getAttribute('id') === 'nav-apns-tab') {
      apnsRef.value.updatePayload(values.PayloadTemplate.DefaultAPNS)
    } else if (activeTab.getAttribute('id') === 'nav-fcm-tab') {
      let deviceType = fcmRef.value.getCurrentDeviceType()
      switch (deviceType) {
        case values.FCMDeviceType.Android:
          fcmRef.value.updatePayload(values.PayloadTemplate.DefaultFCMAndroid)
          break
        case values.FCMDeviceType.iOS:
          fcmRef.value.updatePayload(values.PayloadTemplate.DefaultFCMiOS)
          break
        case values.FCMDeviceType.Web:
          fcmRef.value.updatePayload(values.PayloadTemplate.DefaultFCMWeb)
          break
      }
    } else {
      notify('You must select APNS or FCM tab to reset payload', false)
    }
  }
});

async function editRecentAPNSMessage(item: any) {
  currentEditMessageId.value = item.id
  currentEditMessageType.value = RemoteNotificationType.APNS
  currentEditNote.value = item.value

  let m = document.getElementById('editNoteModal')
  let modal = new Modal(m!, {});
  modal.show();
}

async function editRecentFCMMessage(item: any) {
  currentEditMessageId.value = item.id
  currentEditMessageType.value = RemoteNotificationType.FCM
  currentEditNote.value = item.value

  let m = document.getElementById('editNoteModal')
  let modal = new Modal(m!, {});
  modal.show();
}

async function updateNote() {
  switch (currentEditMessageType.value) {
    case RemoteNotificationType.APNS:
      await UpdateAPNSMessageNote(currentEditNote.value, currentEditMessageId.value).then(() => {
        notify("Recent message has been updated", true)
        apnsComponentKey.value++
        apnsRef.value.reloadRecentMessages()
      })
      break
    case RemoteNotificationType.FCM:
      await UpdateFCMMessageNote(currentEditNote.value, currentEditMessageId.value).then(() => {
        notify("Recent message has been updated", true)
        fcmComponentKey.value++
        fcmRef.value.reloadRecentMessages()
      })
      break
  }
}

</script>

<template>
  <main>
    <!--    Logs-->
    <div class="modal fade" id="logsModalToggle" aria-hidden="true" aria-labelledby="logsModalToggleLabel" tabindex="-1">
      <Logging/>
    </div>

    <!-- Edit Note -->
    <div class="modal fade" id="editNoteModal" tabindex="-1" aria-labelledby="editNoteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h1 class="modal-title fs-5" id="editNoteModalLabel">Update Recent Message</h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form onsubmit="return false">
              <div class="mb-3">
                <label for="message-note" class="col-form-label">Recipient:</label>
                <input type="text" class="form-control" id="message-note" v-model="currentEditNote">
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Discard</button>
            <button type="button" class="btn btn-primary" v-on:click="updateNote" data-bs-dismiss="modal">Update</button>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="modal-about" data-bs-backdrop="true" data-bs-keyboard="true" tabindex="-1" aria-labelledby="modal-aboutLabel" aria-hidden="false">
      <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content">
          <div class="modal-body">
            <div class=" d-flex align-items-center justify-content-center">
              <img src="../assets/images/icon.svg" class=".img-fluid about-logo" alt="Remotify App"/>
            </div>
            <label  class="text-sm-center">Version {{ values.Bundle.Version }}</label>
            <label  class="text-sm-center">Build with <a href="https://wails.io">Wails</a> and <a href="https://vuejs.org">Vue</a></label>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
    <div class="container-fluid">
      <div class="d-flex justify-content-between align-items-baseline">
        <div class="d-flex justify-content-between">
          <nav class="navbar">
            <div class="nav nav-tabs" id="nav-tab" role="tablist">
              <button class="nav-link active" id="nav-apns-tab" data-bs-toggle="tab" data-bs-target="#nav-home" type="button" role="tab" aria-controls="nav-home" aria-selected="true"><i class="bi bi-apple text-primary" ></i>  &nbsp;APNS</button>
              <button class="nav-link" id="nav-fcm-tab" data-bs-toggle="tab" data-bs-target="#nav-profile" type="button" role="tab" aria-controls="nav-profile" aria-selected="false"><i class="bi bi-fire text-warning" ></i>  &nbsp;FCM</button>
              <button class="nav-link" id="nav-setting-tab" data-bs-toggle="tab" data-bs-target="#nav-contact" type="button" role="tab" aria-controls="nav-contact" aria-selected="false"><i class="bi bi-gear-wide-connected text-primary" ></i>  &nbsp;Settings</button>
            </div>
          </nav>
        </div>
        <div class="d-flex justify-content-between">
<!--          Theme Switch -->
          <div class="form-check form-switch">
            <div class="btn">
              <i class="bi bi-brightness-high-fill" v-on:click="toggleLightTheme" v-if="isDarkMode"></i>
              <i class="bi bi-moon-fill" v-on:click="toggleDarkTheme" v-else></i>
            </div>
          </div>

        </div>
      </div>
      <p/>
      <div class="tab-content" id="nav-tabContent">
        <div class="tab-pane fade show active" id="nav-home" role="tabpanel" aria-labelledby="nav-apns-tab" tabindex="0">
          <APNS ref="apnsRef" :edit-recent-message="editRecentAPNSMessage"/>
        </div>
        <div class="tab-pane fade" id="nav-profile" role="tabpanel" aria-labelledby="nav-fcm-tab" tabindex="1">
          <FCM ref="fcmRef" :edit-recent-message="editRecentFCMMessage"/>
        </div>
        <div class="tab-pane fade" id="nav-contact" role="tabpanel" aria-labelledby="nav-setting-tab" tabindex="2">
          <Settings :key="settingComponentKey"/>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.about-logo {
  object-position: center;
  object-fit: cover;
  height: 100px;
  width: 100px;
}

.remotify-tab-icon {
  object-position: center;
  object-fit: cover;
  height: 16px;
  width: 16px;
}
</style>
