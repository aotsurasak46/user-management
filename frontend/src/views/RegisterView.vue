<script setup>
import { ref , onMounted , watch } from 'vue'
import { useUserAccountStore } from '../stores/userAccount'
import { useRouter} from 'vue-router'
import { toast } from 'vue3-toastify'
const userAccountStore = useUserAccountStore()
const router = useRouter()
const validationErrors = ref({})
const isValid = ref(false)
const isLoading = ref(false)


const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

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

const validateConfirmPassword = (password, confirmPassword) => {
  return password === confirmPassword;
};



watch([name, email, password], () => {
  if (validationErrors.email) validationErrors.email = '';
  if (validationErrors.password) validationErrors.password = '';
  if (validationErrors.name) validationErrors.name = '';
});

watch([password, confirmPassword], () => {
  if (!validateConfirmPassword(password.value, confirmPassword.value)) {
    validationErrors.confirmPassword = 'Passwords do not match';
  } else {
    validationErrors.confirmPassword = '';
  }
});

watch([email,password], () => {
    if (validationErrors.value.email) validationErrors.value.email = '';
})

const handleRegister = async () => {
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
  if (!validatePassword(password.value)) {
      isValid.value = false
      validationErrors.value.password = 'Please enter password'
  }
  if(!validateConfirmPassword(password.value, confirmPassword.value)){
      isValid.value = false
      validationErrors.value.confirmPassword = 'Passwords do not match'
  }
  
  if(isValid.value){
    isLoading.value = true
    try {
      const isSuccess = await userAccountStore.register(email.value, password.value, name.value)
      if(isSuccess){
        localStorage.setItem('showToast', 'registered')
        router.push({name: 'login'})
      }else{
        console.log('Login failed')
      }
    } catch (error) {
      toast.error(error.message, {
        autoClose: 3000,
        position: 'top-right',
      })
      console.log('Login failed', error)
    }finally {
      isLoading.value = false
    }
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-50 px-4">
      <div class="w-full max-w-md bg-white p-8 rounded-xl shadow">
        <h2 class="text-2xl font-bold text-center mb-2">Create your account</h2>
        <p class="text-center text-gray-600 mb-6">
          Get started by creating your account
        </p>
        <div class="space-y-4">
          <div>
            <label class="flex text-sm font-medium mb-1 justify-between" >
                Name
                <label class="text-sm text-red-500" v-if="validationErrors.name">{{ validationErrors.name }}</label>
            </label>
            <input
              type="text"
              v-model="name"
              placeholder="Enter your name"
              class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black"
            />
          </div>
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
                  placeholder="Enter your password"
                  class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black"
              />
          </div>
          <div>
              <label class="flex text-sm font-medium mb-1 justify-between items-center">
              Confirm Password
              <label class="text-sm text-red-500" v-if="validationErrors.confirmPassword">{{ validationErrors.confirmPassword }}</label>
              </label>
              <input
                  type="password"
                  v-model="confirmPassword"
                  placeholder="Enter your confirm password"
                  class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-black"
              />
          </div>
          <button
              @click="handleRegister"
              :disabled="isLoading"
              class="w-full bg-black text-white py-2 rounded-md transition mt-3
                      hover:bg-gray-800 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              >
              {{ isLoading ? 'Loading...' : 'Sign up' }}
          </button>
          <p class="text-center text-sm mt-1">
              Already have an account ?
              <RouterLink :to="{ name: 'login' }" class="text-black underline">
              Login
              </RouterLink>
          </p>
          
      </div>
    </div>
  </div>
</template>

