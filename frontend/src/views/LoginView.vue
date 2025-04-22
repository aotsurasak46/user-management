<script setup>
import {ref, onMounted, watch} from 'vue'
import {useUserAccountStore} from '../stores/userAccount'
import {useRouter} from 'vue-router'
import {toast} from 'vue3-toastify'


const userAccountStore = useUserAccountStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const validationErrors = ref({})

const isValid = ref(false)
const isLoading = ref(false)

const validateEmail = (email) => {
    return email.includes('@')
}

const validateEmailIsEmpty = (email) => {
    return email.trim().length !== 0;
}

const validatePassword = (password) => {
    return password && password.trim().length > 0;
};

watch([email,password], () => {
    if (validationErrors.value.email) validationErrors.value.email = '';
    if (validationErrors.value.password) validationErrors.value.password = '';
})

const handleLogin = async () => {
    isValid.value = true
    validationErrors.value = {}
    if (!validateEmailIsEmpty(email.value)) {
        isValid.value = false
        validationErrors.value.email = 'Email can not be empty'
    }
    else if (!validateEmail(email.value)) {
        isValid.value = false
        validationErrors.value.email = 'Email is not valid'
    }
    if (!validatePassword(password.value)) {
        isValid.value = false
        validationErrors.value.password = 'Please enter password'
    }
    if(isValid.value){
        isLoading.value = true
        try {
            const isSuccess = await userAccountStore.login(email.value, password.value)
            if (isSuccess){
                router.push({ name: 'user-list' })
            }else {
                console.log('Login failed')
            }
        } catch (error) {
            toast.error(error.message, {
                autoClose: 3000,
                position: 'top-right',
            })
            console.log('Login failed', error)
        } finally {
            isLoading.value = false
        }
    }
}


onMounted( async () => {
    if (localStorage.getItem('showToast') === 'registered') {
        toast.success('Account created successfully!', {
            autoClose: 3000,
            position: 'top-right',
        })
        localStorage.removeItem('showToast')
    }
    try {
        await userAccountStore.checkAuth()
        if (userAccountStore.isAuthenticated) {
            router.push({ name: 'user-list' })
        }
    } catch (error) {
        console.log(error)
    }
})


</script>
<template>
    <div class="flex items-center justify-center min-h-screen bg-gray-50 px-4">
      <div class="w-full max-w-md bg-white p-8 rounded-xl shadow">
        <h2 class="text-2xl font-bold text-center mb-2">Welcome back</h2>
        <p class="text-center text-gray-600 mb-6">
          Enter your email below to login to your account
        </p>
        <div class="space-y-4">
          <div>
            <label class="flex text-sm font-medium mb-1 justify-between" >
                Email
                <label class="text-sm text-red-500" v-if="validationErrors.email">{{ validationErrors.email }}</label>
            </label>
            <input
              type="email"
              v-model="email"
              placeholder="m@example.com"
              class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black"
            />
          </div>
            <div>
                <label class="flex text-sm font-medium mb-1 justify-between items-center">
                Password
                <label class="text-sm text-red-500" v-if="validationErrors.password">{{ validationErrors.password }}</label>
                </label>
                <input
                    type="password"
                    v-model="password"
                    class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black"
                />
            </div>
            <button
                @click="handleLogin"
                :disabled="isLoading"
                class="w-full bg-black text-white py-2 rounded-md transition mt-3
                        hover:bg-gray-800 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                {{ isLoading ? 'Logging in...' : 'Login' }}
            </button>
            <p class="text-center text-sm mt-1">
                Don't have an account?
                <RouterLink :to="{ name: 'register' }" class="text-black underline">
                Sign up
                </RouterLink>
            </p>
        </div>
      </div>
    </div>
</template>
  