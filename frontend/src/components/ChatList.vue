<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useChatStore } from '@/stores/chat';
import ChatItem from './ChatItem.vue';
const chatStore = useChatStore()
const isLoading = ref(false)
const errorMessage = ref('')
const searchValue = ref('')

const filteredChatList = computed(() => {
    if(!chatStore.chatList){
        return []
    }
    let filtered = [...chatStore.chatList]

    if (searchValue.value.trim()) {
        const search = searchValue.value.toLowerCase()
        filtered = filtered.filter(chat =>
        chat.user.name.toLowerCase().includes(search) || chat.user.email.toLowerCase().includes(search)
        )
    }
    return filtered
})


onMounted(async() => {
    fetchChats()
})

watch(chatStore.chatList, () => {
    fetchChats()
})

const fetchChats = async () => {
  try {
    isLoading.value = true
    await chatStore.getChatList()
  } catch (error) {
    errorMessage.value = error.message
    console.log('Failed to load chat history : ', error)
  }finally{
    isLoading.value = false
  }
}

const selectChatUser = (user) => {
    chatStore.setSelectedChatUser(user)
}

const displayTimeStamp = (timestamp) => {
    const date = new Date(timestamp);
    const now = new Date();

    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const isToday = date.toDateString() === today.toDateString();
    if (isToday) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }

    const oneWeekAgo = new Date(today);
    oneWeekAgo.setDate(today.getDate() - 7);

    if (date > oneWeekAgo) {
        return date.toLocaleDateString([], { weekday: 'short' }); 
    }

    if (date.getFullYear() === now.getFullYear()) {
        return date.toLocaleDateString([], { month: 'short', day: 'numeric' }); 
    }

    return date.toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' }); 
};
</script>

<template>
    <div class="w-70 sm:w-100 border-r border-gray-300 h-screen flex flex-col">
        <div class="p-4 ">
            <div class="h-12 md:hidden"></div>
            <div class="flex items-center border px-3 py-1 rounded-md w-full lg:w-full ">
                <input type="text" placeholder="Search" v-model="searchValue" class="w-full outline-none">
                <i class="pi pi-search ml-2 text-gray-500"></i>
            </div>
        </div>
        <div class="flex-1 overflow-y-auto 
            [&::-webkit-scrollbar]:w-2
            [&::-webkit-scrollbar-track]:rounded-full
            [&::-webkit-scrollbar-track]:bg-gray-100
            [&::-webkit-scrollbar-thumb]:rounded-full
            [&::-webkit-scrollbar-thumb]:bg-gray-300">
            <div v-if="!chatStore.isConnected" class="text-red-500 text-center my-4">
                Disconnected. <button @click="chatStore.connect()" class="hover:underline hover:cursor-pointer">Retry</button>
            </div>
            <ul>
                <li 
                v-for="chat in filteredChatList"
                :key="chat.user.ID"
                @click="selectChatUser(chat.user)"
                >
                    <ChatItem :user-name="chat.user.name" :last-message="chat.last_message" :timestamp="displayTimeStamp(chat.timestamp)" :is-selected="chatStore.selectedChatUser?.ID === chat.user.ID"></ChatItem>
                </li>
            </ul>
            
        </div>
    </div>
</template>

