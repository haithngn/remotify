<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import {defineProps, onMounted, onUnmounted, ref, watch} from "vue";
import {SendAPNS, GetRecentAPNSMessages, RemoveRecentAPNSMessage, GetSettings} from "../../wailsjs/go/main/App";
import RecentItem from "../components/RecentItem.vue";
import JsonEditorVue from 'json-editor-vue'
import { notify } from "../components/notification";
import 'vanilla-jsoneditor/themes/jse-theme-dark.css'
import * as rt from "../../wailsjs/runtime";
import {values} from "../../wailsjs/go/models";

const deviceToken = ref('')
const apnsID = ref('')
const expiredAt = ref('')
const collapseID = ref('')
const priority = ref(values.APNSPriority.Immediately)
const pushType = ref(values.APNSPushType.Alert)
const bundleId = ref('')
const payload = ref('')
const sendButtonTittle = ref('Send')
const invalidate = ref(true)
const isSending = ref(false)
const jsonEditorHeight = ref("435px")
const json_editor_theme = ref('jse-theme-light')

// Token Dropdown
const tokenDropdownComponentKey = ref(0);

// Bundle Dropdown
const bundleIdDropdownComponentKey = ref(0);

const forceRerenderDropdown = () => {
  tokenDropdownComponentKey.value += 1;
  bundleIdDropdownComponentKey.value += 1;
};

function updatePayload(newValue: any) {
  payload.value = newValue
}

function reloadRecentMessages() {
  forceRerenderDropdown()
}

const props = defineProps({
  editRecentMessage: Function,
})

defineExpose({
  updatePayload,
  reloadRecentMessages
});

watch(deviceToken,(_) => {
  validate()
})

watch(bundleId,(_) => {
  validate()
})

watch(payload,(_) => {
  validate()
})

onMounted(() => {
  GetRecentAPNSMessages().then(data => {
    let history = JSON.parse(data)
    if (history.data.length > 0) {
      deviceToken.value = history.data[0].device_token
      bundleId.value = history.data[0].bundle_id
      payload.value = history.data[0].payload || values.PayloadTemplate.DefaultAPNS
      apnsID.value = history.data[0].apns_id === null ? '' : history.data[0].apns_id
      collapseID.value = history.data[0].collapse_id === null ? '' : history.data[0].collapse_id
      expiredAt.value = history.data[0].expired_at === null ? '' : history.data[0].expired_at
      priority.value = history.data[0].priority === null ? values.APNSPriority.Immediately : history.data[0].priority
      pushType.value = history.data[0].push_type === null ? values.APNSPushType.Alert : history.data[0].push_type
    } else {
      deviceToken.value = ''
      bundleId.value = ''
      payload.value = values.PayloadTemplate.DefaultAPNS
      apnsID.value = ''
      collapseID.value = ''
      expiredAt.value = ''
      priority.value   = values.APNSPriority.Immediately
      pushType.value = values.APNSPushType.Alert
    }
  })

  resize()

  validate()

  GetSettings().then(data => {
    let settings = JSON.parse(data)
    if (settings.data.theme_mode !== undefined) {
      json_editor_theme.value = settings.data.theme_mode === 'dark' ? 'jse-theme-dark' : 'jse-theme-light'
    }
  })

  //Listen Window resize event
  window.addEventListener('resize', resize);
})

onUnmounted(() => {
      window.removeEventListener('resize', resize);
    }
)

function resize() {
  let wh = window.innerHeight
  // 340px is the height of the header, and now hard coded
  //TODO: Enhance to form to make it fit its parent
  jsonEditorHeight.value = (wh - 340) + "px"
}

function sendAndSave() {
  isSending.value = true
  sendButtonTittle.value = 'Sending...'
  SendAPNS(deviceToken.value,
      bundleId.value,
      payload.value,
      apnsID.value,
      collapseID.value,
      expiredAt.value,
      priority.value,
      pushType.value,
      true
  ).then(result => {
    let response = JSON.parse(result)
    if (response.code !== 200) {
      notify('Oops! Something went wrong, ' + response.error.message, false)
      sendButtonTittle.value = 'Send'
      isSending.value = false
      return
    }
    notify('Notification has been sent to the device.')
    sendButtonTittle.value = 'Sent!'
    setTimeout(function() {
      isSending.value = false
      sendButtonTittle.value = 'Send'
    }, 3000);
    forceRerenderDropdown()
  })
}
function send() {
  isSending.value = true
  sendButtonTittle.value = 'Sending...'
  SendAPNS(deviceToken.value,
      bundleId.value,
      payload.value,
      apnsID.value,
      collapseID.value,
      expiredAt.value,
      priority.value,
      pushType.value,
      false
  ).then(result => {
    let response = JSON.parse(result)
    if (response.code !== 200) {
      notify('Oops! Something went wrong, ' + response.error.message, false)
      sendButtonTittle.value = 'Send'
      isSending.value = false
      return
    }
    notify('Notification has been sent to the device.')
    sendButtonTittle.value = 'Sent!'
    setTimeout(function() {
      isSending.value = false
      sendButtonTittle.value = 'Send'
    }, 3000);
    forceRerenderDropdown()
  })
}

function clear () {
  deviceToken.value = ''
  bundleId.value = ''
  payload.value = ''
}

function validate() {
  invalidate.value = (deviceToken.value.length == 0 || bundleId.value.length === 0 || payload.value.length === 0)
}

function selectRecentMessage(message: any) {
  deviceToken.value = message.device_token
  bundleId.value = message.bundle_id
  payload.value = message.payload === null ? values.PayloadTemplate.DefaultAPNS : message.payload
  apnsID.value = message.apns_id === null ? '' : message.apns_id
  collapseID.value = message.collapse_id === null ? '' : message.collapse_id
  expiredAt.value = message.expired_at === null ? '' : message.expired_at
  priority.value = message.priority === null ? values.APNSPriority.Immediately : message.priority
  pushType.value = message.push_type === null ? values.APNSPushType.Alert : message.push_type
  validate()
}

async function getRecentTokens() {
  let data = await GetRecentAPNSMessages()
  let history = JSON.parse(data)
  let result = []
  if (history.data.length > 0) {
    result = history.data.map((e: any) => {
      return {
        id: e.id,
        device_token: e.device_token,
        bundle_id: e.bundle_id,
        payload: e.payload,
        apns_id: e.apns_id,
        collapse_id: e.collapse_id,
        expired_at: e.expired_at,
        priority: e.priority,
        push_type: e.push_type,
        value: e.note,
        sent_at: e.created_at
      }
    })
  }

  return result
}

async function getRecentBundleIDs() {
  let data = await GetRecentAPNSMessages()
  let history = JSON.parse(data)
  let result = []
  if (history.data.length > 0) {
    result = history.data.map((e: any) => {
      return {
        id: e.id,
        device_token: e.device_token,
        bundle_id: e.bundle_id,
        payload: e.payload,
        apns_id: e.apns_id,
        collapse_id: e.collapse_id,
        expired_at: e.expired_at,
        priority: e.priority,
        push_type: e.push_type,
        value: e.note,
        sent_at: e.created_at
      }
    })
  }

  return result
}

async function clearRecentItems(item: any) {
  let id = item.id
  if (id !== undefined) {
    try {
      await RemoveRecentAPNSMessage(id)
      forceRerenderDropdown()
    } catch (e) {
      console.log(e)
    }
  }
}

async function editRecentItem(item: any) {
  props.editRecentMessage != null ? props.editRecentMessage(item) : null
}

rt.EventsOn('onChangeDarkMode', (event) => {
  json_editor_theme.value = 'jse-theme-dark'
});

rt.EventsOn('onChangeLightMode', (event) => {
  json_editor_theme.value = 'jse-theme-light'
});

//Cre: https://stackoverflow.com/a/2117523
function uuidv4() {
  return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c =>
      (+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4).toString(16)
  );
}
</script>

<template>
  <div class="container-fluid">

    <form>
      <div class="input-group mb-3">
        <span class="input-group-text">Device Token</span>
        <input class="form-control" type="text" placeholder="device token" v-model="deviceToken">
        <RecentItem :didSelect="selectRecentMessage"
                    :getItems="getRecentTokens"
                    :key="tokenDropdownComponentKey"
                    :editItem="editRecentItem"
        :clearItem="clearRecentItems"/>
        <span class="input-group-text">Bundle ID</span>
        <input class="form-control" type="text" placeholder="com.example.mobile" v-model="bundleId">
        <RecentItem :didSelect="selectRecentMessage"
                    :getItems="getRecentBundleIDs"
                    :key="bundleIdDropdownComponentKey"
                    :clearItem="clearRecentItems"
                    :editItem="editRecentItem"
        />
      </div>

      <div class="input-group mb-3">
        <span class="input-group-text">APNS ID</span>
        <input class="form-control" type="text" placeholder="APNS Id" v-model="apnsID">
        <span class="input-group-text" @click="() => apnsID = uuidv4()">
          <i class="bi bi-braces" ></i>
        </span>
        <span class="input-group-text" @click="() => apnsID = ''">
          <i :class="apnsID != '' ? 'bi bi-trash text-danger' : 'bi bi-trash text-secondary'" ></i>
        </span>
        <span class="input-group-text">Collapse ID</span>
        <input class="form-control" type="text" placeholder="Collapse ID" v-model="collapseID">
        <span class="input-group-text" @click="() => collapseID = ''">
          <i :class="collapseID != '' ? 'bi bi-trash text-danger' : 'bi bi-trash text-secondary'" ></i>
        </span>
      </div>
      <div class="input-group mb-3">
        <span class="input-group-text">Expired At</span>
        <input class="form-control" type="datetime-local" placeholder="Expired At" v-model="expiredAt">
        <span class="input-group-text" @click="() => expiredAt = ''">
          <i :class="expiredAt != '' ? 'bi bi-trash text-danger' : 'bi bi-trash text-secondary'" ></i>
        </span>

        <span class="input-group-text">Priority</span>
        <input class="form-control" type="text" placeholder="Priority" v-model="priority" disabled>
        <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButtonPriority"
                data-bs-toggle="dropdown" aria-expanded="false"></button>
        <ul class="dropdown-menu" aria-labelledby="dropdownMenuButtonPushTypeList">
          <li>
            <a @click="priority = 10" class="dropdown-item" href="#">10</a>
            <a @click="priority = 5" class="dropdown-item" href="#">5</a>
          </li>
        </ul>
        <span class="input-group-text">Push Type</span>
        <input class="form-control" type="text" placeholder="Push Type" v-model="pushType" disabled>
        <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButtonPriority"
                data-bs-toggle="dropdown" aria-expanded="false"></button>
        <ul class="dropdown-menu" aria-labelledby="dropdownMenuButtonPushTypeList">
          <li>
            <a @click="pushType = values.APNSPushType.Alert" class="dropdown-item" href="#">alert</a>
            <a @click="pushType = values.APNSPushType.Background" class="dropdown-item" href="#">background</a>
            <a @click="pushType = values.APNSPushType.Location" class="dropdown-item" href="#">location</a>
            <a @click="pushType = values.APNSPushType.VoIP" class="dropdown-item" href="#">voip</a>
            <a @click="pushType = values.APNSPushType.Complication" class="dropdown-item" href="#">complication</a>
            <a @click="pushType = values.APNSPushType.FileProvider" class="dropdown-item" href="#">fileprovider</a>
            <a @click="pushType = values.APNSPushType.MDM" class="dropdown-item" href="#">mdm</a>
          </li>
        </ul>
      </div>
        <label>Payload</label>
      <div>
        <JsonEditorVue :class="json_editor_theme" :style="{height: jsonEditorHeight}"
            v-model="payload"
        />
      </div>
      <p/>
      <div class="mb-4 d-flex justify-content-between">
        <div class="d-flex justify-content-between">
          <button @click="clear" type="button" class="btn btn-danger" :disabled="isSending"> Reset <i class="bi bi-eraser text-white" ></i></button>
          &nbsp;&nbsp;
          <button @click="sendAndSave" type="button" class="btn btn-primary" :disabled="invalidate||isSending">
            Push & Save <i class="bi bi-send text-white" ></i>
          </button>
        </div>
        <div class="d-flex justify-content-between">
          <button @click="send" type="button" class="btn btn-primary" :disabled="invalidate||isSending">
            Push
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<style scoped>
textarea {
  resize: none;
  display: inline-block;
}

</style>