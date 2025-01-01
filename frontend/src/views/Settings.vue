<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import {Ref, reactive, onMounted, watch, ref} from 'vue'
import {
  GetSettings, UploadAPNSToken, UploadServiceAccount,
  SetKeyID, SetTeamID, UploadAPNSCertificate, SetAPNSType, SetDecryptPassword, SetLegacyAPNSEnvironment,
  RevealFileLocation, SetJWTAPNSEnvironment
} from "../../wailsjs/go/main/App"
import {is} from "@babel/types";
import {toast} from "vue3-toastify";
import { notify } from '../components/notification';

const p8 = ref<string>('')
const keyId = ref<string>('')
const teamId = ref<string>('')
const serviceAccount = ref<string>('')
const p12 = ref<string>('')
const decryptPassword = ref<string>('')
const legacyAPNSEnvironment = ref('')
const jwtAPNSEnvironment = ref('')
const sandboxMode = ref('sandbox')
const productionMode = ref('production')
const jwtSandboxMode = ref('sandbox')
const jwtProductionMode = ref('production')
const apnsType = ref('')
const apnsModeToken = ref('token')
const apnsModeLegacy = ref('legacy')

onMounted(() => {
  load()
})

function load() {
  GetSettings().then(data => {
    let settings = JSON.parse(data)
    p8.value = settings.data.p8
    serviceAccount.value = settings.data.service_account
    keyId.value = settings.data.key_id
    teamId.value = settings.data.team_id
    p12.value = settings.data.p12
    decryptPassword.value = settings.data.decrypt_password
    legacyAPNSEnvironment.value = settings.data.apns_mode === 'sandbox' ? 'sandbox' : 'production'
    jwtAPNSEnvironment.value = settings.data.jwt_apns_mode === 'sandbox' ? 'sandbox' : 'production'
    apnsType.value = settings.data.apns_type === 'token' ? 'token' : 'legacy'
  })
}

watch(legacyAPNSEnvironment, (value) => {
  SetLegacyAPNSEnvironment(value).then(() => {})
})

watch(jwtAPNSEnvironment, (value) => {
  SetJWTAPNSEnvironment(value).then(() => {})
})

watch(apnsType, (value) => {
  SetAPNSType(value).then(() => {})
})

watch(decryptPassword, (value) => {
  SetDecryptPassword(value).then(() => {})
})

watch(keyId, (newValue) => {
  SetKeyID(newValue).then(() => {})
})

watch(teamId, (newValue) => {
  SetTeamID(newValue).then(() => {})
})

function uploadAPNSToken() {
  UploadAPNSToken().then(() => load())
}

function uploadServiceAccount() {
  UploadServiceAccount().then(() => load())
}

function revealServiceAccountFileLocation() {
  RevealFileLocation(serviceAccount.value).then(result => {
    let response = JSON.parse(result)
    if (response.code !== 200) {
      notify(response.message, false)
    }
  })
}

function uploadAPNSCertificate() {
  UploadAPNSCertificate().then(() => load())
}

function revealAPNSCertificateFileLocation() {
  RevealFileLocation(p12.value).then(result => {
    let response = JSON.parse(result)
    console.log(response)
    if (response.code !== 200) {
      notify(response.message, false)
    }
  })
}

function revealP8FileLocation() {
  RevealFileLocation(p8.value).then(result => {
    let response = JSON.parse(result)
    console.log(response)
    if (response.code !== 200) {
      notify(response.message, false)
    }
  })
}

</script>

<template>
  <div class="container-fluid">
    <div class="card">
      <h5 class="card-header"><i class="bi bi-fire text-warning" ></i> &nbsp; Firebase Cloud Messaging</h5>
      <div class="card-body">
          <!-- APNS Setting  -->
          <div class="input-group mb-3">
            <button @click="uploadServiceAccount" type="button" class="btn btn-primary">
              <i class="bi bi-upload"></i>
              &nbsp; Service Account
            </button>
            <input type="text" class="form-control" placeholder="service_account.json" disabled v-model="serviceAccount">
            <button @click="revealServiceAccountFileLocation" type="button" class="btn btn-secondary">
              <i class="bi bi-folder"></i> &nbsp; Reveal </button>
          </div>
      </div>
    </div>

    <p/>
    <!-- APNS Setting  -->
    <div class="card">
      <h5 class="card-header">
        <div class="d-flex justify-content-between">
          <span> <i class="bi bi-apple text-primary" ></i>  &nbsp; APNS (Token-base)</span>
          <div>
            <label class="lead text-success">In-use &nbsp;</label>
            <input class="form-check-input" type="radio" name="enableTokenBase" id="enableTokenBase" v-model=apnsType v-bind:value="apnsModeToken">
          </div>
        </div>
      </h5>
      <div class="card-body">
      <div class="input-group mb-3">
        <button @click="uploadAPNSToken" type="button" class="btn btn-primary">
          <i class="bi bi-upload"></i> &nbsp; APNS Token File</button>
        <input type="text" class="form-control" placeholder="key.p8" disabled v-model="p8">
        <button @click="revealP8FileLocation" type="button" class="btn btn-secondary">
        <i class="bi bi-folder"></i> &nbsp; Reveal </button>
      </div>
      <div class="input-group mb-3">
        <span class="input-group-text"><i class="bi bi-key text-primary" ></i></span>
        <input type="text" class="form-control" placeholder="Key ID" v-model="keyId">
        <span class="input-group-text"><i class="bi bi-bag text-primary" ></i></span>
        <input type="text" class="form-control" placeholder="Team ID" v-model="teamId">
        <span class="input-group-text">
            Sandbox &nbsp; &nbsp;
            <input class="form-check-input" type="radio" name="jwtSandboxModeRadio" id="jwtSandboxMode"
                   v-model="jwtAPNSEnvironment" v-bind:value="jwtSandboxMode"
            >
          </span>
        <span class="input-group-text">
            Production &nbsp; &nbsp;
            <input class="form-check-input" type="radio" name="jwtProductionModeRadio" id="jwtProductionMode"
                   v-model=jwtAPNSEnvironment v-bind:value="jwtProductionMode"
            >
          </span>
      </div>
      </div>
    </div>
    <p/>
    <!-- APNS Legacy  -->
    <div class="card">
        <h5 class="card-header">
          <div class="d-flex justify-content-between">
            <span><i class="bi bi-apple text-secondary" ></i> &nbsp; APNS (Legacy)</span>
            <div>
              <label class="lead text-success">In-use &nbsp;</label>
              <input class="form-check-input" type="radio" name="enableLegacy" id="enableLegacy" v-model=apnsType v-bind:value="apnsModeLegacy">
            </div>
          </div>
        </h5>
      <div class="card-body">
        <div class="input-group mb-3">
          <button @click="uploadAPNSCertificate" type="button" class="btn btn-primary">
            <i class="bi bi-upload"></i> &nbsp; APNS Certificate</button>
          <input type="text" class="form-control" placeholder="certificate.p12" disabled v-model="p12">
          <button @click="revealAPNSCertificateFileLocation" type="button" class="btn btn-secondary">
          <i class="bi bi-folder"></i> &nbsp; Reveal </button>
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text"><i class="bi bi-key text-primary" ></i></span>
          <input type="text" class="form-control" placeholder="Decrypt Password" v-model="decryptPassword">
          <span class="input-group-text">
            Sandbox &nbsp; &nbsp;
            <input class="form-check-input" type="radio" name="sandboxModeRadio" id="sandboxMode" v-model="legacyAPNSEnvironment" v-bind:value="sandboxMode">
          </span>
          <span class="input-group-text">
            Production &nbsp; &nbsp;
            <input class="form-check-input" type="radio" name="productionModeRadio" id="productionMode" v-model=legacyAPNSEnvironment v-bind:value="productionMode">
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
input#productionMode {
  transform: scale(1.2);
}
input#sandboxMode {
  transform: scale(1.2);
}
</style>