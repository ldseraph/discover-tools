<template>
  <q-page class="q-pa-md">
    <div class="row">
      <div class="col">
        <q-table title="设备列表" :rows="rows" :columns="columns" row-key="name">
          <template v-slot:body-cell-dhcp="props">
            <q-td :props="props">
              <div class="flex flex-nowrap items-center">
                <div class="pr-4"> {{ props.value ? '是' : '否' }}</div>
              </div>
            </q-td>
          </template>

          <template v-slot:body-cell-ip4="props">
            <q-td :props="props">
              <div class="flex flex-nowrap items-center">
                <div class="pr-4">
                  <q-list dense padding>
                    <q-item v-for="value in props.value" :key="value">
                      <q-item-section>
                        <div class="flex">
                          <div class="px-2">地址: {{value.address}}</div>
                          <div class="px-2">子网掩码：{{value.prefix}}</div>
                          <div class="px-2">网关：{{value.gateway}}</div>
                        </div>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </div>
              </div>
            </q-td>
          </template>

          <!-- <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <div class="flex flex-nowrap items-center">
                <div v-if="props.value == 'online'" class="w-3 h-3 rounded-full bg-green-500"></div>
                <div v-else-if="props.value == 'offline'" class="w-3 h-3 rounded-full bg-red-500"></div>
                <div class="pl-2">
                  {{ props.value == 'online'? '在线': '离线' }}
                </div>
              </div>
            </q-td>
          </template> -->

          <template v-slot:body-cell-edit="props">
            <q-td :props="props">
              <div class="flex flex-nowrap items-center">
                <q-btn flat color="primary" label="修改" @click="edit(props.row)" />
              </div>
            </q-td>
          </template>
        </q-table>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar'
import { ref } from 'vue'
import { EventsOn } from '@/api/wailsjs/runtime'
import type { Node } from '@/components/node.ts'
import EditNode from '@/components/EditNode.vue'

const $q = useQuasar()

EventsOn('node', function (nodes: Node[]) {
  rows.value = []
  nodes && nodes.forEach(function (node) {
    rows.value.push(node);
  })
})

interface Columns {
  name: string;
  label: string;
  field: string;
  required?: boolean;
  align?: 'left' | 'right' | 'center';
  sortable?: boolean;
  sortOrder?: 'ad' | 'da';
  headerStyle?: string;
  headerClasses?: string;
}

const columns: Columns[] = [
  {
    name: 'machine-id',
    required: true,
    label: '设备ID',
    align: 'left',
    field: 'uuid',
    sortable: true
  },
  { name: 'dhcp', align: 'left', label: 'DHCP', field: 'dhcp' },
  { name: 'ip4', align: 'left', label: 'ip地址', field: 'ip4' },
  { name: 'edit', align: 'left', label: '', field: 'edit' }
]

const rows = ref<Node[]>([
  // {
  //   uuid: 'aaaaaa',
  //   ip4: [{
  //     address: '192.168.101.3',
  //     prefix: 24,
  //     gateway: '192.168.101.1',
  //   }],
  //   dhcp: true,
  // }
])

function edit(row:Node) {
  $q.dialog({
    component: EditNode,
    componentProps: {
      node: row,
    }
  })
}

</script>
