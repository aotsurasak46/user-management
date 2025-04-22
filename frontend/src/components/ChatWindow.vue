<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useChatStore } from '@/stores/chat';
import { useUserAccountStore } from '@/stores/userAccount';
import { toast } from 'vue3-toastify';
import ChatInput from './ChatInput.vue';
const chatStore = useChatStore()
const userAccountStore = useUserAccountStore()
const isLoading = ref(false)
const errorMessage = ref('')

const inputMessage = ref('')
const chatContainer = ref(null)

watch(
  () => chatStore.chatMessages.length,
  () => {
    nextTick(() => {
      if (chatContainer.value) {
        chatContainer.value.scrollTop = chatContainer.value.scrollHeight
      }
    })
  }
)
watch(
    () => chatStore.selectedChatUser,
    () => {
        fetchChatMessages()
    },
    { deep: true }
);

onMounted(async() => {
    fetchChatMessages()
})


const handleSendMessage = () => {
    const message = inputMessage.value.trim()
    chatStore.sendMessage(message)
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

const fetchChatMessages = async () => {
  try {
    isLoading.value = true
    await chatStore.getChatHistory()
  } catch (error) {
    errorMessage.value = error.message
    console.log('Failed to load chat history : ', error)
  }finally{
    await nextTick()
    if (chatContainer.value) {
        chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
    isLoading.value = false
  }
}

const isOwnMessage = (msg) => {
  const userId = Number(userAccountStore.user.id)
  const fromId = Number(msg.from_id) 
  return fromId === userId
}



</script>
<template>
    <div class="w-full max-h-screen">
        <div v-if="chatStore.selectedChatUser" class="flex flex-col h-screen">
            <div class="flex p-4 bg-white border-b border-gray-300 ">
                <div class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-500">
                    <span>{{ chatStore.selectedChatUser.name.charAt(0).toUpperCase()  }}</span>
                </div>
                <div class="flex flex-col justify-center px-4 min-w-0 flex-1 overflow-hidden">
                    <span class="text-sm text-gray-900 font-medium truncate">{{ chatStore.selectedChatUser.name }}</span>
                    <span class="text-xs text-gray-500 font-normal truncate">{{ chatStore.selectedChatUser.email }}</span>
                </div>
            </div>
            <div class="flex-1 overflow-hidden overflow-y-auto
                p-2 [&::-webkit-scrollbar]:w-2
                [&::-webkit-scrollbar-track]:rounded-full
                [&::-webkit-scrollbar-track]:bg-gray-100
                [&::-webkit-scrollbar-thumb]:rounded-full
                [&::-webkit-scrollbar-thumb]:bg-gray-300" 
                ref="chatContainer"> 
                <div
                v-for="msg in chatStore.chatMessages"
                :key="msg.tempId || msg.ID"
                :class="[
                    'flex',
                    isOwnMessage(msg) ? 'justify-end' : 'justify-start'
                ]"
                >
                    <div
                        :class="[
                        'max-w-xs px-4 py-2 rounded-2xl text-sm my-0.5 break-words overflow-hidden',
                        isOwnMessage(msg)
                            ? 'bg-black text-white rounded-br-none'
                            : 'bg-gray-300 text-gray-800 rounded-bl-none'
                        ]"
                    >
                        <p class="line-clamp-25">{{ msg.content }}</p>
                        <div class="text-xs text-right mt-1 opacity-70">
                        <span>{{ displayTimeStamp(msg.timestamp) }}</span>
                        <span v-if="isOwnMessage(msg)">
                            <template v-if="msg.status === 'sending'">
                                <span class="ml-1 italic text-yellow-500">sending...</span>
                            </template>
                            <template v-else-if="msg.status === 'sent'">
                                <span class="ml-1 text-green-500">✓</span>
                            </template>
                        </span>
                        </div>
                    </div>
                </div>
            </div>
            <ChatInput v-model:input-message="inputMessage" v-model:chat-container="chatContainer" @send-message="handleSendMessage"></ChatInput>
        </div>
        <div v-else class="flex flex-col h-full justify-center items-center">
            <img alt="Logo project" class="logo" src="../assets/images/messages.svg" width="200" />
            <h1 class="text-2xl">Your Messages</h1>
            <h1 class="text-lg text-gray-600">Select a chat to start messaging</h1>
        </div>
    </div>  
</template>
  