import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useUserAccountStore } from './userAccount'
import axios from 'axios'

const WS_URL = import.meta.env.VITE_WS_URL 
const BASE_URL = import.meta.env.VITE_API_BASE_URL 

export const useChatStore = defineStore('chat', {
  state: () => ({
    socket:  null,
    isConnected: false,
    reconnectAttempts: 0,
    maxReconnectAttempts : 5,

    chatList: [],
    selectedChatUser: null,
    chatMessages: [],
  }),
  actions:{
    attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max reconnect attempts reached. Giving up.');
            return;
        }
        console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);
        setTimeout(() => {
            this.reconnectAttempts++;
            this.connect();
        }, 2000); 
    },

    connect() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) return;
    
        this.socket = new WebSocket(WS_URL); 
    
        this.socket.onopen = () => {
            this.isConnected = true; 
            this.reconnectAttempts = 0; 
            console.log('WebSocket connected!');
        };
    
        this.socket.onmessage = (event) => {    
            try {
                const payload = JSON.parse(event.data);
    
                if (payload.type === 'sent' && payload.data?.tempId) {
                    const index = this.chatMessages.findIndex(msg => msg.tempId === payload.data.tempId);
                    if (index !== -1) {
                        if (index !== -1) {
                            this.chatMessages[index] = {
                                id: payload.data.ID,            
                                from_id: payload.data.from_id,
                                to_id: payload.data.to_id,
                                content: payload.data.content,
                                timestamp: payload.data.timestamp,      
                                tempId: payload.data.tempId,
                                status: 'sent'
                            }
                            this.getChatList()
                        }
                    }
                    return;
                }
    
                const data = payload.data;
                if (
                    this.selectedChatUser &&
                    (data.from_id === this.selectedChatUser.ID || data.to_id === this.selectedChatUser.ID)
                ) {
                    this.chatMessages.push({
                        from_id: data.fromId,
                        to_id: data.toId,
                        content: data.content,
                        timestamp: data.timestamp,
                        status: 'received',
                    });
                }

                this.getChatList()

                
    
            } catch (error) {
                console.error('Error parsing message:', error);
            }
        };
    
        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.attemptReconnect();
        };
    
        this.socket.onclose = () => {
            this.isConnected = false; 
            console.log('WebSocket connection closed.');
            this.attemptReconnect(); 
        };
    },


    closeConnection() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
            this.isConnected = false;
            console.log('WebSocket connection closed manually.');
        }
    },
    
    sendMessage(content) {
        const authStore = useUserAccountStore();
      
        const fromId = authStore.user?.id;
        if (!fromId) {
          console.error('User is not authenticated or user ID is missing');
          return;
        }

        const message = {
          tempId: Date.now() + '-' + Math.random().toString(36).substring(2, 8),
          from_id: fromId,
          to: this.selectedChatUser.ID,
          content: content,
          status: 'sending',
        };

        this.chatMessages.push(message);
      
        const payload = {
          to: message.to,
          content: message.content,
          tempId: message.tempId,
        };
        
        try {
          this.socket.send(JSON.stringify(payload));
        } catch (err) {
          console.error('Failed to send message:', err);
        }
      },
      

    async getChatHistory() {
        try {
            if (this.selectedChatUser === null) return
            const response = await axios.get(`${BASE_URL}/api/v1/messages/${this.selectedChatUser.ID}` ,{
            withCredentials: true,
        })
            this.chatMessages = response.data
        }catch(error){
            console.log('error', error)
            let errorMessage = 'Something went wrong. Please try again.'
            if (error.response && error.response.data && error.response.data.error) {
                errorMessage = error.response.data.error
            }
            throw new Error(errorMessage)
        }
    },

    async getChatList() {
        try {
            const response = await axios.get(`${BASE_URL}/api/v1/conversations` ,{
            withCredentials: true,
        })
            this.chatList = response.data
        }catch(error){
            console.log('error', error)
            let errorMessage = 'Something went wrong. Please try again.'
            if (error.response && error.response.data && error.response.data.error) {
                errorMessage = error.response.data.error
            }
            throw new Error(errorMessage)
        }
    },

    setSelectedChatUser(user) {
        this.selectedChatUser = user
    },

  }
})
