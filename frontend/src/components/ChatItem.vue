 <script setup>

import { defineProps } from 'vue';

const props = defineProps({
  userName: { type: String, required: true },
  lastMessage: { type: String, required: true },
  timestamp: { type: String, required: true },
  isSelected: {type: Boolean, required: true  }
})


</script>

<template>
    <div class="flex cursor-pointer mb-2 p-3 justify-between overflow-hidden" :class="[isSelected ? 'bg-gray-200':'']">
        <div class="flex">
            <div class="flex min-w-12 h-12 rounded-full mr-3  items-center justify-center text-black" :class="[isSelected ? 'bg-white':'bg-gray-200']">
                <span>{{ props.userName.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="chat-info">
                <div class="text-md font-bold">{{ props.userName }}</div>
                <div class="text-sm text-gray-500 truncate overflow-hidden text-ellipsis whitespace-nowrap max-w-30">
                    {{ props.lastMessage }}
                </div>
            </div>
        </div>
        <div class="meta">
            <span class="timestamp text-sm">{{ props.timestamp }}</span>
        </div>
    </div>
</template>