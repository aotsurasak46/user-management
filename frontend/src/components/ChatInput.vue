<script setup>
import { ref,  nextTick } from 'vue'
import { toast } from 'vue3-toastify'
const inputMessage = defineModel('inputMessage')
const chatContainer = defineModel('chatContainer')

const emit = defineEmits(['sendMessage'])

const inputRef = ref(null)

const autoResize = () => {
  const el = inputRef.value
  if (el) {
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 200) + 'px'
    nextTick(() => {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
      inputRef.value.scrollTop = inputRef.value.scrollHeight
    })
  }
}

const handleEnter = (e) => {
  if (e.shiftKey) {
    inputMessage.value += '\n'
  } else{
    handleSendMessage()
  }
  nextTick(() => autoResize())
}

const handleSendMessage = () => {
  const message = inputMessage.value.trim()
  if (!message) return
  if (message.length > 1200) {
      toast.error('Message too long. Max allowed is 1200 characters.', {
        autoClose: 3000,
        position: 'top-right',
      })
      return
  }
  emit('sendMessage', message)
  inputMessage.value = ''
  nextTick(() => autoResize())
}

</script>
<template>
    <div class="border-t border-gray-300 p-4 bg-white flex">
      <div class="relative w-full">
        <div class="flex flex-row items-end">
            <textarea
              ref="inputRef"
              v-model="inputMessage"
              @input="autoResize"
              placeholder="Message"
              @keydown.enter.prevent="handleEnter"
              class="flex-1 max-h-40 min-h-[48px] resize-none overflow-auto rounded-md border 
              border-gray-200 px-3 py-2 focus:outline-none focus:ring-1 focus:ring-black 
              text-sm leading-5 transition-all
              [&::-webkit-scrollbar]:w-2
              [&::-webkit-scrollbar-track]:rounded-full
                [&::-webkit-scrollbar-track]:bg-gray-100
                [&::-webkit-scrollbar-thumb]:rounded-full
                [&::-webkit-scrollbar-thumb]:bg-gray-300"
            ></textarea>
            <div>
              <button
                class="bg-black text-white px-3 py-1 rounded hover:bg-gray-800 disabled:opacity-50 hover:cursor-pointer hover:disabled:cursor-default hover:disabled:bg-black ml-2"
                :disabled="!inputMessage.trim()"
                @click="handleSendMessage"
              >
                <i class="pi pi pi-send text-white"></i>
              </button>
            </div>
        </div>
      </div>
    </div>
</template>
