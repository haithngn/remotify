<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import {defineProps, onMounted, onUnmounted, ref, watch} from "vue";
import {GetRecentFCMMessages, GetSettings, RemoveRecentFCMMessage, SendFCM} from "../../wailsjs/go/main/App";
import RecentItem from "../components/RecentItem.vue";
import JsonEditorVue from "json-editor-vue";
import { notify } from "../components/notification";
import 'vanilla-jsoneditor/themes/jse-theme-dark.css'
import * as rt from "../../wailsjs/runtime";
import {values} from "../../wailsjs/go/models";

const deviceToken = ref('')
const payloadData = ref('')
const sendButtonTittle = ref('Send')
const invalidate = ref(true)
const isSending = ref(false)
const jsonEditorHeight = ref("480px")
const json_editor_theme = ref('jse-theme-light')
const deviceType = ref(values.FCMDeviceType.Android)

// Bundle Dropdown
const tokenDropdownComponentKey = ref(0);

const forceRerenderDropdowns = () => {
  tokenDropdownComponentKey.value += 1;
};

const props = defineProps({
  editRecentMessage: Function,
})


watch(deviceToken,(_) => {
  validate()
})

watch(payloadData,(_) => {
  validate()
})

onMounted(() => {
  GetRecentFCMMessages().then(data => {
    let history = JSON.parse(data)
    if (history.data.length > 0) {
      deviceToken.value = history.data[0].device_token
      payloadData.value = history.data[0].payload_data
      deviceType.value = history.data[0].device_type
    } else {
      deviceToken.value = ''
      payloadData.value = values.PayloadTemplate.DefaultFCMAndroid
      deviceType.value = values.FCMDeviceType.Android
    }
  })
  resize();

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
})

function updatePayload(newValue: any) {
  payloadData.value = newValue
}

function getCurrentDeviceType() {
  return deviceType.value
}

function reloadRecentMessages() {
  forceRerenderDropdowns()
}

defineExpose({
  updatePayload,
  getCurrentDeviceType,
  reloadRecentMessages
});

function resize() {
  let wh = window.innerHeight
  //TODO: Enhance to form to make it fit its parent
  jsonEditorHeight.value = (wh - 266) + "px"
}

function sendAndSave() {
  isSending.value = true
  sendButtonTittle.value = 'Sending...'
  SendFCM(deviceToken.value, deviceType.value, payloadData.value, true).then(result => {
    let response = JSON.parse(result)
    if (response.code !== 200) {
      notify(response.error.message, false)
      isSending.value = false
      sendButtonTittle.value = 'Send'
      return
    }
    notify('Notification has been sent to the device.')
    sendButtonTittle.value = 'Sent!'
    setTimeout(function() {
      isSending.value = false
      sendButtonTittle.value = 'Send'
    }, 3000);
    forceRerenderDropdowns()
  })
}

function send() {
  isSending.value = true
  sendButtonTittle.value = 'Sending...'
  SendFCM(deviceToken.value, deviceType.value, payloadData.value, false).then(result => {
    let response = JSON.parse(result)
    if (response.code !== 200) {
      notify(response.error.message, false)
      isSending.value = false
      sendButtonTittle.value = 'Send'
      return
    }
    notify('Notification has been sent to the device.')
    sendButtonTittle.value = 'Sent!'
    setTimeout(function() {
      isSending.value = false
      sendButtonTittle.value = 'Send'
    }, 3000);
    forceRerenderDropdowns()
  })
}

function clear () {
  deviceToken.value = ''
  payloadData.value = ''
}

function validate() {
  invalidate.value = (deviceToken.value.length == 0 || payloadData.value.length === 0)
}

function selectMessage(message: any) {
  deviceToken.value = message.device_token
  payloadData.value = message.payload_data
  deviceType.value = message.device_type
}

async function getRecentTokens() {
  let data = await GetRecentFCMMessages()
  let history = JSON.parse(data)
  let result = []
  if (history.data.length > 0) {
    result = history.data.map((e: any) => {
      return {
        id: e.id,
        device_type: e.device_type,
        device_token: e.device_token,
        payload_data: e.payload_data,
        value: e.note,
        sent_at: e.created_at
      }
    })
  }

  return result
}

async function editRecentItem(item: any) {
  props.editRecentMessage != null ? props.editRecentMessage(item) : null
}

async function clearRecentItems(item: any) {
  let id = item.id
  if (id !== undefined) {
    try {
      await RemoveRecentFCMMessage(id)
      forceRerenderDropdowns()
    } catch (e) {
      console.log(e)
    }
  }
}

rt.EventsOn('onChangeDarkMode', (event) => {
  json_editor_theme.value = 'jse-theme-dark'
});

rt.EventsOn('onChangeLightMode', (event) => {
  json_editor_theme.value = 'jse-theme-light'
});
</script>

<template>
  <div class="container-fluid">

    <form>
      <div class="input-group mb-3">
        <span class="input-group-text">Device Token</span>
        <input class="form-control" type="text" placeholder="device token" v-model="deviceToken">
        <RecentItem
            :didSelect="selectMessage"
            :getItems="getRecentTokens"
            :key="tokenDropdownComponentKey"
            :editItem="editRecentItem"
            :clearItem="clearRecentItems"
        />
      </div>
      <div class="input-group mb-3">
        <span class="input-group-text">
          <i class="bi bi-bag text-primary" ></i> &nbsp; Device
        </span>
        <span class="input-group-text">
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="inlineRadioOptions" id="deviceAndroidId"
                   v-model="deviceType"
                   :value="values.FCMDeviceType.Android" 
                   :checked="deviceType === values.FCMDeviceType.Android"
            >
            <label class="form-check-label" for="deviceAndroidId">Android</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="inlineRadioOptions" id="deviceIOSId"
                   v-model="deviceType"
                   :value="values.FCMDeviceType.iOS" 
                   :checked="deviceType === values.FCMDeviceType.iOS"
            >
            <label class="form-check-label" for="deviceIOSId">iOS</label>
          </div>
          <div class="form-check form-check-inline">
            <input class="form-check-input" type="radio" name="inlineRadioOptions" id="deviceWebId"
                   v-model="deviceType"
                   :value="values.FCMDeviceType.Web" 
                   :checked="deviceType === values.FCMDeviceType.Web"
            >
            <label class="form-check-label" for="deviceWebId">Web</label>
          </div>
        </span>
      </div>
      <div>
        <JsonEditorVue :class="json_editor_theme"
                       :style="{height: jsonEditorHeight}"
                       v-model="payloadData"
        />
      </div>
      <p/>
      <div class="mb-4 d-flex justify-content-between">
        <div class="d-flex justify-content-between">
          <button @click="clear" type="button" class="btn btn-danger" :disabled="isSending"> Reset <i class="bi bi-eraser text-white" ></i></button>&nbsp;&nbsp;&nbsp;&nbsp;
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