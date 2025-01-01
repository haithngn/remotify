<script setup lang="ts">
import { ref, defineProps, onMounted } from 'vue'

const props = defineProps({
  didSelect: Function,
  getItems: Function,
  clearItem: Function,
  editItem: Function
})

const items = ref(<any>[])

onMounted(() => {
  if (props.getItems !== undefined) {
    props.getItems().then((data: any) => {
      items.value = data
    }).catch((err: any) => {
      console.log(err)
    })
  }
})

function click(event : any, item: any) {
  if (event.type === 'click' && event.pointerType !== '') {
    console.log('clicked', item)
    props.didSelect != null ? props.didSelect(item) : null
  }
}

function deleteItem(item: any) {
  props.clearItem != null ? props.clearItem(item) : null
}

function editItem(item: any) {
  props.editItem != null ? props.editItem(item) : null
}

</script>

<template>
    <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton"
            data-bs-toggle="dropdown" aria-expanded="false" :disabled="items.length === 0"></button>
    <ul class="dropdown-menu" aria-labelledby="dropdownMenuItem">
      <li v-for="item in items" v-bind:key="item.value">
        &nbsp;&nbsp; <i class="bi bi-clock text-muted">&nbsp;&nbsp;</i>
        <span v-on:click="click($event, item)" class="timestamp text-muted text-sm">{{item.sent_at}} | </span>
        <span v-on:click="click($event, item)">{{item.value}}</span>  |
        <i class="bi bi-pencil" v-on:click="editItem(item)"> &nbsp; </i>
        <i class="bi bi-x" v-on:click="deleteItem(item)"></i>
        <hr class="dropdown-divider">
      </li>
    </ul>
</template>

<style scoped>

</style>