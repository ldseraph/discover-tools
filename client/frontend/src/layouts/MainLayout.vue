<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-bar style="--wails-draggable:drag">
        <q-icon name="laptop_chromebook" />
        <div>{{ config.app.name }}</div>

        <q-space />

        <q-btn dense flat icon="minimize" @click="WindowMinimise" />
        <q-btn dense flat icon="crop_square"
          @click="WindowToggleMaximise" />
        <q-btn dense flat icon="close" @click="Quit" />
      </q-bar>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import config from '@/config'
import { WindowMinimise, WindowToggleMaximise, Quit } from '@/api/wailsjs/runtime'
</script>
