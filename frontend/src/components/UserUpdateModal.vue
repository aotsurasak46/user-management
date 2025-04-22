<script setup>
import { defineProps, ref, reactive , onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { Button } from 'primevue';
import { toast } from 'vue3-toastify'
const props = defineProps({
    userId: Number,
    isEditMode: Boolean,
    toggleModal: Function
})
const isLoading = ref(false)
const isValid = ref(true)
const validationErrors = ref({})
const userStore = useUserStore()

const name = ref('')
const email = ref('')
const role = ref('User')
const password = ref('')

const validateName = (email) => {
    return email.trim().length !== 0;
}
const validateEmail = (email) => {
    return email.includes('@')
}
const validateEmailIsEmpty = (email) => {
    return email.trim().length !== 0;
}
const validatePassword = (password) => {
    return password && password.trim().length > 0;
};

onMounted(async () => {
    if(props.isEditMode){
        console.log('edit')
        try {
        const userId = parseInt(props.userId)
        await userStore.loadUser(userId)
        name.value = userStore.selectedUser.name
        email.value = userStore.selectedUser.email
        role.value = userStore.selectedUser.role.charAt(0).toUpperCase() + userStore.selectedUser.role.slice(1)
        } catch (error) {
        console.log('error', error)
        }
    }else{
        console.log('create')
    }
})



const handleSubmit = async () => {
    isValid.value = true
    validationErrors.value = {}
    if (!validateName(name.value)) {
      isValid.value = false
      validationErrors.value.name = 'Name can not be empty'
    }
    if (!validateEmailIsEmpty(email.value)) {
        isValid.value = false
        validationErrors.value.email = 'Email can not be empty'
    }
    else if (!validateEmail(email.value)) {
        isValid.value = false
        validationErrors.value.email = 'Email is not valid'
    }
    if (!props.isEditMode && !validatePassword(password.value)) {
        isValid.value = false
        validationErrors.value.password = 'Please enter password'
    }
    if(!isValid.value){
        return;
    }

    if (props.isEditMode && props.userId !== null) {
        await handleEditUser()
    }else{
        await handleCreateUser()
    }
}

const handleCreateUser = async () => {
    try {
        isLoading.value = true
        const isSuccess = await userStore.createUser(name.value, email.value , password.value, role.value.toLowerCase())
        if(isSuccess){
            toast.success('User created successfully!', {
                autoClose: 3000,
                position: 'top-right',
            })
        }else{
            console.log('Create user failed', error)
        }
    } catch (error) {
        toast.error(error.message, {
            autoClose: 3000,
            position: 'top-right',
        })
        console.log('Create user failed', error)
        isLoading.value = false
    } finally {
        isLoading.value = false
        closeModal(true)
    }
}

const handleEditUser = async () => {
    try {
        isLoading.value = true
        const isSuccess = await userStore.editUser(name.value, email.value , role.value.toLowerCase(), props.userId)
        if(isSuccess){
            toast.success('User updated successfully!', {
                autoClose: 3000,
                position: 'top-right',
            })
        }else{
            console.log('Update user failed', error)
        }
    } catch (error) {
        toast.error(error.message, {
            autoClose: 3000,
            position: 'top-right',
        })
        console.log('Update user failed', error)
    } finally {
        isLoading.value = false
        closeModal(true)
    }
}

const closeModal = (shouldRefresh) => {
    props.toggleModal(null, false, shouldRefresh) 
}

</script>
<template>
    <div class="fixed top-0 left-0 right-0 bottom-0 z-50 bg-black/20 flex items-center justify-center">
        <div class="container bg-white w-[50vh] rounded-xl p-6">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold text-black">{{ props.isEditMode ? 'Edit user data' : 'Add new user' }}</h2>
                <Button icon="pi pi-times" @click="closeModal(false)"></Button>
            </div>
            <div>
                <label class="flex text-sm font-medium mb-2 justify-between items-center">
                    Name
                    <label class="text-sm text-red-500" v-if="validationErrors.name">{{ validationErrors.name }}</label>
                </label>
                <input
                    type="text"
                    v-model="name"
                    class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black mb-2"
                />
            </div>
            <div>
                <label class="flex text-sm font-medium mb-2 justify-between items-center">
                    Email
                    <label class="text-sm text-red-500" v-if="validationErrors.email">{{ validationErrors.email }}</label>
                </label>
                <input
                    type="text"
                    v-model="email"
                    class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black mb-2"
                />
            </div>
            <div v-if="!props.isEditMode">
                <label class="flex text-sm font-medium mb-2 justify-between items-center">
                    Password
                    <label class="text-sm text-red-500" v-if="validationErrors.password">{{ validationErrors.password }}</label>
                </label>
                <input
                    type="password"
                    v-model="password"
                    class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black mb-2"
                />
            </div>
            <div>
                <label class="flex text-sm font-medium mb-2 justify-between items-center">
                Role
                </label>
                <select v-model="role" class='mb-2'>
                    <option>User</option>
                    <option>Admin</option>
                </select>
            </div>
            <div class="flex">
                <button
                    @click="handleSubmit"
                    :disabled="isLoading"
                    class="w-full bg-black text-white py-2 rounded-md transition mt-3 mr-1
                            hover:bg-gray-800 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                    {{ isLoading ? 'Loading...' : props.isEditMode ? 'Update' : 'Create' }}
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

