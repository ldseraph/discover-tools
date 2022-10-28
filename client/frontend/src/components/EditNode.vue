<template>
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card>
      <q-card-section class="bg-primary text-white">
        <div class="flex flex-nowrap text-h6">设备ID: {{ props.node.uuid }}</div>
      </q-card-section>
      <q-separator />
      <q-card-actions>
        <q-tabs v-model="tab" class="text-teal">
          <q-tab name="dhcp" label="DHCP" />
          <q-tab name="static" label="静态地址" />
        </q-tabs>
      </q-card-actions>

      <q-separator />

      <q-tab-panels v-model="tab" animated>
        <q-tab-panel name="dhcp">
          DHCP分配地址
        </q-tab-panel>

        <q-tab-panel name="static">
          <q-list dense padding>
            <q-item v-for="value, index in ip4" :key="value">
              <q-item-section>
                <div class="flex items-center">
                  <div class="px-2">地址: {{ value.address }}</div>
                  <div class="px-2">子网掩码：{{ value.prefix }}</div>
                  <div class="px-2">网关：{{ value.gateway }}</div>
                  <q-btn class="px-2" color="primary" label="删除" @click="onDelete(index)" />
                </div>
              </q-item-section>
            </q-item>
          </q-list>

          <q-form @submit="onAdd" class="q-gutter-md">
            <q-input filled v-model="addIP4.address" label="ip地址" :rules="[
              val => val ? isIPv4(val) : 'ip地址错误'
            ]" />

            <q-input filled label="子网掩码" type="number" v-model="addIP4.prefix" lazy-rules :rules="[
              val => val !== null && val !== '' || '请输入子网掩码',
              val => val > 0 && val <= 32 || '请输入正确的子网掩码(0-32)'
            ]" />

            <q-input filled v-model="addIP4.gateway" label="网关地址" :rules="[
              val => val ? isIPv4(val) : '网关地址错误'
            ]" />

            <div>
              <q-btn label="添加ip" type="submit" color="primary" />
            </div>
          </q-form>

        </q-tab-panel>

      </q-tab-panels>

      <q-card-actions align="right">
        <q-btn color="primary" label="修改" @click="onOKClick" />
        <q-btn color="primary" label="取消" @click="onDialogCancel" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
``
<script setup lang="ts">
import { isIPv4 } from 'is-ip';
import { useDialogPluginComponent } from 'quasar'
import type { Node, IP4Address } from './node.ts'
import { SetEth } from '@/api/wailsjs/go/main/NodeManager'
// import { main } from '@/api/wailsjs/go/models.ts'

import { ref } from 'vue'

const props = defineProps<{
  node: Node;
}>()

const tab = ref(props.node.dhcp == true ? 'dhcp' : 'static')
const ip4 = ref(props.node.ip4 || [])

function onDelete(index: number) {
  ip4.value.splice(index, 1);
}

const addIP4 = ref<IP4Address>({
  prefix: 24
})

function onAdd() {
  let index = ip4.value.findIndex((element: IP4Address) => element.address == addIP4.value.address)
  if (index >= 0) {
    ip4.value.splice(index, 1);
  }
  addIP4.value.prefix = +addIP4.value.prefix
  ip4.value.push({ ...addIP4.value })
}

defineEmits([
  ...useDialogPluginComponent.emits
])

const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } = useDialogPluginComponent()

function onOKClick() {
  let newvalue = {
    uuid: props.node.uuid,
    ip4: [],
    dhcp: tab.value == 'dhcp' ? true : false,
  }

  ip4.value.forEach(function (node) {
    newvalue.ip4.push(node);
  })

  console.log(newvalue)

  SetEth(newvalue)

  onDialogOK()
}

</script>