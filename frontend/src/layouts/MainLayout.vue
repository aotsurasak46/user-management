<script setup>
import Sidebar from '@/components/Sidebar.vue';
import { useChatStore } from '../stores/chat';
import { onMounted, onUnmounted } from 'vue';
const chatStore = useChatStore()
onMounted(() => {
    chatStore.connect()
})

onUnmounted(() => {
    chatStore.closeConnection()
})
</script>
<template>
    <div class="flex h-screen">
        <Sidebar/>
        <main class="flex-1 bg-gray-50">
            <RouterView />
        </main>
    </div>
</template>