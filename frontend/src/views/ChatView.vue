<script setup>
import ChatList from '@/components/ChatList.vue';
import ChatWindow from '@/components/ChatWindow.vue';
</script>
<template>
    <div class="w-full">
        <div class="flex">
            <ChatList/>
            <ChatWindow/>
        </div>
    </div>
</template>