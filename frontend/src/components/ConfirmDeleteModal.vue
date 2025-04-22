<script setup>
import { defineProps, ref, reactive , onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { Button } from 'primevue';
import { toast } from 'vue3-toastify'
const props = defineProps({
    userId: Number,
    toggleModal: Function
})
const isLoading = ref(false)
const userStore = useUserStore()
const userEmail = ref('')

onMounted( async () => {
    try {
        const userId = parseInt(props.userId)
        await userStore.loadUser(userId)
        userEmail.value = userStore.selectedUser.email
    } catch (error) {
        console.log('error', error)
    }
})

const handleSubmit = async () => {
    try {
        isLoading.value = true
        const isSuccess = await userStore.deleteUser(props.userId)
        if(isSuccess){
            toast.success('User deleted successfully!', {
                autoClose: 3000,
                position: 'top-right',
            })
        }else{
            console.log('Delete user failed', error)
        }
    } catch (error) {
        toast.error(error.message, {
            autoClose: 3000,
            position: 'top-right',
        })
        console.log('Delete user failed', error)
        isLoading.value = false
    } finally {
        isLoading.value = false
        closeModal(true)
    }
}

const closeModal = (shouldRefresh) => {
    props.toggleModal(null,shouldRefresh) 
}

</script>
<template>
    <div class="fixed top-0 left-0 right-0 bottom-0 z-50 bg-black/20 flex items-center justify-center">
        <div class="container bg-white w-[50vh] rounded-xl p-6">
            <div class="flex justify-between items-center mb-3">
                <h2 class="text-xl font-bold text-black">Confirm Delete</h2>
                <Button icon="pi pi-times" @click="closeModal(false)"></Button>
            </div>
            <h3 class="text-gray-600">Are you sure you want to delete this user? This action cannot be undone.</h3>
            <h3 class="text-gray-600">Email : {{ userEmail }}</h3>
            <div class="flex">
                <button
                    @click="handleSubmit"
                    :disabled="isLoading"
                    class="w-full bg-red-600 text-white py-2 rounded-md transition mt-3 ml-1
                            hover:bg-red-800 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                    {{isLoading ? 'Deleting...' : 'Confirm'}}
                </button>
                <button
                    @click="closeModal(false)"
                    :disabled="isLoading"
                    class="w-full bg-white text-black py-2 rounded-md transition mt-3 ml-1
                            hover:bg-gray-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                    Cancel
                </button>
            </div>
            
        </div>
    </div>
</template>

